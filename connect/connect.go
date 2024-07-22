package main

import (
	"log"
	"net"
	"os"
)

const (
	address = "127.0.0.1:8080"
)

func main() {
	if len(os.Args) != 2 {
		panic("flags required: -c to connect to server")
	}

	switch os.Args[1] {
	case "-c":
		connectToServer()
	default:
		panic("invalid flag")
	}
}

func connectToServer() {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		panic(err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for i := 0; i < 10; i++ {
		send := []byte("[sending from client]")
		_, err := conn.Write(send)
		if err != nil {
			panic(err)
		}
		log.Println("Sent: ", string(send))

		buffer := make([]byte, 1024)
		_, err = conn.Read(buffer)
		if err != nil {
			panic(err)
		}

		log.Println("Recieved: ", string(buffer))
	}
}
