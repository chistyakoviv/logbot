package subscriptions

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

type RepositoryInterface interface {
	Create(ctx context.Context, in *model.Subscription) (*model.Subscription, error)
	Find(ctx context.Context, token string, chatId int64) (*model.Subscription, error)
	FindByChatId(ctx context.Context, chatId int64) ([]*model.Subscription, error)
	FindByToken(ctx context.Context, token string) ([]*model.Subscription, error)
	Exists(ctx context.Context, token string) (bool, error)
	Delete(ctx context.Context, token string, chatId int64) (*model.Subscription, error)
}

type repository struct {
	db db.Client
	sq sq.StatementBuilderType
}

func NewRepository(db db.Client, sq sq.StatementBuilderType) RepositoryInterface {
	return &repository{
		db: db,
		sq: sq,
	}
}
