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

var validTimes []int = []int{60, 180, 300, 600}
var validIncrements []int = []int{0, 1, 2, 10}

type Game struct {
	ID        string
	White     *Player
	Black     *Player
	Moves     chan *Inbound // Moves requests sent from both white and black
	Board     *chess.Board
	Start     chan struct{}
	WhiteTime int // number of seconds in the game for white
	BlackTime int
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
		WhiteTime: time,
		BlackTime: time,
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
		Action:    action,
		Move:      move,
		PlayerID:  pid, // Player who made the move
		GameID:    g.ID,
		FEN:       string(g.Board.FEN()),
		Message:   message,
		Turn:      g.Board.Turn(),
		WhiteTime: g.WhiteTime,
		BlackTime: g.BlackTime,
	}
}

func (g *Game) play() {
	<-g.Start // Wait for game to start

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if g.Board.Turn() == chess.WHITE && g.WhiteTime < 0 ||
				g.Board.Turn() == chess.BLACK && g.BlackTime < 0 {
				// game over
				out := g.Out(GAME_OVER, "", "", "on time")
				g.Lobby.Clean(g.ID, g.White.ID, g.Black.ID)
				g.White.Move <- out
				g.Black.Move <- out
				return
			} else {
				out := g.Out(TIME_UPDATE, "", "", "")
				g.White.Move <- out
				g.Black.Move <- out
			}
			if g.Board.Turn() == chess.WHITE {
				g.WhiteTime -= 1
			} else {
				g.BlackTime -= 1
			}
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
					var out *Outbound
					if status, over := g.Board.GameOver(); over {
						out = g.Out(GAME_OVER, move, pid, status)
						g.Lobby.Clean(g.ID, g.White.ID, g.Black.ID)
					} else {
						if g.Board.Turn() == chess.WHITE {
							g.BlackTime += g.Increment
						} else {
							g.WhiteTime += g.Increment
						}
						out = g.Out(MOVE, move, pid, "")
					}
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
