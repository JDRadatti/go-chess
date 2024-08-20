package api

import (
	"encoding/json"
	"github.com/JDRadatti/reptile/internal/chess"
	"github.com/JDRadatti/reptile/internal/websocket"
	"github.com/google/uuid"
	"log"
	"net/http"
)

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
	Color    chess.Player
}

// HandlePlay handles online game requests
func HandlePlay(w http.ResponseWriter, r *http.Request, lobby *websocket.Lobby) {
	var gameRequest *GameRequest = &GameRequest{}
	err := json.NewDecoder(r.Body).Decode(gameRequest)
	if err != nil {
		log.Println("err", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	playerID := gameRequest.PlayerID
	if err := uuid.Validate(playerID); err != nil {
		playerID = websocket.GenerateID()
	}

	player := lobby.GetOrCreatePlayer(playerID, gameRequest.Time, gameRequest.Increment)
	if player.Game == nil {
		lobby.PlayerPool <- player

		<-player.InGame // wait for match making.
		if player.Game == nil {
			http.Error(w, websocket.MATCHMAKING_ERROR, http.StatusBadRequest)
		}
	} 

	payload := GameAccepted{
		GameID:   player.Game.ID,
		PlayerID: player.ID,
		Color:    player.Game.ColorFromPID(player.ID),
	}

	marshled, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error: %v", err)
	}

	if _, err := w.Write(marshled); err != nil {
		log.Printf("error: %v", err)
	}
}
