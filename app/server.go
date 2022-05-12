package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
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

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer c.Close()

	timeoutDuration := 5 * time.Second
	bufReader := bufio.NewReader(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(timeoutDuration))

		bytes, err := bufReader.ReadBytes('\n')
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", bytes)

		if _, err := c.Write([]byte("+PONG\r\n")); err != nil {
			log.Fatal(err)
		}
	}
}
