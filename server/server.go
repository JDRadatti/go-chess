package main

import (
	"fmt"
	"log"
	"net"
	"os"
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

	port := ":8080"
	ln, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	log.Println("Listening on port: ", port)
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
