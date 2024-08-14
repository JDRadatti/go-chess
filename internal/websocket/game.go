package websocket

import (
	"errors"
	"fmt"
	"github.com/JDRadatti/reptile/internal/chess"
	"github.com/google/uuid"
	"log"
	"slices"
	"time"
)

type State int8

var validTimes []int = []int{60, 180, 600}
var validIncrements []int = []int{0, 1, 5, 10}

type Game struct {
	ID        string
	White     *Player
	Black     *Player
	Moves     chan *Inbound // Moves requests sent from both white and black
	Board     *chess.Board
	Start     chan struct{}
	Time      int // number of seconds in the game
	Increment int // number of seconds to add when player moves
	Lobby     *Lobby
}

func newGame(l *Lobby, time int, increment int) *Game {

	gameID, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}

	if !slices.Contains(validTimes, time) {
		time = 180
	}

	if !slices.Contains(validIncrements, increment) {
		increment = 0
	}

	board := chess.NewBoardClassic()
	newGame := &Game{
		ID:        gameID.String()[:8],
		Moves:     make(chan *Inbound),
		Start:     make(chan struct{}),
		Board:     &board,
		Time:      time,
		Increment: increment,
		Lobby:     l,
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

func (g *Game) Out(action Action, move string, pid string, message string) *Outbound {
	return &Outbound{
		Action:   action,
		Move:     move,
		PlayerID: pid, // Player who made the move
		GameID:   g.ID,
		FEN:      string(g.Board.FEN()),
	}
}

func (g *Game) play() {
	<-g.Start // Wait for game to start

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if g.Time < 0 {
				// game over
				out := g.Out(GAME_OVER, "", "", "on time")
				g.Lobby.Clean(g.ID, g.White.ID, g.Black.ID)
				g.White.Move <- out
				g.Black.Move <- out
				return
			}
			g.Time -= 1 + g.Increment
		case moveRequest := <-g.Moves:

			pid := moveRequest.PlayerID
			if !g.ValidPID(pid) {
				out := g.Out(INVALID_MOVE, "", pid, "")
				if pid == g.White.ID {
					g.White.Move <- out
				} else {
					g.Black.Move <- out
				}
			} else {

				// try move
				move, valid := g.Board.Move(moveRequest.Move)
				if valid {
					out := g.Out(MOVE, move, pid, "")
					g.White.Move <- out
					g.Black.Move <- out
				} else {
					out := g.Out(INVALID_MOVE, "", pid, "")
					if pid == g.White.ID {
						g.White.Move <- out
					} else {
						g.Black.Move <- out
					}
				}
			}

		default:
			continue
		}
	}
}
