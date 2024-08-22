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

type testCase struct {
	name      string
	time      int
	inc       int
	gameID    string
	playerID  [2]string
	join      [2]bool // join[i] true if player i should join before handshake
	inbounds  [2][]*Inbound
	outbounds [2][]Outbound // "*" as outbound.playerID means any valid uuid
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
	inputs := []testCase{
		{
			name:     "VALID: both players already in game. valid join request.",
			gameID:   "0",
			playerID: [2]string{"1", "2"},
			join:     [2]bool{true, true},
			time:     time,
			inc:      increment,
			inbounds: [2][]*Inbound{
				{
					{
						Action:   JOIN,
						PlayerID: "1",
					},
				},
				{
					{
						Action:   JOIN,
						PlayerID: "2",
					},
				},
			},
			outbounds: [2][]Outbound{
				{
					{
						Action:    JOIN_SUCCESS,
						FEN:       string(startingFEN),
						WhiteTime: time,
						BlackTime: time,
						Increment: increment,
						PlayerID:  "1",
						GameID:    "0",
					},
				},
				{
					{
						Action:    JOIN_SUCCESS,
						FEN:       string(startingFEN),
						WhiteTime: time,
						BlackTime: time,
						Increment: increment,
						PlayerID:  "2",
						GameID:    "0",
					},
				},
			},
		},
		{
			name:     "VALID: one player not already in game. valid join request. random id.",
			gameID:   "0",
			playerID: [2]string{"1", "2"},
			join:     [2]bool{true, false},
			time:     time,
			inc:      increment,
			inbounds: [2][]*Inbound{
				{
					{
						Action:   JOIN,
						PlayerID: "1",
					},
				},
				{
					{
						Action:   JOIN,
						PlayerID: "2",
					},
				},
			},
			outbounds: [2][]Outbound{
				{
					{
						Action:    JOIN_SUCCESS,
						FEN:       string(startingFEN),
						WhiteTime: time,
						BlackTime: time,
						Increment: increment,
						PlayerID:  "1",
						GameID:    "0",
					},
				},
				{
					{
						Action:    JOIN_SUCCESS,
						FEN:       string(startingFEN),
						WhiteTime: time,
						BlackTime: time,
						Increment: increment,
						PlayerID:  "*",
						GameID:    "0",
					},
				},
			},
		},
	}

	for _, tt := range inputs {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLobby()
			game := NewGame(l, tt.time, tt.inc)
			game.id = GameID(tt.gameID)

			for i := range tt.inbounds {
				// create connection
				wsHandler := &WSHandler{
					Lobby:  l,
					GameID: GameID(tt.gameID),
				}
				s, conn := newWSServer(t, wsHandler)
				defer conn.Close()
				defer s.Close()

				if tt.join[i] {
					l.Join(PlayerID(tt.playerID[i]), game)
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
