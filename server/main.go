//server

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/kless/goconfig/config"
	"io"
	"log"
	"net"
)

var clientCipher cipher.Block
var serverCipher cipher.Block
var listen string
var iv = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

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

func handleStream(t net.Conn, f net.Conn, s cipher.Stream) {
	in := make([]byte, 1024*4)
	out := make([]byte, 1024*4)
	for {
		n, err := f.Read(in)
		if err == nil || err == io.EOF {
			s.XORKeyStream(out, in[:n])
			_, err := t.Write(out[:n])
			if err != nil {
				f.Close()
			}
		} else {
			t.Close()
			break
		}
	}
}

func handle(conn net.Conn) {
	ccfb := cipher.NewCFBDecrypter(clientCipher, iv)
	scfb := cipher.NewCFBEncrypter(serverCipher, iv)
	out := make([]byte, 1024)
	readAndDecode(conn, out[:4], ccfb)
	if out[1] != 1 {
		log.Print("unsupport command")
		return
	}
	addr := new(net.TCPAddr)
	addr.IP = make([]byte, 4)
	switch out[3] { //atyp
	case 1: //ipv4
		readAndDecode(conn, out[:4], ccfb)
		copy(addr.IP, out[:4])
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
		copy(addr.IP, addrs[0].To4())
	default:
		log.Print("unsupport address type")
		return
	}
	readAndDecode(conn, out[:2], ccfb)
	addr.Port = nl2int(out)

	sConn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Print("cannot connect to server", addr.String())
		encodeAndWrite(conn, out[:10], scfb)
		conn.Close()
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
	go handleStream(sConn, conn, ccfb)
	go handleStream(conn, sConn, scfb)
}

func main() {
	//read config
	c, _ := config.ReadDefault("config.ini")
	listen, _ = c.String("server", "listen")
	log.Printf("listen:%s", listen)
	ck, _ := c.String("encrypto", "client-key")
	log.Printf("client-key:%s", ck)
	sk, _ := c.String("encrypto", "server-key")
	log.Printf("server-key:%s", sk)
	clientCipher, _ = aes.NewCipher([]byte(ck))
	serverCipher, _ = aes.NewCipher([]byte(sk))

	lsn, e := net.Listen("tcp", listen)
	if e != nil {
		log.Fatal(e)
	}
	for {
		conn, e := lsn.Accept()
		if e != nil {
			log.Print(e)
		}
		go handle(conn)
	}
}
