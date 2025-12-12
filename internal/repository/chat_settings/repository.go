package chat_settings

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

type RepositoryInterface interface {
	Create(ctx context.Context, in *model.ChatSettings) (*model.ChatSettings, error)
	UpdateCollapsePeriod(ctx context.Context, chatId int64, period time.Duration) (*model.ChatSettings, error)
	UpdateMuteUntil(ctx context.Context, chatId int64, muteUntil time.Time) (*model.ChatSettings, error)
	Find(ctx context.Context, chatId int64) (*model.ChatSettings, error)
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
