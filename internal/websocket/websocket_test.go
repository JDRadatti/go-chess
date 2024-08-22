package websocket

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JDRadatti/reptile/internal/chess"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// there can be x number of players. x should be the length of all
// slices
type testCase struct {
	name      string
	time      int
	inc       int
	gameID    string
	playerID  []PlayerID
	join      []bool // join[i] true if player i should join before handshake
	inbounds  [][]*Inbound
	outbounds [][]Outbound // "*" as outbound.playerID means any valid uuid
}

// valid handshakes:
// gameID in lobby, no playerID but game not full
// gameID in lobby, invalid playerID but game not full
//
// invalid handshakes:
// no handshake message
// invalid handshake message
// gameID not in lobby/invalid
// gameID in lobby but game is over
// gameID in lobby but invalid/not found playerID
func TestHandshake(t *testing.T) {
	board := chess.NewBoardClassic()
	startingFEN := board.FEN()
	time := 180
	increment := 0
	players := 3
	playerIDs := make([]PlayerID, players)
	joins := make([]*Inbound, players)
	success := make([]Outbound, players)
	fail := make([]Outbound, players)
	player := []chess.Player{chess.WHITE, chess.BLACK, chess.INVALID_PLAYER}
	for i := range players {
		playerID := GeneratePlayerID()
		j := &Inbound{
			Action:   JOIN,
			PlayerID: playerID,
		}
		s := Outbound{
			Action:    JOIN_SUCCESS,
			FEN:       string(startingFEN),
			WhiteTime: time,
			BlackTime: time,
			Increment: increment,
			PlayerID:  playerID,
			GameID:    "0",
			Player:    player[i],
		}
		f := Outbound{
			Action: JOIN_FAIL,
		}
		playerIDs[i] = playerID
		joins[i] = j
		success[i] = s
		fail[i] = f
	}

	inputs := []testCase{
		{
			name:     "VALID: both players already in game. valid join request.",
			gameID:   "0",
			playerID: []PlayerID{playerIDs[0], playerIDs[1]},
			join:     []bool{true, true},
			time:     time,
			inc:      increment,
			inbounds: [][]*Inbound{
				{joins[0]},
				{joins[1]},
			},
			outbounds: [][]Outbound{
				{success[0]},
				{success[1]},
			},
		},
		{
			name:     "VALID: one player not already in game. valid join request. random id.",
			gameID:   "0",
			playerID: []PlayerID{playerIDs[0], playerIDs[1]},
			join:     []bool{true, false},
			time:     time,
			inc:      increment,
			inbounds: [][]*Inbound{
				{joins[0]},
				{joins[1]},
			},
			outbounds: [][]Outbound{
				{success[0]},
				{success[1]},
			},
		},
		{
			name:     "VALID: two valid join requests and a third fail when joining same gameID.",
			gameID:   "0",
			playerID: []PlayerID{playerIDs[0], playerIDs[1], playerIDs[2]},
			join:     []bool{true, false, false},
			time:     time,
			inc:      increment,
			inbounds: [][]*Inbound{
				{joins[0]},
				{joins[1]},
				{joins[2]},
			},
			outbounds: [][]Outbound{
				{success[0]},
				{success[1]},
				{fail[2]},
			},
		},
	}

	for _, tt := range inputs {
		l := NewLobby()
		game := NewGame(l, tt.time, tt.inc)
		game.id = GameID(tt.gameID)

		t.Run(tt.name, func(t *testing.T) {
			for i := range tt.inbounds {
				// create connection
				wsHandler := &WSHandler{
					Lobby:  l,
					GameID: GameID(tt.gameID),
				}
				s, conn := newWSServer(t, wsHandler)

				if tt.join[i] {
					l.Join(tt.playerID[i], game)
					game.addPlayerID(tt.playerID[i])
				}
				for j, inbound := range tt.inbounds[i] {
					sendMessage(t, conn, inbound)
					joinSuccess := receiveWSMessage(t, conn)

					// if playerID = *, set PlayerID to the one returned by
					// server. the server will send the client a new uuid if
					// 1. the given uuid is not in a different game
					// 2. the given game is not full
					if tt.outbounds[i][j].PlayerID == "*" {
						tt.outbounds[i][j].PlayerID = joinSuccess.PlayerID
						if err := uuid.Validate(string(tt.outbounds[i][j].PlayerID)); err != nil {
							assert.Fail(t, "server sent invalid playerID")
						}
					}
					assert.Equal(t, tt.outbounds[i][j], joinSuccess)
					conn.Close()
					s.Close()
				}

			}

			// Clean and test Clean worked
			for i := range tt.inbounds {
				l.Clean(game.id, PlayerID(tt.playerID[i]), "")
				_, ok := l.Players[PlayerID(tt.playerID[i])]
				assert.Equal(t, false, ok)
				_, ok = l.Games[GameID(tt.gameID)]
				assert.Equal(t, false, ok)
			}
		})
	}
}

func newWSServer(t *testing.T, h http.Handler) (*httptest.Server, *websocket.Conn) {
	t.Helper()

	s := httptest.NewServer(h)
	wsURL := httpToWs(t, s.URL)

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	return s, ws
}

func httpToWs(t *testing.T, url string) string {
	t.Helper()

	if s, found := strings.CutPrefix(url, "https"); found {
		return "wss" + s
	} else if s, found := strings.CutPrefix(url, "http"); found {
		return "ws" + s
	}
	return url
}

func sendMessage(t *testing.T, conn *websocket.Conn, msg *Inbound) {
	t.Helper()

	m, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, m); err != nil {
		t.Fatalf("%v", err)
	}
}

func receiveWSMessage(t *testing.T, conn *websocket.Conn) Outbound {
	t.Helper()

	_, m, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	var response Outbound
	err = json.Unmarshal(m, &response)
	if err != nil {
		t.Fatal(err)
	}

	return response
}
