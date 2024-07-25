package api

import (
	"github.com/JDRadatti/reptile/internal/websocket"
	"log"
	"net/http"
    "encoding/json"
)

type GameOptions struct {
    Time int
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
	log.Println(options)
    gameID := "123"
    w.Write([]byte(gameID))
}
