package subscriptions

import (
	"time"

	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/google/uuid"
)

type SubscriptionRow struct {
	ID        uuid.UUID `db:"id"`
	ChatID    int64     `db:"chat_id"`
	Token     string    `db:"token"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *SubscriptionRow) Values() []any {
	return []any{
		r.ID,
		r.ChatID,
		r.Token,
		r.CreatedAt,
	}
}

func ToModel(r *SubscriptionRow) *model.Subscription {
	if r == nil {
		return nil
	}
	return &model.Subscription{
		ID:        r.ID,
		ChatID:    r.ChatID,
		Token:     r.Token,
		CreatedAt: r.CreatedAt,
	}
}

// The object exists inside a repository method for a short time, no need to keep a reference.
func FromModel(m *model.Subscription) SubscriptionRow {
	if m == nil {
		return SubscriptionRow{}
	}
	return SubscriptionRow{
		ID:        m.ID,
		ChatID:    m.ChatID,
		Token:     m.Token,
		CreatedAt: m.CreatedAt,
	}
}
