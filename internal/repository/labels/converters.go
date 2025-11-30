package labels

import (
	"time"

	"github.com/chistyakoviv/logbot/internal/model"
)

type LabelsRow struct {
	ChatId    int64     `db:"chat_id"`
	Username  string    `db:"username"`
	Labels    []string  `db:"labels"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r *LabelsRow) Values() []any {
	return []any{
		r.ChatId,
		r.Username,
		r.Labels,
		r.UpdatedAt,
	}
}

func ToModel(r *LabelsRow) *model.Label {
	if r == nil {
		return nil
	}
	return &model.Label{
		ChatId:    r.ChatId,
		Username:  r.Username,
		Labels:    r.Labels,
		UpdatedAt: r.UpdatedAt,
	}
}

func FromModel(m *model.Label) *LabelsRow {
	if m == nil {
		return nil
	}
	return &LabelsRow{
		ChatId:    m.ChatId,
		Username:  m.Username,
		Labels:    m.Labels,
		UpdatedAt: m.UpdatedAt,
	}
}
