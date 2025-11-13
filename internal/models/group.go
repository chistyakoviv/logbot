package models

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID        uuid.UUID
	ChatID    int64
	Token     string
	CreatedAt time.Time
}
