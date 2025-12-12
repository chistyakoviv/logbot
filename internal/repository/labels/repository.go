package labels

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

type RepositoryInterface interface {
	FindByKey(ctx context.Context, in *model.LabelKey) (*model.Label, error)
	FindByChatIdAndLabel(ctx context.Context, chatId int64, label string) ([]*model.Label, error)
	FindAllByChat(ctx context.Context, chatId int64) ([]*model.Label, error)
	Create(ctx context.Context, in *model.Label) (*model.Label, error)
	Update(ctx context.Context, in *model.Label) (*model.Label, error)
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
