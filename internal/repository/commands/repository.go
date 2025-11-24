package commands

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

type RepositoryInterface interface {
	Create(ctx context.Context, in *model.Command) (*model.Command, error)
	Update(ctx context.Context, in *model.Command) (*model.Command, error)
	FindByKey(ctx context.Context, in *model.CommandKey) (*model.Command, error)
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
