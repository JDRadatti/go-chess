package api

import (
	"encoding/json"
	"github.com/JDRadatti/reptile/internal/websocket"
	"net/http"
)

type GameOptions struct {
	Time      int
	Increment int
}

// HandlePlay handles online game requests
func HandlePlay(w http.ResponseWriter, r *http.Request, lobby *websocket.Lobby) {

	var options GameOptions
	err := json.NewDecoder(r.Body).Decode(&options)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//player := websocket.NewPlayer(lobby)
	//lobby.PlayerPool <- player
    //go player.WaitForGame(w)
}
