package websocket

import (
	"errors"
	"github.com/google/uuid"
	"log"
)

type State int8

type Game struct {
	ID    string
	White *Player
	Black *Player
	Moves chan *MoveRequest // Moves requests sent from both white and black
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
		Moves: make(chan *MoveRequest),
		Start: make(chan struct{}),
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
			log.Println(moveRequest)
		}
	}
}
