package websocket

import (
	"encoding/json"
	"github.com/JDRadatti/reptile/internal/chess"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"time"
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
	messageType              = websocket.TextMessage
	matchmakingMaxWait       = 180 // seconds
	writeWait                = 10 * time.Second
	pongWait                 = 60 * time.Second
	pingPeriod               = (pongWait * 9) / 10
	maxMessageSize     int64 = 512
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

func (p *Player) LeaveGame() {
	p.Game = nil
	p.Time = 0
	p.Increment = 0
	p.InGame = make(chan struct{})
	p.Move = make(chan *Outbound)
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
	timer := time.NewTimer(time.Duration(matchmakingMaxWait) * time.Second)
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		timer.Stop()
		p.Conn.Close()
	}()
	log.Println("WAITING.. game requested")
	select {
	case <-timer.C:
		log.Println("Timer finished")
		joinSuccess := &Outbound{
			Action:  MATCHMAKING_ERROR,
			Message: "Matchmaking took too long.",
		}
		message, err := json.Marshal(joinSuccess)
		if err != nil {
			return
		}

		if err := p.Conn.WriteMessage(messageType, []byte(message)); err != nil {
			log.Printf("error: %v", err)
		}
		// remove from game
		return
	case <-p.Game.Start: // wait for game to start
		log.Println("game started")
	}
	log.Println("GAME STARTED")

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
			p.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			message, err := json.Marshal(out)
			if err != nil {
				log.Printf("error: %v", err)
				return
			}
			if err := p.Conn.WriteMessage(messageType, []byte(message)); err != nil {
				log.Printf("error: %v", err)
				return
			}

		case <-ticker.C:
			p.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
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
		p.Conn.Close()
		log.Println("READ CLOSE")
		p.LeaveGame()
	}()

	p.Conn.SetReadLimit(maxMessageSize)
	p.Conn.SetReadDeadline(time.Now().Add(pongWait))
	p.Conn.SetPongHandler(func(string) error { p.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

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
			return
		}
	}
}
