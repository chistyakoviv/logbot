package model

import "time"

type Label struct {
	ChatId    int64
	Username  string
	Labels    []string
	UpdatedAt time.Time
}

type LabelKey struct {
	ChatId   int64
	Username string
}
