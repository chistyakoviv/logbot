package model

import "time"

type ChatSettings struct {
	ChatId         int64
	CollapsePeriod time.Duration
	MuteUntil      time.Time
	SilenceUntil   time.Time
	UpdatedAt      time.Time
}

func (cs *ChatSettings) IsMuted(now time.Time) (bool, time.Duration) {
	timeRemaining := cs.MuteUntil.Sub(now)
	return !cs.MuteUntil.IsZero() && now.Before(cs.MuteUntil), timeRemaining
}

func (cs *ChatSettings) IsSilenced(now time.Time) (bool, time.Duration) {
	timeRemaining := cs.SilenceUntil.Sub(now)
	return !cs.SilenceUntil.IsZero() && now.Before(cs.SilenceUntil), timeRemaining
}

func (cs *ChatSettings) IsCollapsed(now time.Time, lastSentTimestamp time.Time) (bool, time.Duration) {
	timeSinceLastSent := now.Sub(lastSentTimestamp)
	return cs.CollapsePeriod > 0 && timeSinceLastSent < cs.CollapsePeriod, timeSinceLastSent
}
