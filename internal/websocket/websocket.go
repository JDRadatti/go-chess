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
    // check for client id
    // if player id valid, send current game or add to pool
    // if no player id, create new player and add to pool

    player := NewPlayer(l, conn)
    //if player.Game != nil {
    //    // Player already in game, check gameIDs match
    //    // check client id's match
    //    // send redirect to the current url
    //    log.Printf("player %s already in game %s but requested %s", player.ID, player.Game.ID, gameID)
    //}
	go player.write()
	go player.read(l)
}
