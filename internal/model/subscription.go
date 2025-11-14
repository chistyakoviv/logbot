package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID        uuid.UUID
	ChatID    int64
	Token     string
	CreatedAt time.Time
}

type SubscriptionInfo struct {
	ChatID int64
	Token  string
}
