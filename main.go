package main

import (
	. "./handler"
	"crypto/aes"
	"crypto/cipher"
	"flag"
	"github.com/kless/goconfig/config"
	"log"
	"net"
)

var mode, server, listen, ck, sk string
var clientCipher, serverCipher cipher.Block
var handler Handler

func loadConfig() {
	//read config
	c, _ := config.ReadDefault("config.ini")
	mode, _ = c.String("main", "mode")
	server, _ = c.String("client", "server")
	listen, _ = c.String(mode, "listen")
	ck, _ = c.String("encrypto", "client-key")
	sk, _ = c.String("encrypto", "server-key")
}

func loadFlags() {
	flag.StringVar(&mode, "mode", mode, "server or client")
	flag.StringVar(&server, "server", server, "the remote server")
	flag.StringVar(&listen, "listen", listen, "the ip and port to bind")
	flag.StringVar(&ck, "client-key", ck, "the client key")
	flag.StringVar(&sk, "server-key", sk, "the server key")
	flag.Parse()
	log.Printf("mode:%s", mode)
	log.Printf("server:%s", server)
	log.Printf("listen:%s", listen)
	log.Printf("client-key:%s", ck)
	log.Printf("server-key:%s", sk)
}

func main() {
	loadConfig()
	loadFlags()

	clientCipher, _ = aes.NewCipher([]byte(ck))
	serverCipher, _ = aes.NewCipher([]byte(sk))

	switch mode {
	case "server":
		handler = &Server{
			ClientCipher: clientCipher,
			ServerCipher: serverCipher,
		}
	case "client":
		handler = &Client{
			Server:       server,
			ClientCipher: clientCipher,
			ServerCipher: serverCipher,
		}
	}

	lsn, e := net.Listen("tcp", listen)
	if e != nil {
		log.Fatal(e)
	}
	for {
		conn, e := lsn.Accept()
		if e != nil {
			log.Print(e)
		}
		go handler.Handle(conn)
	}
}
