package websocket

import (
)

type Lobby struct {
    Games map[string]*Game
}

func NewLobby() *Lobby {
    return &Lobby{
        Games: make(map[string]*Game),
    }
}

func (l *Lobby) getGame(id string) (*Game, bool) {
    game, ok := l.Games[id]
    return game, ok
}
