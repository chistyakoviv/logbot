package chat_settings

import (
	"time"

	"github.com/chistyakoviv/logbot/internal/model"
)

type ChatSettingsRow struct {
	ChatId         int64         `db:"chat_id"`
	CollapsePeriod time.Duration `db:"collapse_period"`
	MuteUntil      time.Time     `db:"mute_until"`
	SilenceUntil   time.Time     `db:"silence_until"`
	UpdatedAt      time.Time     `db:"updated_at"`
}

func (r *ChatSettingsRow) Values() []any {
	return []any{
		r.ChatId,
		r.CollapsePeriod,
		r.MuteUntil,
		r.SilenceUntil,
		r.UpdatedAt,
	}
}

func ToModel(r *ChatSettingsRow) *model.ChatSettings {
	if r == nil {
		return nil
	}
	return &model.ChatSettings{
		ChatId:         r.ChatId,
		CollapsePeriod: r.CollapsePeriod,
		MuteUntil:      r.MuteUntil,
		SilenceUntil:   r.SilenceUntil,
		UpdatedAt:      r.UpdatedAt,
	}
}

func FromModel(m *model.ChatSettings) ChatSettingsRow {
	if m == nil {
		return ChatSettingsRow{}
	}
	return ChatSettingsRow{
		ChatId:         m.ChatId,
		CollapsePeriod: m.CollapsePeriod,
		MuteUntil:      m.MuteUntil,
		SilenceUntil:   m.SilenceUntil,
		UpdatedAt:      m.UpdatedAt,
	}
}
