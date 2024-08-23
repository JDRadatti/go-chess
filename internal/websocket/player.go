package websocket

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

var (
	messageType              = websocket.TextMessage
	matchmakingMaxWait       = 180 // seconds
	writeWait                = 10 * time.Second
	pongWait                 = 60 * time.Second
	pingPeriod               = (pongWait * 9) / 10
	maxMessageSize     int64 = 512
)

type PlayerID string

func (pid PlayerID) validate() bool {
	if err := uuid.Validate(string(pid)); err != nil {
		return true
	}
	return false
}

type Player struct {
	id    PlayerID
	game  *Game
	lobby *Lobby
	conn  *websocket.Conn
	send  chan *Outbound
}

func NewPlayer(l *Lobby, c *websocket.Conn, g *Game) *Player {
	return &Player{
		lobby: l,
		conn:  c,
		game:  g,
		send:  make(chan *Outbound),
	}
}

func GeneratePlayerID() PlayerID {
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Printf("error %s", err)
	}
	return PlayerID(uuid.String())
}

// write message from the Game to the websocket
// All writes to websocket MUST be in this function to avoid
// concurrent write errors
func (p *Player) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		p.conn.Close()
	}()

	for {
		select {
		case out := <-p.send:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if message, ok := marshal(out); ok {
				if err := p.conn.WriteMessage(messageType, message); err != nil {
					log.Printf("error: %v", err)
					return
				}
			}
		case <-ticker.C:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// read message from the websocket and notify the Game
// All reads from websocket MUST be in this function to avoid
// concurrent read errors
func (p *Player) read() {
	defer func() {
		p.game.leave <- p
		p.conn.Close()
	}()

	p.conn.SetReadLimit(maxMessageSize)
	p.conn.SetReadDeadline(time.Now().Add(pongWait))
	p.conn.SetPongHandler(func(string) error { p.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		if p.game.state != playing {
			continue
		}

		_, message, err := p.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		if in, ok := unmarshal(message); ok {
			if p.id != in.PlayerID {
				log.Println("invalid player id", p.id, in.PlayerID)
				continue // Soft handle invalid ids
			}
			switch in.Action {
			case MOVE:
				p.game.move <- in
			case RESIGN:
				p.game.resign <- in
			case DRAW_DENY:
				fallthrough
			case DRAW_REQUEST:
				fallthrough
			case DRAW_ACCEPT:
				p.game.draw <- in
			case ABORT:
				p.game.abort <- in
			}
		}
	}
}
