package logs

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/google/uuid"
)

type RepositoryInterface interface {
	Create(ctx context.Context, log *model.Log) (*model.Log, error)
	FindAllByToken(ctx context.Context, token uuid.UUID) ([]*model.Log, error)
	Delete(ctx context.Context, id int) error
	DeleteByToken(ctx context.Context, token uuid.UUID) error
	DeleteByHash(ctx context.Context, hash string) error
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
