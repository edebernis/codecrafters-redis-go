package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
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

	cmd      []string
	bulkSize *int64
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	bufReader := bufio.NewReader(conn)
	handler := &handler{conn: conn}

	for {
		input, err := bufReader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatal(err)
		}

		if err := handler.handleInput(input); err != nil {
			fmt.Println(err)
		}
	}
}

func (h *handler) handleInput(input string) error {
	input = strings.TrimRight(input, "\r\n")

	// First input of the command
	if h.cmd == nil {
		if input[0] != '*' {
			return errors.New("client must send a RESP array")
		}
		len, err := strconv.ParseInt(input[1:], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid array size %s: %w", input[1:], err)
		}
		h.cmd = make([]string, 0, int(len))
		return nil
	}

	if h.bulkSize == nil {
		if input[0] != '$' {
			return errors.New("RESP array must be contained of bulk strings only")
		}
		len, err := strconv.ParseInt(input[1:], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid array size %s: %w", input[1:], err)
		}
		h.bulkSize = &len
		return nil
	}

	if int(*h.bulkSize) != len(input) {
		return fmt.Errorf("bulk string size is incorrect. Bulk size %d != input size %d", *h.bulkSize, len(input))
	}

	h.cmd = append(h.cmd, input)
	h.bulkSize = nil

	// Array items remaining
	if len(h.cmd) != cap(h.cmd) {
		return nil
	}

	return h.doCommand()
}

func (h *handler) doCommand() error {
	defer func() { h.cmd = nil }()

	switch strings.ToLower(h.cmd[0]) {
	case "ping":
		return h.write(
			newSimpleString("PONG"),
		)
	case "echo":
		return h.write(
			newBulkString(h.cmd[1]),
		)
	default:
		return fmt.Errorf("unknown command: %s", h.cmd)
	}
}

func (h *handler) write(cmd string) error {
	_, err := h.conn.Write([]byte{cmd + "\r\n"})
	return err
}

func newSimpleString(str string) {
	return "+" + str
}

func newBulkString(str string) {
	return fmt.Sprintf("$%d\r\n%s", len(str), str)
}
