package handler

import (
	"crypto/cipher"
	"log"
	"net"
)

var iv = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15} //length must be 16

func HandleStream(t net.Conn, f net.Conn, s cipher.Stream) {
	in := make([]byte, 1024*4)
	out := make([]byte, 1024*4)
	for {
		n, err := f.Read(in)
		if n > 0 {
			s.XORKeyStream(out, in[:n])
			_, err2 := t.Write(out[:n])
			if err != nil || err2 != nil {
				return
			}
		} else {
			return
		}

	}
}

func D(v ...interface{}) {
	if false {
		log.Print(v...)
	}
}
