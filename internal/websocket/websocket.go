package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
    "time"
)

type Action string

const (
	JOIN Action = "join"
	MOVE Action = "move"
)

const (
	JOIN_SUCCESS = "join success"
)

type Inbound struct {
	Action   Action
	Move     string
	PlayerID string
	GameID   string
	Time     time.Duration
	Color    int
}

type Outbound struct {
	Action   Action
	Move     string
	PlayerID string
	GameID   string
	Time     time.Duration
	Color    int
	Message  string
}

type WSHandler struct {
	Lobby  *Lobby
	GameID string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (ws *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	if player, ok := ws.handshake(conn); ok {
		go player.write()
		go player.read()
	} else {
		log.Println("unregistered player. must make post request to /play before joining websocket")
	}
}

func (ws *WSHandler) handshake(conn *websocket.Conn) (*Player, bool) {

	// Get player from lobby using PlayerID sent from client.
	// Client MUST send the playerID on connect
	_, message, err := conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("error: %v", err)
		}
		return
	}

    in := &Inbound{}
	err = json.Unmarshal(message, in)
	if err != nil {
		log.Printf("error: %v", err)
        conn.Close()
        return
	}
    if in.Action != JOIN {
        log.Printf("action must be join")
        conn.Close()
        return
    }

	if player, ok := l.GetPlayer(in.PlayerID); ok {
        player.Conn = conn
		go player.write()
		go player.read()
	} else {
		log.Println("unregistered player. must make post request to /play before joining websocket")
	}
}
