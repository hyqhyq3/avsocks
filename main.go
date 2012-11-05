package main

import (
	. "./handler"
	"crypto/aes"
	"crypto/cipher"
	"flag"
	"./github.com/kless/goconfig/config"
	"log"
	"net"
)

var mode, server, client_listen, server_listen, ck, sk, http_proxy string
var clientCipher, serverCipher cipher.Block
var handler Handler

func loadConfig() {
	//read config
	c, _ := config.ReadDefault("config.ini")
	mode, _ = c.String("main", "mode")
	server, _ = c.String("client", "server")
	client_listen, _ = c.String("client", "listen")
	server_listen, _ = c.String("server", "listen")
	ck, _ = c.String("encrypto", "client-key")
	sk, _ = c.String("encrypto", "server-key")
	http_proxy, _ = c.String("client", "http-proxy")
}

func loadFlags() {
	flag.StringVar(&mode, "mode", mode, "server or client")
	flag.StringVar(&server, "server", server, "the remote server")
	flag.StringVar(&client_listen, "client-listen", client_listen, "the ip and port of client to bind")
	flag.StringVar(&server_listen, "server-listen", server_listen, "the ip and port of server to bind")
	flag.StringVar(&ck, "client-key", ck, "the client key")
	flag.StringVar(&sk, "server-key", sk, "the server key")
	flag.StringVar(&http_proxy, "http-proxy", http_proxy, "the ip and port of http proxy server ")
	flag.Parse()
	log.Printf("mode:%s", mode)
	if mode == "client" {
		log.Printf("server:%s", server)
		log.Printf("client-listen:%s", client_listen)
		if http_proxy != "" {
			log.Printf("http-proxy:%s", http_proxy)
		}
	} else {
		log.Printf("server-listen:%s", server_listen)
	}

	log.Printf("client-key:%s", ck)
	log.Printf("server-key:%s", sk)
}

func main() {
	loadConfig()
	loadFlags()

	clientCipher, _ = aes.NewCipher([]byte(ck))
	serverCipher, _ = aes.NewCipher([]byte(sk))

	var listen string
	switch mode {
	case "server":
		handler = &Server{
			ClientCipher: clientCipher,
			ServerCipher: serverCipher,
		}
		listen = server_listen
	case "client":
		handler = &Client{
			Server:       server,
			ClientCipher: clientCipher,
			ServerCipher: serverCipher,
		}
		listen = client_listen
	}

	lsn, e := net.Listen("tcp", listen)
	if e != nil {
		log.Fatal(e)
	}
	defer lsn.Close()

	if mode == "client" && http_proxy != "" {
		go startHTTPProxyServer(http_proxy)
	}

	for {
		conn, e := lsn.Accept()
		if e != nil {
			log.Print(e)
			continue
		}
		go handler.Handle(conn)
	}
}
