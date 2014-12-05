#pragma once
#include "avsocks.h"

struct socks5_auth_method_request
{
	uint8_t version;
	uint8_t nmethods;
	std::vector<uint8_t> methods;
};

struct socks5_auth_method_reply
{
	uint8_t version;
	uint8_t method;
};

union socks5_address
{
	uint8_t ipv4[4];
	struct  
	{
		uint8_t len;
		char host[255];
	} domain;
	uint8_t ipv6[16];
};

struct socks5_request
{
	uint8_t version;
	uint8_t cmd;
	uint8_t reserved;
	uint8_t atype;
	socks5_address addr;
	uint16_t port;
};

struct socks5_reply
{
	uint8_t version;
	uint8_t cmd;
	uint8_t reserved;
	uint8_t atype;
	socks5_address addr;
	uint16_t port;
};

class socks5_proxy
	: public boost::enable_shared_from_this<socks5_proxy>
	, boost::noncopyable
{
	asio::io_service& io_;
	tcp::socket socket_;
	asio::coroutine client_coro_;
	asio::coroutine proxy_coro_;
	asio::streambuf send_to_client_buffer_;
	asio::streambuf recv_from_client_buffer_;
	socks5_auth_method_request socks5_auth_method_request_;
	socks5_auth_method_reply socks5_auth_method_reply_;
	socks5_request socks5_request_;

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

protected:
	void pipe(tcp::socket& local, tcp::socket& remote);

public:
	socks5_proxy(asio::io_service& io);

	tcp::socket& socket();

	void start(boost::system::error_code ec = boost::system::error_code(), std::size_t bytes_transfered = 0);
};

