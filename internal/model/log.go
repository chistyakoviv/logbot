package model

import (
	"time"

	"github.com/google/uuid"
)

type Log struct {
	Id        int
	Token     uuid.UUID
	Data      string
	Label     string
	Hash      string
	CreatedAt time.Time
}

type LogInfo struct {
	Token uuid.UUID
	Data  string
	Label string
}
