package commands

import (
	"time"

	"github.com/chistyakoviv/logbot/internal/model"
)

type CommandRow struct {
	Name      string                 `db:"name"`
	UserId    int64                  `db:"user_id"`
	ChatId    int64                  `db:"chat_id"`
	Stage     int                    `db:"stage"`
	Data      map[string]interface{} `db:"data"`
	UpdatedAt time.Time              `db:"updated_at"`
}

func (r *CommandRow) Values() []any {
	return []any{
		r.Name,
		r.UserId,
		r.ChatId,
		r.Stage,
		r.Data,
		r.UpdatedAt,
	}
}

func ToModel(r *CommandRow) *model.Command {
	if r == nil {
		return nil
	}
	return &model.Command{
		Name:      r.Name,
		UserId:    r.UserId,
		ChatId:    r.ChatId,
		Stage:     r.Stage,
		Data:      r.Data,
		UpdatedAt: r.UpdatedAt,
	}
}

func FromModel(m *model.Command) CommandRow {
	if m == nil {
		return CommandRow{}
	}
	return CommandRow{
		Name:      m.Name,
		UserId:    m.UserId,
		ChatId:    m.ChatId,
		Stage:     m.Stage,
		Data:      m.Data,
		UpdatedAt: m.UpdatedAt,
	}
}
