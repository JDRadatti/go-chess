package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		panic("flags required: -c to connect to server")
	}

	switch os.Args[1] {
	case "-c":
		conn := connectToServer()
		for i := 0; i < 10; i++ {
			doSomething(conn, i)
		}
	default:
		panic("invalid flag")
	}
}

func connectToServer() net.Conn {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	return conn
}

func doSomething(conn net.Conn, i int) {
	fmt.Fprintf(conn, "do something %d\n", i)
	response := make([]byte, 1024)
	_, err := conn.Read(response)
	if err != nil {
		panic(err)
	}

	fmt.Println("response from server: ", response)
}
