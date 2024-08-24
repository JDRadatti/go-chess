package api

import (
	"github.com/JDRadatti/reptile/internal/websocket"
	"log"
	"net/http"
)

func HandleLobby(w http.ResponseWriter, l *websocket.Lobby) {
	payload := []byte(l.String())
	if _, err := w.Write(payload); err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
