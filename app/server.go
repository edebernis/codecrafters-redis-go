package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("failed to bind to port 6379: ", err.Error())
		os.Exit(1)
	}

	_, err = l.Accept()
	if err != nil {
		fmt.Println("error accepting connection: ", err.Error())
		os.Exit(1)
	}
}
