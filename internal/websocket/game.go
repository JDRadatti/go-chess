package websocket

import (
	"errors"
	"fmt"
	"github.com/JDRadatti/reptile/internal/chess"
	"github.com/google/uuid"
	"log"
)

type State int8

type Game struct {
	ID    string
	White *Player
	Black *Player
	Moves chan *Inbound // Moves requests sent from both white and black
	Board chess.Board
	Start chan struct{}
}

const (
	WHITE int = iota
	BLACK
	ERROR
)

func newGame() *Game {

	gameID, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}

	newGame := &Game{
		ID:    gameID.String()[:8],
		Moves: make(chan *Inbound),
		Start: make(chan struct{}),
		Board: chess.NewBoardClassic(),
	}
	go newGame.play()
	return newGame
}

func (g *Game) addPlayer(p *Player) error {
	if g.White == nil {
		g.White = p
	} else if g.Black == nil {
		g.Black = p
	} else {
		return errors.New("game full")
	}
	if g.White != nil && g.Black != nil {
		close(g.Start) // Tell game to start
	}
	p.Game = g
	return nil
}

func (g *Game) ColorFromPID(pid string) int {
	if g.White != nil && g.White.ID == pid {
		return WHITE
	} else if g.Black != nil && g.Black.ID == pid {
		return BLACK
	} else {
		log.Printf("Player ID not found %s", pid)
		return ERROR
	}
}

func (g *Game) String() string {
	return fmt.Sprintf("white: %s, black %s \n moves: %v",
		g.White.ID, g.Black.ID)
}

func (g *Game) play() {
	<-g.Start // Wait for game to start
	for {
		select {
		case moveRequest := <-g.Moves:

			pid := moveRequest.PlayerID
			if pid != g.White.ID && pid != g.Black.ID {
				log.Println("invalid move: invalid playerID")
				break
			}

			// try move
			move, valid := g.Board.Move(moveRequest.Move)
			if !valid {
				break
			}

			out := &Outbound{
				Action:   MOVE,
				Move:     move,
				PlayerID: pid, // Player who made the move
				GameID:   g.ID,
				FEN:      string(g.Board.FEN()),
			}

			// relay move to both players
			g.White.Move <- out
			g.Black.Move <- out
		default:
			continue
		}
	}
}
