package api

import (
	"github.com/google/uuid"
	"log"
	"net/http"
)

func HandleToken(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Printf("error %s", err)
	}
	w.Write([]byte(uuid.String()))
}
