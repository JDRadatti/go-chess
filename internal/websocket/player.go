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
	playerID, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}

	return &Player{
		ID:     playerID.String()[:8],
		Lobby:  l,
		Conn:   conn,
		Move:   make(chan string),
		InGame: make(chan struct{}),
	}
}

func (p *Player) AddConn(conn *websocket.Conn) {
	p.Conn = conn
}

// wait for match making to add player to a game
func (p *Player) WaitForGame() {
	<-p.InGame // Wait until matchmaking finishes

	payload := GameAccepted{
		PlayerID: p.ID,
		GameID:   p.Game.ID,
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
func (p *Player) read() {
	<-p.InGame
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
        p.Game.Moves <- payload
	}
}
