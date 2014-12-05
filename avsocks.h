#include <cstdint>
#include <cstddef>
#include <iostream>
#include <boost/asio.hpp>
namespace asio = boost::asio;
using asio::ip::tcp;
#include <boost/noncopyable.hpp>
#include <boost/bind.hpp>
#include <boost/shared_ptr.hpp>
#include <boost/make_shared.hpp>
#include <boost/enable_shared_from_this.hpp>