package model

import (
	"github.com/google/uuid"
	"time"
)

type Lobby struct {
	Games []Game
}

type Game struct {
	ID          uuid.UUID // Unique id for this game
	InitialTime uint16    // Time control in seconds. Max is 65,535 seconds
	StartTime   time.Time // Time this game started
	White       Player // Player playing with the white pieces
	Black       Player // Player playing with the black pieces
}

type Player struct {
	name string
}
