#include "avsocks.h"
#include "avsocks_server.h"

int main(int argc, char* argv[])
{
	asio::io_service io;

	auto server = boost::make_shared<avsocks_server>(io, 8000);
	asio::spawn(io, boost::bind(&avsocks_server::start, server, _1));
	io.run();
}