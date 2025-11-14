package groups

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/models"
)

type IRepository interface {
	Create(ctx context.Context, in *models.Group) (*models.Group, error)
	FindByToken(ctx context.Context, token string) (*models.Group, error)
	DeleteByToken(ctx context.Context, token string) (*models.Group, error)
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
