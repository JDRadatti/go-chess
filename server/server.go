package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

const (
	address = "127.0.0.1:8080"
)

func main() {
	if len(os.Args) != 2 {
		panic("flags required: -s to start server")
	}

	switch os.Args[1] {
	case "-s":
		startServer()
	default:
		panic("invalid flag")
	}
}

func startServer() error {

	ln, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	log.Println("Listening from: ", address)
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go rogerThat(conn)
	}
}

// Sends the client a copy of what the client sent it
func rogerThat(conn net.Conn) {
	request := make([]byte, 1024)
	_, err := conn.Read(request)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Server recieved: %s", request)
}
