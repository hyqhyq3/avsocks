package handler

import (
	"crypto/cipher"
	"io"
	"net"
)

var iv = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15} //length must be 16

func HandleStream(t net.Conn, f net.Conn, s cipher.Stream) {
	in := make([]byte, 1024*4)
	out := make([]byte, 1024*4)
	for {
		n, err := f.Read(in)
		if err == nil || err == io.EOF {
			s.XORKeyStream(out, in[:n])
			_, err := t.Write(out[:n])
			if err != nil {
				f.Close()
				break
			}
		} else {
			t.Close()
			break
		}
	}
}
