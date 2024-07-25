package websocket

import (
	"log"
	"net/http"
)

type Lobby struct {
	Games      map[string]*Game // Current running games (has both players)
	GamePool   chan *Game       // Current waiting games (only one player)
	PlayerPool chan *Player     // Players who are waiting for a game
}

func NewLobby() *Lobby {
	return &Lobby{
		Games:      make(map[string]*Game),
		GamePool:   make(chan *Game),
		PlayerPool: make(chan *Player),
	}
}

// JoinGame handles the match making for online games
// JoinGame will create a game if none are available
func (l *Lobby) JoinGame(r *http.Request) {

}

func (l *Lobby) getGame(id string) (*Game, bool) {
	game, ok := l.Games[id]
	return game, ok
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
}

func (l *Lobby) Run() {
	for {
        log.Println("lobby")
		select {
		case player := <-l.PlayerPool:
			// Player is waiting
			l.findGame(player)
		}
	}
}
