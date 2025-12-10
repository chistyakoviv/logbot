package subscriptions

import (
	"context"
	"time"

	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/google/uuid"
)

func (s *service) Subscribe(ctx context.Context, in *model.SubscriptionInfo) (*model.Subscription, error) {
	sub := &model.Subscription{
		Id:          uuid.New(),
		ChatId:      in.ChatId,
		Token:       in.Token,
		ProjectName: in.ProjectName,
		CreatedAt:   time.Now(),
	}
	out, err := s.subscriptionsRepository.Create(ctx, sub)

	if err != nil {
		return nil, err
	}

	return out, nil
}
