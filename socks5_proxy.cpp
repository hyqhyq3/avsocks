#include "socks5_proxy.h"
#include <boost/lexical_cast.hpp>

socks5_proxy::socks5_proxy(asio::io_service& io)
	: io_(io)
	, socket_(io_)
{

}

tcp::socket& socks5_proxy::socket()
{
	return socket_;
}

void socks5_proxy::go(asio::yield_context yield)
{
	auto self = shared_from_this();
	try {
		std::array<uint8_t, 256> buf;
		asio::async_read(socket_, asio::buffer(buf, 2), yield);
		if (buf[0] != 5)
		{
			return;
		}
		int methods = buf[1];
		asio::async_read(socket_, asio::buffer(buf, methods), yield);
		// find no auth method
		if (!std::find(&buf[0], &buf[methods], 0))
		{
			std::cout << "no suitable auth method" << std::endl;
			return;
		}
		std::array<uint8_t, 2> auth = { 5, 0 };
		asio::async_write(socket_, asio::buffer(auth), yield);

		// read client request header
		std::size_t n = asio::async_read(socket_, asio::buffer(header_, 4), yield);
		if (header_[0] != 5)
		{
			return;
		}
		// cmd
		switch (header_[1])
		{
		case kConnect:
			handle_connect(std::move(yield));
			break;
		case kBind:
			break;
		case kUdpAssociate:
			break;
		}
	}
	catch (std::exception e)
	{
		std::cout << e.what() << std::endl;
		socket_.close();
	}
}

void socks5_proxy::handle_connect(asio::yield_context yield)
{
	tcp::endpoint destination;
	if (header_[3] == kIPv4)
	{
		uint32_t ip = 0;
		uint16_t port = 0;
		asio::async_read(socket_, asio::buffer(&ip, 4), yield);
		asio::async_read(socket_, asio::buffer(&port, 2), yield);
		
		port = htons(port);

		destination = tcp::endpoint(asio::ip::address_v4(ip), port);
	}
	else if (header_[3] == kDomain)
	{
		uint8_t len = 0;
		std::string host;
		uint16_t port = 0;
		asio::async_read(socket_, asio::buffer(&len, 1), yield);
		host.resize(len);
		asio::async_read(socket_, asio::buffer(&host[0], len), yield);
		asio::async_read(socket_, asio::buffer(&port, 2), yield);
		
		port = htons(port);

		tcp::resolver resolver(io_);
		destination = *resolver.async_resolve(tcp::resolver::query(host, boost::lexical_cast<std::string>(port)), yield);
	}
	else if (header_[3] = kIPv6)
	{
		std::array<uint8_t, 16> ip;
		uint16_t port = 0;
		asio::async_read(socket_, asio::buffer(ip, 16), yield);
		asio::async_read(socket_, asio::buffer(&port, 2), yield);

		port = htons(port);

		destination = tcp::endpoint(asio::ip::address_v6(ip), port);
	}
	std::cout << "connecting " << destination << std::endl;
	boost::system::error_code ec;
	tcp::socket transport(io_);
	transport.async_connect(destination, yield[ec]);
	if (ec)
	{
		std::cout << ec.message() << std::endl;
		return;
	}
	std::cout << "connected" << std::endl;

	auto local_endpoint = transport.local_endpoint();
	std::array<uint8_t, 4> reply = { 5, 0, 0 };
	if (local_endpoint.protocol() == tcp::v4())
	{
		reply[3] = kIPv4;
		std::array<uint8_t, 4> bind_address = local_endpoint.address().to_v4().to_bytes();
		uint16_t bind_port = local_endpoint.port();
		std::vector<asio::const_buffer> bufs = { asio::buffer(reply), asio::buffer(bind_address), asio::buffer(&bind_port, 2) };
		asio::async_write(socket_, asio::buffer(bufs), yield);
	}
	else if (local_endpoint.protocol() == tcp::v6())
	{
		reply[3] = kIPv6;
		std::array<uint8_t, 16> bind_address = local_endpoint.address().to_v6().to_bytes();
		uint16_t bind_port = local_endpoint.port();
		std::vector<asio::const_buffer> bufs = { asio::buffer(reply), asio::buffer(bind_address), asio::buffer(&bind_port, 2) };
		asio::async_write(socket_, asio::buffer(reply), yield);
	}

	asio::spawn(yield, boost::bind(&socks5_proxy::pipe, this, _1, boost::ref(transport), boost::ref(socket_)));
	asio::spawn(yield, boost::bind(&socks5_proxy::pipe, this, _1, boost::ref(socket_), boost::ref(transport)));
}

void socks5_proxy::pipe(asio::yield_context yield, tcp::socket& local, tcp::socket& remote)
{
	auto self = shared_from_this();
	try {
		std::array<uint8_t, 1024> buf;
		while (true)
		{
			std::size_t n = local.async_read_some(asio::buffer(buf), yield);
			asio::async_write(remote, asio::buffer(buf, n), yield);
		}
	}
	catch (std::exception e)
	{
		local.close();
		remote.close();
	}
}
