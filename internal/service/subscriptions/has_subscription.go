package subscriptions

import (
	"context"
)

func (s *service) HasSubscription(ctx context.Context, token string) (bool, error) {
	return s.subscriptionsRepository.Exists(ctx, token)
}
