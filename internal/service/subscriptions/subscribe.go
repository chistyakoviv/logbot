package subscriptions

import (
	"context"
	"time"

	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/google/uuid"
)

func (s *service) Subscribe(ctx context.Context, in *model.SubscriptionInfo) (*model.Subscription, error) {
	sub := &model.Subscription{
		ID:        uuid.New(),
		ChatID:    in.ChatID,
		Token:     in.Token,
		CreatedAt: time.Now(),
	}
	return s.subscriptionsRepository.Create(ctx, sub)
}
