package api

import (
	"encoding/json"
	"github.com/JDRadatti/reptile/internal/websocket"
	"github.com/google/uuid"
	"log"
	"net/http"
)

// HandlePlay handles online game requests
func HandlePlay(w http.ResponseWriter, r *http.Request, lobby *websocket.Lobby) {

	gameRequest := &websocket.GameRequest{}
	err := json.NewDecoder(r.Body).Decode(gameRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	playerID := gameRequest.PlayerID
	if err := uuid.Validate(string(gameRequest.PlayerID)); err != nil {
		playerID = websocket.GeneratePlayerID()
		gameRequest.PlayerID = playerID
	}

	payload, err := json.Marshal(lobby.Match(gameRequest))
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if _, err := w.Write(payload); err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
