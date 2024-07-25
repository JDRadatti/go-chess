package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ServeWebSocket(w http.ResponseWriter, r *http.Request, l *Lobby) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	player := newPlayer(l, conn)
    l.PlayerPool <- player

	go player.write()
	go player.read()
}
