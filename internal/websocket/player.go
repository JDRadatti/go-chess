package websocket

import (
	"encoding/json"
	"github.com/JDRadatti/reptile/internal/chess"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
)

type Player struct {
	ID        string
	Game      *Game
	Lobby     *Lobby
	Conn      *websocket.Conn
	Move      chan *Outbound
	InGame    chan struct{}
	Time      int
	Increment int
}

var (
	messageType        = websocket.TextMessage
	matchMakingMaxWait = 60 // seconds
)

func NewPlayer(l *Lobby, conn *websocket.Conn, time int, increment int) *Player {
	return &Player{
		Lobby:     l,
		Conn:      conn,
		Move:      make(chan *Outbound),
		InGame:    make(chan struct{}),
		Time:      time,
		Increment: increment,
	}
}

func GenerateID() string {
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Printf("error %s", err)
	}
	return uuid.String()
}

// write message from the Game to the websocket
// All writes to websocket MUST be in this function to avoid
// concurrent write errors
func (p *Player) write() {
	<-p.Game.Start // wait for game to start
	var fen string
	var turn chess.Player
	if p.Game.Board != nil {
		fen = string(p.Game.Board.FEN())
		turn = p.Game.Board.Turn()
	} else {
		fen = "RNBQKBNR/PPPPPPPP/8/8/8/8/pppppppp/rnbqkbnr"
	}
	message, err := json.Marshal(&Outbound{
		Action: GAME_START,
		GameID: p.Game.ID,
		FEN:    fen,
		Turn:   turn,
	})
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	if err := p.Conn.WriteMessage(messageType, []byte(message)); err != nil {
		log.Printf("error: %v", err)
		return
	}

	for {
		select {
		case out := <-p.Move:
			message, err := json.Marshal(out)
			if err != nil {
				log.Printf("error: %v", err)
				return
			}
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
	<-p.Game.Start // wait for game to start
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
			log.Println("invalid player id", p.ID, in.PlayerID)
			continue // Soft handle invalid ids
		}
		switch in.Action {
		case MOVE:
			p.Game.Moves <- in
		default:
			p.Conn.Close()
		}
	}
}
