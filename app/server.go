package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(c net.Conn) {
			defer c.Close()

			var b bytes.Buffer
			if _, err := io.Copy(&b, c); err != nil {
				if !errors.Is(err, io.EOF) {
					log.Fatal(err)
				}
			}
			fmt.Println(b.String())

			if _, err := c.Write([]byte("+PONG\r\n")); err != nil {
				log.Fatal(err)
			}
		}(conn)
	}
}
