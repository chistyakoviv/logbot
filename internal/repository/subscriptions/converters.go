package subscriptions

import (
	"time"

	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/google/uuid"
)

type SubscriptionRow struct {
	Id        uuid.UUID `db:"id"`
	ChatId    int64     `db:"chat_id"`
	Token     string    `db:"token"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *SubscriptionRow) Values() []any {
	return []any{
		r.Id,
		r.ChatId,
		r.Token,
		r.CreatedAt,
	}
}

func ToModel(r *SubscriptionRow) *model.Subscription {
	if r == nil {
		return nil
	}
	return &model.Subscription{
		Id:        r.Id,
		ChatId:    r.ChatId,
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
		Id:        m.Id,
		ChatId:    m.ChatId,
		Token:     m.Token,
		CreatedAt: m.CreatedAt,
	}
}
