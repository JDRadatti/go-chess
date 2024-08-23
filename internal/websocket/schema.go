package websocket

import (
	"encoding/json"
	"github.com/JDRadatti/reptile/internal/chess"
	"log"
)

const ( // incoming action
	JOIN         = "join"
	MOVE         = "move"
	RESIGN       = "resign"
	DRAW         = "draw" // both players must send to accept draw
	DRAW_REQUEST = "draw_request"
	DRAW_ACCEPT  = "draw_accept"
	DRAW_DENY    = "draw_deny"
	ABORT        = "abort"
)

const ( // outgoing status
	JOIN_SUCCESS   = "join_success"
	JOIN_FAIL      = "join_fail"
	GAME_START     = "game_start"
	GAME_END       = "game_end"
	MOVE_SUCCESS   = "move_success"
	MOVE_FAIL      = "move_fail"
	RESIGN_SUCCESS = "resign_success"
	DRAW_SUCCESS   = "draw_success"
	TIME_UPDATE    = "time_update"
)

type Inbound struct {
	Action   string
	Move     string
	PlayerID PlayerID
	GameID   GameID
}

type Outbound struct {
	Action    string
	Move      string
	FEN       string
	PlayerID  PlayerID
	GameID    GameID
	WhiteTime int
	BlackTime int
	Increment int
	Player    chess.Player
	Turn      chess.Player
}

// GameRequest is sent from the client when wanting to join a game
type GameRequest struct {
	PlayerID  PlayerID
	Time      int
	Increment int
}

// GameResponse is sent from the client after joining a game
type GameResponse struct {
	PlayerID PlayerID
	GameID   GameID
	Player   chess.Player
}

func unmarshal(message []byte) (*Inbound, bool) {
	in := &Inbound{}
	err := json.Unmarshal(message, in)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, false
	}
	return in, true
}

func marshal(o *Outbound) ([]byte, bool) {

	message, err := json.Marshal(o)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, false
	}
	return message, true
}
