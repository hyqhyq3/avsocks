#pragma once
#include "avsocks.h"

class socks5_proxy
	: public boost::enable_shared_from_this<socks5_proxy>
	, boost::noncopyable
{
	asio::io_service& io_;
	tcp::socket socket_;

	enum CommandType {
		kConnect = 1,
		kBind = 2,
		kUdpAssociate = 3,
	};

	enum AddressType {
		kIPv4 = 1,
		kDomain = 3,
		kIPv6 = 4,
	};

	std::array<uint8_t, 4> header_;

	void handle_connect(asio::yield_context yield);

protected:
	void pipe(asio::yield_context yield, tcp::socket& local, tcp::socket& remote);

public:
	socks5_proxy(asio::io_service& io);

	tcp::socket& socket();

	void go(asio::yield_context yield);
};

