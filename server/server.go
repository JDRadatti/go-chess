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

var count int = 0

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
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	log.Println("Listening from: ", address)
	for {
		buffer := make([]byte, 1024)
		_, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			panic(err)
		}
		log.Println("server recieved", string(buffer))
		go rogerThat(conn, addr)
	}
}

// Sends the client a copy of what the client sent it
func rogerThat(conn net.PacketConn, addr net.Addr) {
	response := fmt.Sprintf("Num of responses %d", count)
    log.Println("server sent: ", string(response))
	conn.WriteTo([]byte(response), addr)
	count++
}
