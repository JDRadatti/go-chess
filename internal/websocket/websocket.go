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

	// Get player from lobby using PlayerID sent from client.
	// Client MUST send the playerID on connect
	_, message, err := conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("error: %v", err)
		}
		return
	}

	log.Println("message: ", string(message))

	if player, ok := l.GetPlayer(string(message)); ok {
		go player.write()
		go player.read()
    }
    //if player.Game != nil {
	//    // Player already in game, check gameIDs match
	//    // check client id's match
	//    // send redirect to the current url
	//    log.Printf("player %s already in game %s but requested %s", player.ID, player.Game.ID, gameID)
	//}
}
