package main

import (
	"flag"
	"log"
	"net/http"
    "github.com/JDRadatti/reptile/internal/websocket"
)

var addr = flag.String("addr", ":3000", "http server address")

func serveHome(lobby *websocket.Lobby) {
    router := http.NewServeMux()
	router.HandleFunc("GET /game/{id}", func(w http.ResponseWriter, r *http.Request) {
        if _, ok := r.Header["Upgrade"]; ok {
            idString := r.PathValue("id")
            websocket.ServeWebSocket(w, r, lobby, idString)
        } else {
            http.ServeFile(w, r, "app/dist/index.html")
        }

	})
	router.Handle("/", http.FileServer(http.Dir("app/dist")))
	log.Println("http server listening from", *addr)
	log.Panic(http.ListenAndServe(*addr, router))
}

func main() {
	flag.Parse()
    lobby := websocket.NewLobby()
	serveHome(lobby)
}
