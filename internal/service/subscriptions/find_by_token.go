package subscriptions

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) FindByToken(ctx context.Context, token string) (*model.Subscription, error) {
	return s.subscriptionsRepository.FindByToken(ctx, token)
}
