#include "socks5_proxy.h"
#include <boost/lexical_cast.hpp>
#include <boost/asio/yield.hpp>

template<typename Stream, typename Buffer, typename Handler>
struct async_read_socks5_auth_request_op : asio::coroutine
{
	Stream& m_stream;
	Buffer& m_buffer;
	Handler m_handler;

	async_read_socks5_auth_request_op(Stream& stream, Buffer& buffer, Handler& handler) : m_stream(stream), m_buffer(buffer), m_handler(std::move(handler))
	{

	}

	async_read_socks5_auth_request_op(const async_read_socks5_auth_request_op& other) : asio::coroutine(other), m_stream(other.m_stream), m_buffer(other.m_buffer), m_handler(std::move(other.m_handler))
	{

	}

	void operator()(const boost::system::error_code& ec = boost::system::error_code(), std::size_t bytes_transfered = 0)
	{
		if (ec)
		{
			m_handler(ec, 0);
			return;
		}
		reenter(this)
		{
			yield asio::async_read(m_stream, asio::buffer(&m_buffer.version, 1), *this);
			yield asio::async_read(m_stream, asio::buffer(&m_buffer.nmethods, 1), *this);
			yield asio::async_read(m_stream, asio::buffer(m_buffer.methods, m_buffer.nmethods), *this);
			m_handler(ec, 0);
		}
	}

};

template<typename Stream, typename Buffer, typename Handler>
BOOST_ASIO_INITFN_RESULT_TYPE(Handler, void(boost::system::error_code, std::size_t))
async_read_socks5_auth_request(Stream& stream, Buffer& buffer, Handler&& handler)
{
	asio::detail::async_result_init<Handler, void(boost::system::error_code, std::size_t)> init(std::move(handler));
	async_read_socks5_auth_request_op<Stream, Buffer, Handler>(stream, buffer, init.handler)();
	return init.result.get();
}

template<typename Stream, typename Buffer, typename Handler>
struct async_write_socks5_auth_reply_op : asio::coroutine
{
	Stream& m_stream;
	const Buffer& m_buffer;
	Handler m_handler;

	async_write_socks5_auth_reply_op(Stream& stream, const Buffer& buffer, Handler& handler) : m_stream(stream), m_buffer(buffer), m_handler(std::move(handler))
	{

	}

	async_write_socks5_auth_reply_op(const async_write_socks5_auth_reply_op& other) : asio::coroutine(other), m_stream(other.m_stream), m_buffer(other.m_buffer), m_handler(std::move(other.m_handler))
	{

	}

	void operator()(const boost::system::error_code& ec = boost::system::error_code(), std::size_t bytes_transfered = 0)
	{
		if (ec)
		{
			m_handler(ec, 0);
			return;
		}
		reenter(this)
		{
			yield asio::async_write(m_stream, asio::buffer(&m_buffer.version, 1), *this);
			yield asio::async_write(m_stream, asio::buffer(&m_buffer.method, 1), *this);
			m_handler(ec, 0);
		}
	}

};

template<typename Stream, typename Buffer, typename Handler>
BOOST_ASIO_INITFN_RESULT_TYPE(Handler, void(boost::system::error_code, std::size_t))
async_write_socks5_auth_reply(Stream& stream, const Buffer& buffer, Handler&& handler)
{
	asio::detail::async_result_init<Handler, void(boost::system::error_code, std::size_t)> init(std::move(handler));
	async_write_socks5_auth_reply_op<Stream, Buffer, Handler>(stream, buffer, init.handler)();
	return init.result.get();
}

template<typename Stream, typename Buffer, typename Handler>
struct async_read_socks5_request_op : asio::coroutine
{
	Stream& m_stream;
	Buffer& m_buffer;
	Handler m_handler;

	async_read_socks5_request_op(Stream& stream, Buffer& buffer, Handler& handler) : m_stream(stream), m_buffer(buffer), m_handler(std::move(handler))
	{

	}

	async_read_socks5_request_op(const async_read_socks5_request_op& other) : asio::coroutine(other), m_stream(other.m_stream), m_buffer(other.m_buffer), m_handler(std::move(other.m_handler))
	{

	}

	void operator()(const boost::system::error_code& ec = boost::system::error_code(), std::size_t bytes_transfered = 0)
	{
		if (ec)
		{
			m_handler(ec, 0);
			return;
		}
		reenter(this)
		{
			yield asio::async_read(m_stream, asio::buffer(&m_buffer.version, 1), *this);
			yield asio::async_read(m_stream, asio::buffer(&m_buffer.cmd, 1), *this);
			yield asio::async_read(m_stream, asio::buffer(&m_buffer.reserved, 1), *this);
			yield asio::async_read(m_stream, asio::buffer(&m_buffer.atype, 1), *this);
			switch (m_buffer.atype)
			{
			case 1:
				yield asio::async_read(m_stream, asio::buffer(m_buffer.addr.ipv4, 4), *this);
				break;
			case 3:
				yield asio::async_read(m_stream, asio::buffer(&m_buffer.addr.domain.len, 1), *this);
				yield asio::async_read(m_stream, asio::buffer(&m_buffer.addr.domain.host[0], m_buffer.addr.domain.len), *this);
				break;
			case 4:
				yield asio::async_read(m_stream, asio::buffer(m_buffer.addr.ipv6, 16), *this);
				break;
			default:
				break;
			}
			m_handler(ec, 0);
		}
	}

};

template<typename Stream, typename Buffer, typename Handler>
BOOST_ASIO_INITFN_RESULT_TYPE(Handler, void(boost::system::error_code, std::size_t))
async_read_socks5_request(Stream& stream, Buffer& buffer, Handler&& handler)
{
	asio::detail::async_result_init<Handler, void(boost::system::error_code, std::size_t)> init(std::move(handler));
	async_read_socks5_request_op<Stream, Buffer, Handler>(stream, buffer, init.handler)();
	return init.result.get();
}

socks5_proxy::socks5_proxy(asio::io_service& io)
	: io_(io)
	, socket_(io_)
{

}

tcp::socket& socks5_proxy::socket()
{
	return socket_;
}

void socks5_proxy::start(boost::system::error_code ec /*= boost::system::error_code()*/, std::size_t bytes_transfered /*= 0*/)
{
	reenter(client_coro_)
	{
		yield async_read_socks5_auth_request(socket_, socks5_auth_method_request_, boost::bind(&socks5_proxy::start, shared_from_this(), _1, _2));
		// send auth method
		socks5_auth_method_reply_.version = 5;
		socks5_auth_method_reply_.method = 0;
		yield async_write_socks5_auth_reply(socket_, socks5_auth_method_reply_, boost::bind(&socks5_proxy::start, shared_from_this(), _1, _2));
		yield async_read_socks5_request(socket_, socks5_request_, boost::bind(&socks5_proxy::start, shared_from_this(), _1, _2));
		switch ((CommandType)socks5_request_.cmd)
		{
		case kConnect:
			yield async_connect();
			break;
		case kBind:
			break;
		case kUdpAssociate:
			break;
		default:
			break;
		}
	}
	
}
