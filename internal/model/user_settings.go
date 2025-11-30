package model

import (
	"time"

	"github.com/chistyakoviv/logbot/internal/i18n/language"
)

type UserSettings struct {
	UserId    int64
	Username  string
	Lang      int
	UpdatedAt time.Time
}

func (us *UserSettings) Language() string {
	return language.Language(us.Lang).String()
}

type UserSettingsInfo struct {
	Username string
	Lang     int
}
