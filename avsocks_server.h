#pragma once
#include "avsocks.h"

class avsocks_server
	: boost::noncopyable
{
	asio::io_service& io_;
	tcp::acceptor acceptor_;
public:
	avsocks_server(asio::io_service& io, uint16_t port);

	void start(asio::yield_context yield);
};

