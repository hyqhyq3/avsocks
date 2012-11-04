//client

package main

import (
	"crypto/aes"
	"crypto/cipher"

	"errors"
	"github.com/kless/goconfig/config"
	"io"
	"log"
	"net"
)

var clientCipher, serverCipher cipher.Block
var iv = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15} //length must be 16

var listen, server string

func handshake(conn net.Conn) (err error) {
	b := make([]byte, 2)
	_, err = conn.Read(b)
	if err != nil {
		return
	}

	if b[0] != 5 {
		return errors.New("socks version incorrect")
	}

	b = make([]byte, b[1]) //auth methods
	conn.Read(b)
	conn.Write([]byte("\x05\x00")) //no auth
	return
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
	err := handshake(conn)
	if err != nil {
		log.Print("error when handshake")
		conn.Close()
		return
	}
	ccfb := cipher.NewCFBEncrypter(clientCipher, iv)
	scfb := cipher.NewCFBDecrypter(serverCipher, iv)
	sConn, err := net.Dial("tcp", server)
	if err != nil {
		log.Print("cannot connect to server")
		conn.Close()
		return
	}
	go handleStream(sConn, conn, ccfb)
	go handleStream(conn, sConn, scfb)
}

func main() {
	//read config
	c, _ := config.ReadDefault("config.ini")
	server, _ = c.String("client", "server")
	listen, _ = c.String("client", "listen")
	ck, _ := c.String("encrypto", "client-key")
	sk, _ := c.String("encrypto", "server-key")
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
