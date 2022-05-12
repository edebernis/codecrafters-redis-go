package main

import (
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
			input, err := io.ReadAll(c)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(input))

			if _, err := c.Write([]byte("+PONG\r\n")); err != nil {
				log.Fatal(err)
			}

			c.Close()
		}(conn)
	}
}
