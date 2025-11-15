package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	Id        uuid.UUID
	ChatId    int64
	Token     string
	CreatedAt time.Time
}

type SubscriptionInfo struct {
	ChatId int64
	Token  string
}
