package subscriptions

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/chistyakoviv/logbot/internal/repository/subscriptions"
)

type ServiceInterface interface {
	Subscribe(ctx context.Context, in *model.SubscriptionInfo) (*model.Subscription, error)
	Find(ctx context.Context, token string, chatId int64) (*model.Subscription, error)
	FindByChatId(ctx context.Context, chatId int64) ([]*model.Subscription, error)
	HasSubscription(ctx context.Context, token string) (bool, error)
	FindByToken(ctx context.Context, token string) ([]*model.Subscription, error)
	Unsubscribe(ctx context.Context, token string, chatId int64) (*model.Subscription, error)
}

type service struct {
	subscriptionsRepository subscriptions.RepositoryInterface
	txManager               db.TxManager
}

func NewService(subscriptionsRepository subscriptions.RepositoryInterface, txManager db.TxManager) ServiceInterface {
	return &service{
		subscriptionsRepository: subscriptionsRepository,
		txManager:               txManager,
	}
}
