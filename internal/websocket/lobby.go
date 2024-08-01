package websocket

import (
	"log"
)

type Lobby struct {
	Games      map[string]*Game   // Current running games (has both players)
	Players    map[string]*Player // Current Players in a game
	GamePool   chan *Game         // Current waiting games (only one player)
	PlayerPool chan *Player       // Players who are waiting for a game
}

func NewLobby() *Lobby {
	return &Lobby{
		Games:      make(map[string]*Game),
		Players:    make(map[string]*Player),
		GamePool:   make(chan *Game),
		PlayerPool: make(chan *Player),
	}
}

func (l *Lobby) GetPlayer(playerID string) *Player {
	return l.Players[playerID]
}

func (l *Lobby) GetGame(id string) *Game {
	return l.Games[id]
}

func (l *Lobby) addGame(game *Game) {
	if _, ok := l.Games[game.ID]; ok {
		panic("cannot add game that aready exists")
	}
	l.Games[game.ID] = game
}

// findGame handles match making.
// eventually want better match making but for now
// it simply matches people based on open games
// TODO add elo, color, and game mode match making
func (l *Lobby) findGame(p *Player) {
	if p.Game != nil {
		log.Printf("already in game %s", p.Game.ID)
	}

	var game *Game
	if len(l.GamePool) >= 1 {
		game = <-l.GamePool
		game.addPlayer(p)
	} else {
		game = newGame()
		game.addPlayer(p)
		l.addGame(game)
	}
	p.Game = game
	close(p.InGame)
}

func (l *Lobby) Run() {
	for {
		select {
		case player := <-l.PlayerPool:
			l.findGame(player)
		}
	}
}
