package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
)

type testCase struct {
	name      string
	playerIDs []string
	moves     []string
}

// Test two players joining
// Test four players joining
// Test one player joinging and timeout
// Test invalid player id joining
// Test valid player id joinikng with incorrect game id
// TestGames: input is a list of a bunch of games, weach witgh a playerIDs, gameID, and moves
// Test multiple games running at once
// TestLobby assumes the player already has a client ID recieved from

// Test invalid handshake (no join message)
// Test not sending the handshake (timeout?)
// How to know when both players are connected
// the /token endpoint (simulated with playerIDs)
func TestLobby(t *testing.T) {
	inputs := []testCase{
		{
			name:      "two player basic game with one move each",
			playerIDs: []string{"0", "1"},
			moves:     []string{"d4", "e4"},
		},
	}

	for _, tt := range inputs {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLobby()
			go l.Run()

			var wg sync.WaitGroup
			for _, playerID := range tt.playerIDs {
				// simulate play requests to lobby
				player := l.GetOrCreatePlayer(playerID)
				l.PlayerPool <- player
				<-player.InGame // wait for match making.

				p, ok := l.Players[playerID]
				assert.Equal(t, ok, true)
				assert.Equal(t, p.ID, playerID)

                // create connection
				wsHandler := &WSHandler{
					Lobby:  l,
					GameID: p.Game.ID,
				}
				s, conn := newWSServer(t, wsHandler)
				defer conn.Close()
				defer s.Close()

                // simulate game
				wg.Add(1)
				go func() {
					defer wg.Done()
					runGame(t, conn, player, tt)
				}()

			}
			wg.Wait()
		})
	}
}

func runGame(t *testing.T, conn *websocket.Conn, player *Player, tt testCase) {
	t.Helper()

	joinRequest := Inbound{
		Action:   JOIN,
		PlayerID: player.ID,
	}
	sendMessage(t, conn, joinRequest)
	joinSuccess := receiveWSMessage(t, conn)
	assert.Equal(t, string(joinSuccess.Action), JOIN_SUCCESS)

	color := player.Game.ColorFromPID(player.ID)

	for i := 0; i < len(tt.moves); i++ { // simulate game
		// Send message when it's player's move
		expectedPlayerID := i % 2
		if expectedPlayerID == color {
			moveRequest := Inbound{
				Action:   MOVE,
				Move:     tt.moves[i],
				PlayerID: player.ID,
				GameID:   player.Game.ID,
			}
			sendMessage(t, conn, moveRequest)
		}
		// receive
		// Should recieve a move confirmation if player sends a move
		// or if opponent sends a move
		expectedOut := Outbound{
			Action:   MOVE,
			Move:     tt.moves[i],
			GameID:   player.Game.ID,
			PlayerID: strconv.Itoa(expectedPlayerID),
		}
		actualOut := receiveWSMessage(t, conn)
		assert.Equal(t, expectedOut, actualOut)
	}

	assert.Equal(t, tt.moves, player.Game.AllMoves)
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

func sendMessage(t *testing.T, conn *websocket.Conn, msg Inbound) {
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
