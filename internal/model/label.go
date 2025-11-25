package model

import "time"

type Label struct {
	ChatId    int64
	UserId    int64
	Labels    []string
	UpdatedAt time.Time
}

type LabelKey struct {
	ChatId int64
	UserId int64
}
