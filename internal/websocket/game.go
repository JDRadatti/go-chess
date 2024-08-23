package websocket

import (
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
	resign        chan *Inbound
	abort         chan *Inbound
	draw          chan *Inbound
	pendingDraw   int
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
		resign:        make(chan *Inbound),
		draw:          make(chan *Inbound),
		abort:         make(chan *Inbound),
		join:          make(chan *Player),
		leave:         make(chan *Player),
		board:         &board,
		timeRemaining: [2]int{time, time},
		players:       [2]*Player{},
		playerIDs:     [2]PlayerID{},
		pendingDraw:   -1,
		increment:     increment,
		lobby:         l,
		state:         waiting,
	}
	go newGame.play()
	return newGame
}

func generateGameID() GameID {
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Printf("error %s", err)
	}
	return GameID(uuid.String()[:8])
}

func (g *Game) clean() {
	g.lobby.Clean(g.id, g.playerIDs[whiteIndex], g.playerIDs[blackIndex])
}

func (g *Game) playerFromID(playerID PlayerID) (*Player, int, bool) {
	if g.playerIDs[whiteIndex] == playerID {
		return g.players[whiteIndex], whiteIndex, true
	} else if g.playerIDs[blackIndex] == playerID {
		return g.players[blackIndex], blackIndex, true
	} else {
		return nil, -1, false
	}
}

func (g *Game) playerType(playerID PlayerID) chess.Player {
	index, _ := g.playerIndex(playerID)
	return chess.Player(index)
}

func (g *Game) playerIndex(playerID PlayerID) (int, bool) {
	if g.playerIDs[whiteIndex] == playerID {
		return whiteIndex, true
	} else if g.playerIDs[blackIndex] == playerID {
		return blackIndex, true
	} else {
		return -1, false
	}
}

func (g *Game) addPlayerID(playerID PlayerID) (int, bool) {
	if g.playerIDs[whiteIndex] == "" {
		g.playerIDs[whiteIndex] = playerID
		return whiteIndex, true
	} else if g.playerIDs[blackIndex] == "" {
		g.playerIDs[blackIndex] = playerID
		return blackIndex, true
	} else {
		return -1, false
	}
}

func (g *Game) currentPlayerIndex() int {
	return int(g.board.Turn())
}

func (g *Game) bothPlayersConnected() bool {
	return g.players[whiteIndex] != nil && g.players[blackIndex] != nil
}

func (g *Game) out(action string, pid PlayerID) *Outbound {
	return &Outbound{
		Action:    action,
		Move:      g.board.LastMove(),
		PlayerID:  pid, // Player who made the move
		GameID:    g.id,
		FEN:       string(g.board.FEN()),
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
			if index, ok := g.playerIndex(player.id); ok {
				g.players[index] = player
			}
			if g.bothPlayersConnected() {
				startOut := g.out(GAME_START, "")
				g.players[whiteIndex].send <- startOut
				g.players[blackIndex].send <- startOut
				g.state = playing
			}
		case player := <-g.leave:
			if index, ok := g.playerIndex(player.id); ok {
				g.players[index] = nil
				close(player.send)
			}
		case <-ticker.C:
			if g.state != playing {
				continue
			}
			currentI := g.currentPlayerIndex()
			if g.timeRemaining[currentI] < 0 {
				out := g.out(GAME_END, "")
				g.players[whiteIndex].send <- out
				g.players[blackIndex].send <- out
				return
			}
			out := g.out(TIME_UPDATE, "")
			g.players[whiteIndex].send <- out
			g.players[blackIndex].send <- out
			g.timeRemaining[currentI]--
		case moveRequest := <-g.move:
			if g.state != playing {
				continue
			}

			player, index, ok := g.playerFromID(moveRequest.PlayerID)
			if !ok || index != g.currentPlayerIndex() {
				continue
			}

			move, valid := g.board.Move(moveRequest.Move)
			if valid {
				out := g.out(MOVE_SUCCESS, player.id)
				out.Move = move
				g.players[whiteIndex].send <- out
				g.players[blackIndex].send <- out
				g.pendingDraw = -1
			}

			if _, over := g.board.GameOver(); over {
				out := g.out(GAME_END, player.id)
				g.players[whiteIndex].send <- out
				g.players[blackIndex].send <- out
				return
			}

			g.timeRemaining[index] += g.increment
		case resignRequest := <-g.resign:
			if index, ok := g.playerIndex(resignRequest.PlayerID); ok {
				out := g.out(RESIGN, g.playerIDs[index])
				g.players[whiteIndex].send <- out
				g.players[blackIndex].send <- out
				return
			}
		case abortRequest := <-g.abort:
			if index, ok := g.playerIndex(abortRequest.PlayerID); ok {
				if !g.board.CanAbort() {
					continue
				}
				out := g.out(ABORT, g.playerIDs[index])
				g.players[whiteIndex].send <- out
				g.players[blackIndex].send <- out
				return
			}
		case drawRequest := <-g.draw:
			if index, ok := g.playerIndex(drawRequest.PlayerID); ok {
				if g.pendingDraw == -1 && drawRequest.Action == DRAW_REQUEST {
					out := g.out(DRAW_REQUEST, g.playerIDs[index])
					g.players[(index+1)%2].send <- out // send to other index
					g.pendingDraw = index
				} else if g.pendingDraw == (index+1)%2 && drawRequest.Action == DRAW_ACCEPT {
					out := g.out(DRAW, g.playerIDs[index])
					g.players[whiteIndex].send <- out
					g.players[blackIndex].send <- out
					return
				} else if drawRequest.Action == DRAW_DENY {
					out := g.out(DRAW_DENY, g.playerIDs[index])
					out.Player = g.playerType(out.PlayerID)
					g.players[whiteIndex].send <- out
					g.players[blackIndex].send <- out
					g.pendingDraw = -1
				}
			}

		}
	}
}
