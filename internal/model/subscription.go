package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const MaxProjectNameLength = 100

var ErrTokenAlreadyExists = errors.New("token already exists")

type Subscription struct {
	Id          uuid.UUID
	ChatId      int64
	Token       string
	ProjectName string
	CreatedAt   time.Time
}

type SubscriptionInfo struct {
	ChatId      int64
	Token       string
	ProjectName string
}
