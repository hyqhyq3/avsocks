package handler

import (
	"net"
)

type Handler interface {
	Handle(net.Conn)
}
