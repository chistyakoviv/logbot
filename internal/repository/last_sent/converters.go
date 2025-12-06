package last_sent

import (
	"time"

	"github.com/chistyakoviv/logbot/internal/model"
)

type LastSentRow struct {
	ChatId    int64     `db:"chat_id"`
	Token     string    `db:"token"`
	Hash      string    `db:"hash"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r *LastSentRow) Values() []any {
	return []any{
		r.ChatId,
		r.Token,
		r.Hash,
		r.UpdatedAt,
	}
}

func ToModel(r *LastSentRow) *model.LastSent {
	if r == nil {
		return nil
	}
	return &model.LastSent{
		ChatId:    r.ChatId,
		Token:     r.Token,
		Hash:      r.Hash,
		UpdatedAt: r.UpdatedAt,
	}
}

func FromModel(m *model.LastSent) LastSentRow {
	if m == nil {
		return LastSentRow{}
	}
	return LastSentRow{
		ChatId:    m.ChatId,
		Token:     m.Token,
		Hash:      m.Hash,
		UpdatedAt: m.UpdatedAt,
	}
}
