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

var (
	validTimes      = []int{60, 180, 300, 600}
	validIncrements = []int{0, 1, 2, 10}
)

const (
	defaultTime      = 300
	defaultIncrement = 0
	whiteIndex       = 0
	blackIndex       = 1
)

type GameState int8

const (
	waiting GameState = iota
	playing
	over
)

type GameID string

type Game struct {
	id            GameID
	players       [2]*Player
	playerIDs     [2]PlayerID
	timeRemaining [2]int
	increment     int // number of seconds to add when player moves
	join          chan *Player
	leave         chan *Player
	move          chan *Inbound // Moves requests sent from both white and black
	board         *chess.Board
	lobby         *Lobby
	state         GameState
}

func NewGame(l *Lobby, time int, increment int) *Game {

	if !slices.Contains(validTimes, time) {
		time = defaultTime
	}

	if !slices.Contains(validIncrements, increment) {
		increment = defaultIncrement
	}

	board := chess.NewBoardClassic()
	newGame := &Game{
		id:            generateGameID(),
		move:          make(chan *Inbound),
		board:         &board,
		timeRemaining: [2]int{time, time},
		increment:     increment,
		lobby:         l,
		state:         waiting,
	}
	go newGame.play()
	return newGame
}

func (g *Game) clean() {
	g.lobby.Clean(g.id, g.playerIDs[whiteIndex], g.playerIDs[blackIndex])
}

func (g *Game) playerFromID(playerID PlayerID) (*Player, int, bool) {
	if g.playerIDs[whiteIndex] == playerID {
		return g.players[whiteIndex], whiteIndex, true
	} else if g.playerIDs[blackIndex] == playerID {
		return g.players[whiteIndex], blackIndex, true
	} else {
		return nil, -1, false
	}
}

func (g *Game) playerIndex(player *Player) (int, bool) {
	if g.playerIDs[whiteIndex] == player.ID {
		return whiteIndex, true
	} else if g.playerIDs[blackIndex] == player.ID {
		return blackIndex, true
	} else {
		return -1, false
	}
}

func generateGameID() GameID {
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Printf("error %s", err)
	}
	return GameID(uuid.String()[:8])
}

func (g *Game) out(action Action, move string, pid PlayerID, message string) *Outbound {
	return &Outbound{
		Action:    action,
		Move:      move,
		PlayerID:  pid, // Player who made the move
		GameID:    g.id,
		FEN:       string(g.board.FEN()),
		Message:   message,
		Turn:      g.board.Turn(),
		WhiteTime: g.timeRemaining[whiteIndex],
		BlackTime: g.timeRemaining[blackIndex],
	}
}

func (g *Game) play() {
	ticker := time.NewTicker(time.Second)
	defer func() {
		ticker.Stop()
		g.clean()
	}()

	for {
		select {
		case player := <-g.join:
			if index, ok := g.playerIndex(player); ok {
				g.players[index] = player
			}
		case player := <-g.leave:
			if index, ok := g.playerIndex(player); ok {
				g.players[index] = nil
				close(player.send)
			}
		case <-ticker.C:
			// check game over
			// send time update
			// decrement time
		case moveRequest := <-g.move:
			if g.state != playing {
				continue
			}

			player, index, ok := g.playerFromID(moveRequest.PlayerID)
			if !ok {
				continue
			}
			//check if move is valid
			// if valid move, check if game over
			// move, valid := g.Board.Move(moveRequest.Move)
			// send move to client (Diff if game over or not)
			// increment players time

			log.Println(player, index)
		}
	}
}
