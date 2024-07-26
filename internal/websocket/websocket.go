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

func ServeWebSocket(w http.ResponseWriter, r *http.Request, l *Lobby, gameID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

    player := NewPlayer(l, conn)
    l.PlayerPool <-player
    if player.Game != nil {
        // Player already in game, check gameIDs match
        // check client id's match
        // send redirect to the current url
        log.Printf("player %s already in game %s but requested %s", player.ID, player.Game.ID, gameID)
    }

	go player.write()
	go player.read()
}
