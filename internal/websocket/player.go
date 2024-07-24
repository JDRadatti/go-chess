package websocket

import (
	"github.com/gorilla/websocket"
	"log"
)

type Player struct {
	Game  *Game
	Lobby *Lobby
	Conn  *websocket.Conn
    Move chan string
}

func newPlayer(g *Game, l *Lobby, c *websocket.Conn) *Player {
	return &Player{
		Game:  g,
		Lobby: l,
		Conn:  c,
	}
}

var (
    messageType = 0
)
// write message from the Game to the websocket
func (p *Player) write() {
	for {
		select {
		case message := <-p.Move:
			if err := p.Conn.WriteMessage(messageType, []byte(message)); err != nil {
				log.Println(err)
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
		log.Println("Connectikon SAYS: ", string(message))
		p.Game.Moves <- string(message)
	}
}
