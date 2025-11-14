package subscriptions

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) Unsubscribe(ctx context.Context, token string) (*model.Subscription, error) {
	return s.subscriptionsRepository.DeleteByToken(ctx, token)
}
