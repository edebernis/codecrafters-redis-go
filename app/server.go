package main

import (
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
			if _, err := c.Write([]byte("+PONG\r\n")); err != nil {
				log.Fatal(err)
			}

			c.Close()
		}(conn)
	}
}
