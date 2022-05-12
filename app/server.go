package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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
	defer conn.Close()

	timeoutDuration := 5 * time.Second
	bufReader := bufio.NewReader(conn)

	var cmd []string
	for {
		conn.SetReadDeadline(time.Now().Add(timeoutDuration))

		bytes, err := bufReader.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatal(err)
		}

		input := strings.Trim(string(bytes), "\r\n")
		cmd = append(cmd, input)
	}

	if err := handleCommand(conn, cmd); err != nil {
		fmt.Println(err)
	}
}

func handleCommand(conn net.Conn, cmd []string) error {
	if _, err := conn.Write([]byte("+PONG\r\n")); err != nil {
		return err
	}

	return nil
}
