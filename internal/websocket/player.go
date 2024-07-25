package websocket

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type Player struct {
	ID    uuid.UUID
	Game  *Game
	Lobby *Lobby
	Conn  *websocket.Conn
	Move  chan string
}

var (
	messageType        = websocket.TextMessage
	matchMakingMaxWait = 60 // seconds
)

func newPlayer(l *Lobby, c *websocket.Conn) *Player {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return &Player{
		ID:    id,
		Lobby: l,
		Conn:  c,
		Move:  make(chan string),
	}
}

// wait for match making to add player to a game
func (p *Player) waitForGame() {

	ticker := time.NewTicker(time.Duration(matchMakingMaxWait))
	defer func() { ticker.Stop() }()

	for {
		if p.Game != nil {
			if err := p.Conn.WriteMessage(messageType, []byte(p.Game.ID)); err != nil {
				log.Printf("error: %v", err)
			}
			return
		}

		select {
		case <-ticker.C:
			if err := p.Conn.WriteMessage(messageType, []byte("here")); err != nil {
				log.Printf("error: %v", err)
			}
			return
		}
	}

}

// write message from the Game to the websocket
func (p *Player) write() {
	p.waitForGame()
	log.Println("FOUND GAME: ", p.Game)
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
	for {
		_, message, err := p.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		log.Println("Connection SAYS: ", string(message))
	}
}
