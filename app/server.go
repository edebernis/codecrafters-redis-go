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

type handler struct {
	conn net.Conn
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	timeoutDuration := 5 * time.Second
	bufReader := bufio.NewReader(conn)
	handler := &handler{conn: conn}

	for {
		conn.SetReadDeadline(time.Now().Add(timeoutDuration))

		bytes, err := bufReader.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatal(err)
		}

		if err := handler.handleCommand(bytes); err != nil {
			fmt.Println(err)
		}
	}
}

func (h *handler) handleCommand(cmd []byte) error {
	cmd = strings.Trim(string(cmd), "\r\n"))

	if strings.ToLower(cmd) == "ping" {
		if _, err := h.conn.Write([]byte("+PONG\r\n")); err != nil {
			return err
		}
	}

	return nil
}
