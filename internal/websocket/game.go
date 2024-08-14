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
	Board *chess.Board
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

	board := chess.NewBoardClassic()
	newGame := &Game{
		ID:    gameID.String()[:8],
		Moves: make(chan *Inbound),
		Start: make(chan struct{}),
		Board: &board,
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

func (g *Game) ColorFromPID(pid string) chess.Player {
	if g.White != nil && g.White.ID == pid {
		return chess.WHITE
	} else if g.Black != nil && g.Black.ID == pid {
		return chess.BLACK
	} else {
		log.Printf("Player ID not found %s", pid)
		return chess.INVALID_PLAYER
	}
}

func (g *Game) ValidPID(pid string) bool {
	return (pid == g.White.ID || pid == g.Black.ID) && g.ColorFromPID(pid) == g.Board.Turn()
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
			if !g.ValidPID(pid) {
				goto SEND_ERROR
			} else {

				// try move
				move, valid := g.Board.Move(moveRequest.Move)
				if valid {
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
				} else {
					log.Println("INVALID MOVE")
					goto SEND_ERROR
				}
			}
		SEND_ERROR:
			out := &Outbound{
				Action:   INVALID_MOVE,
				Move:     "",
				PlayerID: pid, // Player who made the move
				GameID:   g.ID,
				FEN:      string(g.Board.FEN()),
			}
			if pid == g.White.ID {
				g.White.Move <- out
			} else {
				g.Black.Move <- out
			}

		default:
			continue
		}
	}
}
