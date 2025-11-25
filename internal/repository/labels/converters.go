package labels

import (
	"time"

	"github.com/chistyakoviv/logbot/internal/model"
)

type LabelsRow struct {
	ChatId    int64     `db:"chat_id"`
	UserId    int64     `db:"user_id"`
	Labels    []string  `db:"labels"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r *LabelsRow) Values() []any {
	return []any{
		r.ChatId,
		r.UserId,
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
		UserId:    r.UserId,
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
		UserId:    m.UserId,
		Labels:    m.Labels,
		UpdatedAt: m.UpdatedAt,
	}
}
