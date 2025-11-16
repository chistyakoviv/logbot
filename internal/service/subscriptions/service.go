package subscriptions

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/chistyakoviv/logbot/internal/repository/subscriptions"
)

type IService interface {
	Subscribe(ctx context.Context, in *model.SubscriptionInfo) (*model.Subscription, error)
	FindByTokenAndChat(ctx context.Context, token string, chatId int64) (*model.Subscription, error)
	Unsubscribe(ctx context.Context, token string, chatId int64) (*model.Subscription, error)
}

type service struct {
	subscriptionsRepository subscriptions.IRepository
	txManager               db.TxManager
}

func NewService(subscriptionsRepository subscriptions.IRepository, txManager db.TxManager) IService {
	return &service{
		subscriptionsRepository: subscriptionsRepository,
		txManager:               txManager,
	}
}
