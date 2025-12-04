package logs

import (
	"time"

	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/google/uuid"
)

type LogsRow struct {
	Id        int       `db:"id"`
	Token     uuid.UUID `db:"token"`
	Data      string    `db:"data"`
	Label     string    `db:"label"`
	Hash      string    `db:"hash"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *LogsRow) Values() []any {
	return []any{
		r.Id,
		r.Token,
		r.Data,
		r.Label,
		r.Hash,
		r.CreatedAt,
	}
}

func ToModel(r *LogsRow) *model.Log {
	if r == nil {
		return nil
	}
	return &model.Log{
		Id:        r.Id,
		Token:     r.Token,
		Data:      r.Data,
		Label:     r.Label,
		Hash:      r.Hash,
		CreatedAt: r.CreatedAt,
	}
}

func FromModel(m *model.Log) LogsRow {
	if m == nil {
		return LogsRow{}
	}
	return LogsRow{
		Id:        m.Id,
		Token:     m.Token,
		Data:      m.Data,
		Label:     m.Label,
		Hash:      m.Hash,
		CreatedAt: m.CreatedAt,
	}
}
