package user_settings

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

type IRepository interface {
	Create(ctx context.Context, in *model.UserSettings) (*model.UserSettings, error)
	Update(ctx context.Context, in *model.UserSettings) (*model.UserSettings, error)
	Find(ctx context.Context, id int64) (*model.UserSettings, error)
}

type respository struct {
	db db.Client
	sq sq.StatementBuilderType
}

func NewRepository(db db.Client, sq sq.StatementBuilderType) IRepository {
	return &respository{
		db: db,
		sq: sq,
	}
}
