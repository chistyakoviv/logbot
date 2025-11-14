package groups

import (
	"time"

	"github.com/chistyakoviv/logbot/internal/models"
	"github.com/google/uuid"
)

type GroupRow struct {
	ID        uuid.UUID `db:"id"`
	ChatID    int64     `db:"chat_id"`
	Token     string    `db:"token"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *GroupRow) Values() []any {
	return []any{
		r.ID,
		r.ChatID,
		r.Token,
		r.CreatedAt,
	}
}

func ToModel(r *GroupRow) *models.Group {
	if r == nil {
		return nil
	}
	return &models.Group{
		ID:        r.ID,
		ChatID:    r.ChatID,
		Token:     r.Token,
		CreatedAt: r.CreatedAt,
	}
}

// The object exists inside a repository method for a short time, no need to keep a reference.
func FromModel(m *models.Group) GroupRow {
	if m == nil {
		return GroupRow{}
	}
	return GroupRow{
		ID:        m.ID,
		ChatID:    m.ChatID,
		Token:     m.Token,
		CreatedAt: m.CreatedAt,
	}
}
