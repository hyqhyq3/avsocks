#include "avsocks_server.h"
#include "socks5_proxy.h"
#include <boost/asio/yield.hpp>

avsocks_server::avsocks_server(asio::io_service& io, uint16_t port)
	: io_(io)
	, acceptor_(io_, tcp::endpoint(tcp::v4(), port))
{

}

void avsocks_server::start(boost::system::error_code ec /* = boost::system::error_code() */)
{
	reenter(coro_)
	{
		do
		{
			new_client_.reset(new socks5_proxy(io_));
			yield acceptor_.async_accept(new_client_->socket(), callback());
			if (!ec)
			{
				std::cout << "accepted a connection" << std::endl;
				fork new_client_->start();
			}
			else
			{
				break;
			}
		} while (coro_.is_parent());
	}
}
