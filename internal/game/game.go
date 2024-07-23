package game

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func Play(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
    log.Println("REQUEST: ", r.Header)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade Error:", err)
		return
	}

	conn.WriteMessage(websocket.TextMessage, []byte("Hello from server"))
	var time int // total game time in seconds
	for time = 60; time > 0; time-- {
        conn.WriteMessage(websocket.TextMessage, []byte("time remaining(s): " + string(time)))
	}

}
