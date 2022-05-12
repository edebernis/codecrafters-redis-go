package main

import (
	"fmt"
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

			var b []byte
			if _, err := c.Read(b); err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(b))

			if _, err := c.Write([]byte("+PONG\r\n")); err != nil {
				log.Fatal(err)
			}
		}(conn)
	}
}
