package subscriptions

import (
	"context"
	"errors"
	"time"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/google/uuid"
)

func (s *service) Subscribe(ctx context.Context, in *model.SubscriptionInfo) (*model.Subscription, error) {
	var out *model.Subscription
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		sub := &model.Subscription{
			Id:          uuid.New(),
			ChatId:      in.ChatId,
			Token:       in.Token,
			ProjectName: in.ProjectName,
			CreatedAt:   time.Now(),
		}
		_, err = s.subscriptionsRepository.Find(ctx, sub.Token, sub.ChatId)
		if err == nil {
			return model.ErrTokenAlreadyExists
		}
		if !errors.Is(err, db.ErrNotFound) {
			return err
		}

		out, err = s.subscriptionsRepository.Create(ctx, sub)

		return err
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}
