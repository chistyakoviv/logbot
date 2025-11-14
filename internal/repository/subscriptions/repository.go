package subscriptions

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

type IRepository interface {
	Create(ctx context.Context, in *model.Subscription) (*model.Subscription, error)
	FindByToken(ctx context.Context, token string) (*model.Subscription, error)
	DeleteByToken(ctx context.Context, token string) (*model.Subscription, error)
}

type repository struct {
	db db.Client
	sq sq.StatementBuilderType
}

func NewRepository(db db.Client, sq sq.StatementBuilderType) IRepository {
	return &repository{
		db: db,
		sq: sq,
	}
}
