package websocket

// MoveRequest is sent from client to request a move
type MoveRequest struct {
	PlayerID string
	Move     string
}

// MoveAccept is sent from the server when a move is accepted
type MoveAccept struct {
	PlayerID string
	Move     string
}

// MoveAccept is sent from the server when a move is denied
type MoveError struct {
	PlayerID string
	Move     string
	Err      error
}

// GameRequest is sent from the client when wanting to join a game
type GameRequest struct {
	PlayerID  string
	Time      int
	Increment int
}

// GameAccepted is sent from the client when wanting to join a game
type GameAccepted struct {
	PlayerID string
	GameID   string
}

// GameError is sent from the client when joining a game fails
type GameError struct {
	PlayerID string
	Err      string
}
