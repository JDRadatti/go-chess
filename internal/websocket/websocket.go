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

func ServeWebSocket(w http.ResponseWriter, r *http.Request, l *Lobby, id string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

    game := createOrGet(l, id)
    player := newPlayer(game, l, conn)
    err = game.addPlayer(player)
    if err != nil {
        panic(err)
    }
    log.Println("GAME: ", game)

    go player.write()
    go player.read()
}
