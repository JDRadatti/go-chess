package websocket

import (
	"fmt"
	"github.com/JDRadatti/reptile/internal/chess"
	"log"
	"strings"
)

var (
	gameLimit = 10
)

type Lobby struct {
	Games    map[GameID]*Game   // Current running games (has both players)
	Players  map[PlayerID]*Game // Current Players in a game.
	GamePool chan *Game         // Current waiting games (only one player)
}

func NewLobby() *Lobby {
	return &Lobby{
		Games:    make(map[GameID]*Game),
		Players:  make(map[PlayerID]*Game),
		GamePool: make(chan *Game, gameLimit),
	}
}

func (l *Lobby) Clean(gid GameID, pid1 PlayerID, pid2 PlayerID) {
	if _, ok := l.Games[gid]; ok {
		delete(l.Games, gid)
	}
	if _, ok := l.Players[pid1]; ok {
		delete(l.Players, pid1)
	}
	if _, ok := l.Players[pid2]; ok {
		delete(l.Players, pid2)
	}
}

func (l *Lobby) GetGameFromGameID(id GameID) (*Game, bool) {
	game, ok := l.Games[id]
	return game, ok
}

func (l *Lobby) GetGameFromPlayerID(id PlayerID) (*Game, bool) {
	player, ok := l.Players[id]
	return player, ok
}

func (l *Lobby) Join(playerID PlayerID, game *Game) bool {
	if game == nil {
		return false
	}
	if _, ok := l.Players[playerID]; ok {
		return false
	}
	l.Players[playerID] = game

	if _, ok := l.Games[game.id]; !ok {
		l.Games[game.id] = game
	}
	return true
}

func (l *Lobby) Success(pid PlayerID, gid GameID, i int) *GameResponse {
	return &GameResponse{
		PlayerID: pid,
		GameID:   gid,
		Player:   chess.Player(i),
	}
}

func (l *Lobby) Fail() *GameResponse {
	return &GameResponse{
		PlayerID: "",
		GameID:   "",
		Player:   -1,
	}
}

func (l *Lobby) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("#games %d, #players %d, #gamepool %d\n", len(l.Games), len(l.Players), len(l.GamePool)))
	for _, game := range l.Games {
		builder.WriteString(fmt.Sprintf("-- %s --\n", string(game.id)))
		for i, pid := range game.playerIDs {
			if pid != "" {
				builder.WriteString(fmt.Sprintf("%d %s\n", i, string(pid)))
			}
		}
	}

	return builder.String()
}

func (l *Lobby) Match(request *GameRequest) *GameResponse {
	if game, ok := l.GetGameFromPlayerID(request.PlayerID); ok {
		log.Printf("player %s already in game %s", request.PlayerID, game.id)
		if index, ok := game.playerIndex(request.PlayerID); ok {
			return l.Success(request.PlayerID, game.id, index)
		} else {
			return l.Fail()
		}
	}
	// TODO: handle different game options
	var game *Game
	select {
	case g := <-l.GamePool:
		if g.state == over {
			game = NewGame(l, request.Time, request.Increment)
			l.GamePool <- game
		} else {
			game = g
		}
	default:
		game = NewGame(l, request.Time, request.Increment)
		l.GamePool <- game
	}

	if index, ok := game.addPlayerID(request.PlayerID); ok {
		l.Join(request.PlayerID, game)
		return l.Success(request.PlayerID, game.id, index)
	} else {
		return l.Fail()
	}
}
