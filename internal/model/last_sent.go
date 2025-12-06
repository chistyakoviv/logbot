package model

import "time"

type LastSent struct {
	ChatId    int64
	Token     string
	Hash      string
	UpdatedAt time.Time
}

type LastSentKey struct {
	ChatId int64
	Token  string
	Hash   string
}
