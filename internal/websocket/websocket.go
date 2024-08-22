package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WSHandler struct {
	Lobby  *Lobby
	GameID GameID
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (ws *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	player, response, ok := ws.handshake(conn)
	if message, success := marshal(response); success {
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Printf("error: %v", err)
		}
	} else {
		log.Println(err)
		return
	}

	if ok {
		go player.write()
		go player.read()
	}

}

func (ws *WSHandler) handshake(conn *websocket.Conn) (*Player, *Outbound, bool) {

	_, message, err := conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("error: %v", err)
		}
		return nil, handshakeFail(), false
	}

	in, ok := unmarshal(message)
	if !ok || in.Action != JOIN {
		log.Printf("error: %v", err)
		return nil, handshakeFail(), false
	}

	if game, ok := ws.Lobby.GetGameFromPlayerID(in.PlayerID); ok {
		if game.id == ws.GameID {
			return NewPlayer(ws.Lobby, conn, game),
				handshakeSuccess(in.PlayerID, game),
				true
		} else {
			return nil, handshakeFail(), false
		}
	}

	// Player not already in game (opened game link)
	if game, ok := ws.Lobby.GetGameFromGameID(ws.GameID); ok {

		if ok := in.PlayerID.validate(); ok {
			in.PlayerID = GeneratePlayerID()
		}

		if _, ok := game.addPlayerID(in.PlayerID); !ok {
			return nil, handshakeFail(), false
		}

		ws.Lobby.Join(in.PlayerID, game)

		return NewPlayer(ws.Lobby, conn, game),
			handshakeSuccess(in.PlayerID, game),
			true
	}

	return nil, handshakeFail(), false
}

func handshakeFail() *Outbound {
	return &Outbound{
		Action: JOIN_FAIL,
	}
}

func handshakeSuccess(pid PlayerID, g *Game) *Outbound {
	return &Outbound{
		Action:    JOIN_SUCCESS,
		PlayerID:  pid,
		GameID:    g.id,
		FEN:       string(g.board.FEN()),
		Turn:      g.board.Turn(),
		WhiteTime: g.timeRemaining[whiteIndex],
		BlackTime: g.timeRemaining[blackIndex],
		Player:    g.playerTypeFromID(pid),
	}
}
