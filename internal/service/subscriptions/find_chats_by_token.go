package subscriptions

import "context"

func (s *service) FindChatsByToken(ctx context.Context, token string) ([]int64, error) {
	return s.subscriptionsRepository.FindChatsByToken(ctx, token)
}
