package main

import (
	"flag"
	"github.com/JDRadatti/reptile/internal/api"
	"github.com/JDRadatti/reptile/internal/websocket"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":3000", "http server address")

func serveHome(lobby *websocket.Lobby) {
	router := http.NewServeMux()

	router.HandleFunc("POST /play", func(w http.ResponseWriter, r *http.Request) {
		api.HandlePlay(w, r, lobby)
	})
	router.HandleFunc("GET /play/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "app/dist/index.html")
	})

	router.HandleFunc("GET /lobby", func(w http.ResponseWriter, r *http.Request) {
		api.HandleLobby(w, lobby)
	})

	router.HandleFunc("GET /game/{id}", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Header["Upgrade"]; ok {
			idString := r.PathValue("id")
			wsHandler := &websocket.WSHandler{
				Lobby:  lobby,
				GameID: websocket.GameID(idString),
			}
			wsHandler.ServeHTTP(w, r)
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
