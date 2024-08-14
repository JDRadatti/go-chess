package websocket

import (
	"log"
)

var (
	maxGames = 10
)

type Lobby struct {
	Games      map[string]*Game   // Current running games (has both players)
	Players    map[string]*Player // Current Players in a game.
	GamePool   chan *Game         // Current waiting games (only one player)
	PlayerPool chan *Player       // Players who are waiting for a game [DEP]
}

func NewLobby() *Lobby {
	return &Lobby{
		Games:      make(map[string]*Game),
		Players:    make(map[string]*Player),
		GamePool:   make(chan *Game, maxGames),
		PlayerPool: make(chan *Player),
	}
}

func (l *Lobby) GetGame(id string) (*Game, bool) {
	game, ok := l.Games[id]
	return game, ok
}

func (l *Lobby) GetPlayer(id string) (*Player, bool) {
	player, ok := l.Players[id]
	return player, ok
}

func (l *Lobby) AddPlayer(player *Player) bool {
	log.Println("ADD PLAEYR", player)
	if _, ok := l.Players[player.ID]; ok {
		log.Println("HERE already playing in diff game")
		return false
	}
	l.Players[player.ID] = player
	return true
}

func (l *Lobby) GetOrCreatePlayer(playerID string) *Player {
	if player, ok := l.GetPlayer(playerID); ok {
		log.Printf("player already in game %s", playerID)
		return player
	}
	player := NewPlayer(l, nil)
	player.ID = playerID // must be validated before this function
	l.Players[playerID] = player
	return player
}

func (l *Lobby) Run() {
	for {
		select {
		case player := <-l.PlayerPool:
			if player.Game != nil {
				log.Printf("player %s already in game %s", player.ID, player.Game.ID)
				continue
			}
			// TODO: handle different game options
			var game *Game
			select {
			case g, ok := <-l.GamePool:
				if ok {
					game = g
				} else {
					panic("GamePool channel closed")
				}
			default:
				game = newGame()
				l.GamePool <- game
			}

			if err := game.addPlayer(player); err != nil {
				log.Printf("error %v", err)
				close(player.InGame)
				continue
			}

			player.Game = game
			l.Games[game.ID] = game
			close(player.InGame)
			// TODO: make sure to remove gameID and playerID when game ends
		}
	}
}
