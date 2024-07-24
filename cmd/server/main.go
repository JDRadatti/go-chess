package main

import (
	"flag"
	"log"
	"net/http"
    "github.com/JDRadatti/reptile/internal/websocket"
)

var addr = flag.String("addr", ":3000", "http server address")

func serveHome() {
    router := http.NewServeMux()
	router.HandleFunc("GET /game/{id}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("game id", r)
        idString := r.PathValue("id")
        log.Println(idString)
        websocket.ServeWebSocket(w, r)
	})
	router.Handle("/", http.FileServer(http.Dir("app/dist")))
	log.Println("http server listening from", *addr)
	log.Panic(http.ListenAndServe(*addr, router))
}

func main() {
	flag.Parse()
	serveHome()
}
