package websocket

import (
	"errors"
	"github.com/google/uuid"
	"log"
)

type State int8

var GameState = struct {
	WAITING   State // Waiting for both players to join
	WHITETURN State // Game is ongoing and it is white's turn
	BLACKTURN State // Game is ongoing and it is black's turn
	WHITEWON  State // Game is finished and white won
	BLACKWON  State // Game is finished and black won
	TIE       State // Game is finisehd and resulted in a tie
}{
	WAITING:   0,
	WHITETURN: 1,
	BLACKTURN: 2,
	WHITEWON:  3,
	BLACKWON:  4,
	TIE:       5,
}

type Game struct {
	ID    string
	White *Player
	Black *Player
	State State
	Moves chan string // Moves send from both white and black
}

func newGame() *Game {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}

	newGame := &Game{
		ID:    id.String(),
		Moves: make(chan string),
		State: GameState.WAITING,
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
	return nil
}

func (g *Game) play() {
	// Wait for both players to join
	for g.State == GameState.WAITING {
		log.Println("WAITING")
		if g.White != nil && g.Black != nil {
			g.State = GameState.WHITETURN
		}
	}
	log.Println("Both palyers have joined", g.State)

	// Start playing the game
	for {
		select {
		case move := <-g.Moves:
			log.Println(move)
		}
		// check if both players are connected
		// connection timeout?

		// check for incoming move requests
		// if move, check if valid and update state and send to
		// both players

		// check for incoming draw or abort requests
		//

	}

}
