//server

package handler

import (
	"crypto/cipher"
	"log"
	"net"
)

type Server struct {
	ClientCipher, ServerCipher cipher.Block
}

func readAndDecode(conn net.Conn, out []byte, s cipher.Stream) {
	in := make([]byte, len(out))
	conn.Read(in)
	s.XORKeyStream(out, in)
}

func encodeAndWrite(conn net.Conn, in []byte, s cipher.Stream) {
	out := make([]byte, len(in))
	s.XORKeyStream(out, in)
	conn.Write(out)
}

func int2nl(i int) []byte {
	b := make([]byte, 2)
	b[0] = byte((i & 0xffff) >> 8)
	b[1] = byte(i & 0xff)
	return b
}

func nl2int(nl []byte) int {
	return int(nl[0])*256 + int(nl[1])
}

func (s *Server) Handle(conn net.Conn) {
	defer conn.Close()
	ccfb := cipher.NewCFBDecrypter(s.ClientCipher, iv)
	scfb := cipher.NewCFBEncrypter(s.ServerCipher, iv)
	out := make([]byte, 1024)
	readAndDecode(conn, out[:4], ccfb)
	if out[1] != 1 {
		log.Print("unsupport command")
		return
	}
	addr := new(net.TCPAddr)
	switch out[3] { //atyp
	case 1: //ipv4
		readAndDecode(conn, out[:4], ccfb)
		addr.IP = net.IPv4(out[0], out[1], out[2], out[3])
	case 3: //domain
		readAndDecode(conn, out[:1], ccfb)
		l := out[0]
		readAndDecode(conn, out[:l], ccfb)
		host := string(out[:l])
		addrs, err := net.LookupIP(host)
		if err != nil || len(addrs) == 0 {
			log.Print("cannot resolve ", host)
			return
		}
		addr.IP = addrs[0]
	default:
		log.Print("unsupport address type")
		return
	}
	readAndDecode(conn, out[:2], ccfb)
	addr.Port = nl2int(out)

	sConn, err := net.DialTCP("tcp", nil, addr)
	defer sConn.Close()
	if err != nil {
		log.Print("cannot connect to server", addr.String())
		encodeAndWrite(conn, out[:10], scfb)
		return
	}
	response := make([]byte, 10)
	response[0] = 5
	response[1] = 0
	response[2] = 0
	response[3] = 1
	copy(response[4:], sConn.LocalAddr().(*net.TCPAddr).IP)
	copy(response[8:], int2nl(sConn.LocalAddr().(*net.TCPAddr).Port))
	encodeAndWrite(conn, response, scfb)
	go HandleStream(conn, sConn, scfb)
	HandleStream(sConn, conn, ccfb)
}
