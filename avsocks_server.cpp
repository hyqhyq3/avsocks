#include "avsocks_server.h"
#include "socks5_proxy.h"

avsocks_server::avsocks_server(asio::io_service& io, uint16_t port)
	: io_(io)
	, acceptor_(io_, tcp::endpoint(tcp::v4(), port))
{

}

void avsocks_server::start(asio::yield_context yield)
{
	while (true)
	{
		auto proxy = boost::make_shared<socks5_proxy>(io_);
		boost::system::error_code ec;
		acceptor_.async_accept(proxy->socket(), yield[ec]);
		if (!ec)
		{
			std::cout << "accepted a connection" << std::endl;
			asio::spawn(io_, boost::bind(&socks5_proxy::go, proxy, _1));
		}
		else
		{
			break;
		}
	}
}
