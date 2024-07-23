package main

import (
	"os"
)

const (
	addressUDP  = "127.0.0.1:8080"
	addressHTTP = ":3000"
)

func main() {
	if len(os.Args) != 2 {
		panic(`at least one of the following flags required.
flags: 
    -udp   to start udp server 
    -http  to start http server
`)
	}

	switch os.Args[1] {
	case "-udp":
		startUDP()
	case "-http":
		startHTTP()
	default:
		panic("invalid flag")
	}
}
