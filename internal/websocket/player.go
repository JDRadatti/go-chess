package websocket

import (
	"encoding/json"
	"github.com/google/uuid"
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

// wait for match making to add player to a game
func (p *Player) WaitForGame() {
	<-p.InGame // Wait until matchmaking finishes

	payload := GameAccepted{
		GameID:   p.Game.ID,
		PlayerID: p.ID,
		Color:    p.Game.ColorFromPID(p.ID),
	}

	marshled, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error: %v", err)
	}

	if err := p.Conn.WriteMessage(websocket.TextMessage, marshled); err != nil {
		log.Printf("error: %v", err)
	}
}

// write message from the Game to the websocket
// All writes to websocket MUST be in this function to avoid
// concurrent write errors
func (p *Player) write() {
	p.WaitForGame()
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
func (p *Player) read(l *Lobby) {

	// First message must be the playerID
	_, message, err := p.Conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("error: %v", err)
		}
		return
	}

	var c ConnectRequest
	err = json.Unmarshal(message, &c)
	if err != nil {
		log.Printf("error: %v", err)
	}

	if err := uuid.Validate(c.PlayerID); err != nil {
		log.Printf("error: %v", err)
	}
	p.ID = c.PlayerID

	l.PlayerPool <- p // Request game from lobby
	<-p.InGame        // Wait for lobby to close, indicating game found

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

		var payload *MoveRequest = &MoveRequest{}
		log.Println(string(message))
		err = json.Unmarshal(message, payload)
		if err != nil {
			log.Printf("error: %v", err)
		}
		log.Println(*payload)
		if p.ID != payload.PlayerID {
			log.Printf("error: %v", err)
			continue // Soft handle invalid ids
		}
		p.Game.Moves <- payload
	}
}
