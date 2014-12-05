#include "avsocks.h"
#include "avsocks_server.h"

int main(int argc, char* argv[])
{
	asio::io_service io;

	boost::make_shared<avsocks_server>(io, 8000)->start();
	io.run();
}