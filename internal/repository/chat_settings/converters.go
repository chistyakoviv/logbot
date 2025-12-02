package chat_settings

import (
	"time"

	"github.com/chistyakoviv/logbot/internal/model"
)

type ChatSettingsRow struct {
	ChatId         int64         `db:"chat_id"`
	CollapsePeriod time.Duration `db:"collapse_period"`
	SilenceUntil   time.Time     `db:"silence_until"`
	UpdatedAt      time.Time     `db:"updated_at"`
}

func (r *ChatSettingsRow) Values() []any {
	return []any{
		r.ChatId,
		r.CollapsePeriod,
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
		SilenceUntil:   m.SilenceUntil,
		UpdatedAt:      m.UpdatedAt,
	}
}
