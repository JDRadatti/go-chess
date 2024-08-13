package websocket

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
)

type State int8

type Game struct {
	ID       string
	White    *Player
	Black    *Player
	Moves    chan *Inbound // Moves requests sent from both white and black
	AllMoves []string
	Start    chan struct{}
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
		ID:       gameID.String()[:8],
		Moves:    make(chan *Inbound),
		Start:    make(chan struct{}),
		AllMoves: []string{},
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
		g.White.ID, g.Black.ID, g.AllMoves)
}

func (g *Game) play() {
	log.Println("waiting for game to start")
	<-g.Start // Wait for game to start
	log.Println("game started")
	for {
		select {
		case moveRequest := <-g.Moves:
			// check player id
			// check game id
			// check if valid player move
			// call gamelogic module to check for valid gamelogic
			pid := moveRequest.PlayerID
			switch pid {
			case g.White.ID:
				if len(g.AllMoves)%2 != 0 { // white moves on evens
					continue // skip moves out of order
				}
			case g.Black.ID:
				if len(g.AllMoves)%2 != 1 { // black moves on odds
					continue // skip moves out of order
				}
			default:
				log.Println("invalid playerID")
				return
			}

			// right now, just relay move to both players
			out := &Outbound{
				Action:   MOVE,
				Move:     moveRequest.Move,
				PlayerID: pid, // Player who made the move
				GameID:   g.ID,
			}
			g.AllMoves = append(g.AllMoves, out.Move)
			g.White.Move <- out
			g.Black.Move <- out
		default:
			continue
		}
	}
}
