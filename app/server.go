package main

import (
	"errors"
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

			if err := c.SetDeadline(time.Now().Add(time.Second * 5)); err != nil {
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

				if _, err := c.Write([]byte("+PONG\r\n")); err != nil {
					log.Fatal(err)
				}
			}
		}(conn)
	}
}
