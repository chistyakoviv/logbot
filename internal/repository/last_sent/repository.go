package last_sent

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

type RepositoryInterface interface {
	Create(ctx context.Context, lastSent *model.LastSent) (*model.LastSent, error)
	Update(ctx context.Context, lastSent *model.LastSent) (*model.LastSent, error)
	LastSent(ctx context.Context, lastSentKey *model.LastSentKey) (time.Time, error)
	FindByKey(ctx context.Context, lastSentKey *model.LastSentKey) (*model.LastSent, error)
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
