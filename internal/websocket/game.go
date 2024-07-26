package websocket

import (
    "github.com/google/uuid"
	"errors"
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

	gameID, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}

	newGame := &Game{
        ID:    gameID.String()[:8],
		Moves: make(chan string),
		State: GameState.WAITING,
	}
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
}
