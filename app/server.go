package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"
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

			if err := c.SetDeadline(3 * time.Second); err != nil {
				log.Fatal(err)
			}

			for {
				var b []byte
				if _, err := c.Read(b); err != nil {
					if errors.Is(err, os.ErrDeadlineExceeded) {
						break
					}
					log.Fatal(err)
				}
				fmt.Println(string(b))

				if _, err := c.Write([]byte("+PONG\r\n")); err != nil {
					log.Fatal(err)
				}
			}
		}(conn)
	}
}
