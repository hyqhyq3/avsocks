package main

import (
	"bufio"
	"errors"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

func int2nl(i int) []byte {
	b := make([]byte, 2)
	b[0] = byte((i & 0xffff) >> 8)
	b[1] = byte(i & 0xff)
	return b
}

func handshakeWithSocksServer(conn net.Conn) bool {
	conn.Write([]byte("\x05\x01\x00"))
	b := make([]byte, 2)
	conn.Read(b)
	return b[0] == 5 && b[1] == 0
}

func connectSocksServer(socks net.Conn, domain string, port int) (err error) {
	handshakeWithSocksServer(socks)
	w := bufio.NewWriter(socks)
	//send connect command
	w.Write([]byte("\x05\x01\x00\x03"))
	w.WriteByte(byte(len(domain)))
	w.Write([]byte(domain))
	w.Write(int2nl(port))
	w.Flush()

	b := make([]byte, 256)
	socks.Read(b[:4])
	if b[0] != 5 || b[1] != 0 {
		return errors.New("cannot connect to server")
	}
	switch b[3] {
	case 1: //ipv4
		socks.Read(b[:6]) //discard the result
	case 3:
		socks.Read(b[:1])
		socks.Read(b[:b[0]+2])
	case 4:
		socks.Read(b[:16])
	}
	return
}

func HandleHTTP(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	req, err := http.ReadRequest(r)
	if err != nil || req == nil {
		log.Print("Invalid HTTP Request: ", err)
		return
	}
	var domain string
	var port int
	slices := strings.Split(req.Host, ":")
	if len(slices) == 1 {
		port = 80
	} else {
		port, _ = strconv.Atoi(slices[1])
	}
	domain = slices[0]
	socks, err := net.Dial("tcp", "localhost:1080")
	if err != nil {
		return
	}
	defer socks.Close()
	err = connectSocksServer(socks, domain, port)
	if err != nil {
		log.Print("cannot connect to server")
		return
	}
	req.Write(socks)
	b := make([]byte, 1024*32)
	var re, we error
	var n int
	go func() {
		b := make([]byte, 1024*4)
		n, re = conn.Read(b)
		if n > 0 {
			_, we = socks.Write(b[:n])
		}
		if re != nil || we != nil {
			return
		}
	}()
	for {
		n, re = socks.Read(b)
		if n > 0 {
			_, we = conn.Write(b[:n])
		}
		if re != nil || we != nil {
			return
		}
	}
}

func startHTTPProxyServer(addr string) {
	lsn, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal(e)
	}
	for {
		conn, e := lsn.Accept()
		if e != nil {
			log.Print(e)
			continue
		}
		go HandleHTTP(conn)
	}
}
