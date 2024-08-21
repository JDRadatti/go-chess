package websocket

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JDRadatti/reptile/internal/chess"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name      string
	gameID    string
	playerID  string
	time      int
	inc       int
	inbounds  []*Inbound
	outbounds []Outbound
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
			name:     "VALID: gameID in lobby, playerID matches",
			gameID:   "0",
			playerID: "1",
			time:     time,
			inc:      increment,
			inbounds: []*Inbound{
				{
					Action:   JOIN,
					PlayerID: "1",
				},
			},
			outbounds: []Outbound{
				{
					Action:   JOIN_SUCCESS,
					FEN:      string(startingFEN),
                    WhiteTime: time,
                    BlackTime: time,
                    Increment: increment,
					PlayerID: "1",
					GameID:   "0",
				},
			},
		},
	}

	for _, tt := range inputs {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLobby()
			game := NewGame(l, tt.time, tt.inc)
			game.id = GameID(tt.gameID)
			l.Join(PlayerID(tt.playerID), game)

			// create connection
			wsHandler := &WSHandler{
				Lobby:  l,
				GameID: GameID(tt.gameID),
			}
			s, conn := newWSServer(t, wsHandler)
			defer conn.Close()
			defer s.Close()

			for i, inbound := range tt.inbounds {
				sendMessage(t, conn, inbound)
				joinSuccess := receiveWSMessage(t, conn)
				assert.Equal(t, tt.outbounds[i], joinSuccess)
			}

			// Clean and test Clean worked
			l.Clean(game.id, PlayerID(tt.playerID), "")
			_, ok := l.Players[PlayerID(tt.playerID)]
			assert.Equal(t, false, ok)
			_, ok = l.Games[GameID(tt.gameID)]
			assert.Equal(t, false, ok)
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
