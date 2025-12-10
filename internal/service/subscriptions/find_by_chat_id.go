package subscriptions

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) FindByChatId(ctx context.Context, chatId int64) ([]*model.Subscription, error) {
	return s.subscriptionsRepository.FindByChatId(ctx, chatId)
}
