package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

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
		// Send join fail
		defer conn.Close()
		joinSuccess := &Outbound{
			Action: JOIN_FAIL,
		}
		message, err := json.Marshal(joinSuccess)
		if err != nil {
			return
		}

		if err := conn.WriteMessage(messageType, []byte(message)); err != nil {
			log.Printf("error: %v", err)
		}
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
		return nil, false
	}

	in := &Inbound{}
	err = json.Unmarshal(message, in)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, false
	}
	if in.Action != JOIN {
		log.Printf("action must be join")
		return nil, false
	}

	player, ok := ws.Lobby.GetPlayer(in.PlayerID)
	if !ok {
		player = NewPlayer(ws.Lobby, conn, 0, 0)
		player.ID = GenerateID()
		if game, ok := ws.Lobby.GetGame(ws.GameID); ok {
			if err = game.addPlayer(player); err != nil {
				log.Println("game full")
				return nil, false
			}
			player.Time = game.WhiteTime
			player.Increment = game.Increment
		}
		if ok = ws.Lobby.AddPlayer(player); !ok {
			log.Println("invalid player")
			return nil, false
		}
	}

	if player.Game == nil || player.Game.ID != ws.GameID {
		log.Println("invalid game id")
		return nil, false
	}

	// Send join success
	joinSuccess := &Outbound{
		Action:    JOIN_SUCCESS,
		PlayerID:  player.ID,
		GameID:    player.Game.ID,
		Color:     player.Game.ColorFromPID(player.ID),
		WhiteTime: player.Game.WhiteTime,
		BlackTime: player.Game.BlackTime,
		Increment: player.Game.Increment,
		Turn:      player.Game.Board.Turn(),
	}
	message, err = json.Marshal(joinSuccess)
	if err != nil {
		log.Println("server error")
		return nil, false
	}

	if err := conn.WriteMessage(messageType, []byte(message)); err != nil {
		log.Printf("error: %v", err)
		return nil, false
	}

	player.Conn = conn
	return player, true
}
