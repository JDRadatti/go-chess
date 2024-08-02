package websocket

import (
    "time"
)
// MoveRequest is sent from client to request a move
type ConnectRequest struct {
	PlayerID string
	Time     time.Time
}

// MoveRequest is sent from client to request a move
type MoveRequest struct {
	PlayerID string
	GameID   string
	Move     string
}

// MoveAccept is sent from the server when a move is accepted
type MoveAccept struct {
	PlayerID string
	GameID   string
	Move     string
}

// MoveAccept is sent from the server when a move is denied
type MoveError struct {
	PlayerID string
	GameID   string
	Move     string
	Err      error
}


// GameError is sent from the client when joining a game fails
type GameError struct {
	PlayerID string
	GameID   string
	Err      string
}
