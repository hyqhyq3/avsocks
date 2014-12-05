#pragma once
#include "avsocks.h"

class socks5_proxy;

class avsocks_server
	: public boost::enable_shared_from_this<avsocks_server>
	, boost::noncopyable
{
	asio::io_service& io_;
	tcp::acceptor acceptor_;
	asio::coroutine coro_;
	boost::shared_ptr<socks5_proxy> new_client_;

	std::function<void(boost::system::error_code)> callback() { return  boost::bind(&avsocks_server::start, shared_from_this(), _1); }
public:
	avsocks_server(asio::io_service& io, uint16_t port);

	void start(boost::system::error_code ec = boost::system::error_code());

};

