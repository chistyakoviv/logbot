package subscriptions

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) FindByTokenAndChat(ctx context.Context, token string, chatId int64) (*model.Subscription, error) {
	return s.subscriptionsRepository.FindByTokenAndChat(ctx, token, chatId)
}
