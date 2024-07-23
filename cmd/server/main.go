package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("-addr", ":3000", "http server address")

func serveHome() {
	http.Handle("/", http.FileServer(http.Dir("../app/dist")))
	log.Println("http server listening from", addr)
	log.Panic(http.ListenAndServe(*addr, nil))
}

func main() {
	flag.Parse()
	serveHome()
}
