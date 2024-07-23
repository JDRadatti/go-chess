package main

import (
	"fmt"
	"log"
	"net"
)

var count int = 0

func startUDP() {
	conn, err := net.ListenPacket("udp", addressUDP)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	log.Println("Listening from: ", addressUDP)
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
