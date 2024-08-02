package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type Player struct {
	ID     string
	Game   *Game
	Lobby  *Lobby
	Conn   *websocket.Conn
	Move   chan string
	InGame chan struct{}
}

var (
	messageType        = websocket.TextMessage
	matchMakingMaxWait = 60 // seconds
)

func NewPlayer(l *Lobby, conn *websocket.Conn) *Player {
	return &Player{
		Lobby:  l,
		Conn:   conn,
		Move:   make(chan string),
		InGame: make(chan struct{}),
	}
}

// write message from the Game to the websocket
// All writes to websocket MUST be in this function to avoid
// concurrent write errors
func (p *Player) write() {
	for {
		select {
		case message := <-p.Move:
			if err := p.Conn.WriteMessage(messageType, []byte(message)); err != nil {
				log.Printf("error: %v", err)
				return
			}
		}
	}
}

// read message from the websocket and notify the Game
// All reads from websocket MUST be in this function to avoid
// concurrent read errors
func (p *Player) read() {
	// From now on, every move must contain a valid playerID
	// Handle move requests
	for {
		_, message, err := p.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

        in := &Inbound{}
		err = json.Unmarshal(message, in)
		if err != nil {
			log.Printf("error: %v", err)
		}
		if p.ID != in.PlayerID {
			log.Printf("error: %v", err)
			continue // Soft handle invalid ids
		}
        switch in.Action {
        case MOVE:
		    p.Game.Moves <- in.Move
        default: 
            p.Conn.Close()
        }
	}
}
