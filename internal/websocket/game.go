package websocket

import (
    "github.com/google/uuid"
	"errors"
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
		g.Black = p
	} else if g.Black == nil {
		g.White = p
	} else {
		return errors.New("game full")
	}
    if g.White != nil && g.Black != nil {
        close(g.Start) // Tell game to start
    }
	return nil
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
