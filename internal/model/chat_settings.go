package model

import "time"

type ChatSettings struct {
	ChatId         int64
	CollapsePeriod time.Duration
	MuteUntil      time.Time
	UpdatedAt      time.Time
}

type ChatSettingsInfo struct {
	CollapsePeriod time.Duration
	MuteUntil      time.Time
}
