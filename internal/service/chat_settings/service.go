package chat_settings

import (
	"context"
	"time"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/chistyakoviv/logbot/internal/repository/chat_settings"
)

type ServiceInterface interface {
	UpdateCollapsePeriod(ctx context.Context, chatId int64, period time.Duration) (*model.ChatSettings, error)
	UpdateMuteUntil(ctx context.Context, chatId int64, muteUntil time.Time) (*model.ChatSettings, error)
	UpdateSilenceUntil(ctx context.Context, chatId int64, silenceUntil time.Time) (*model.ChatSettings, error)
	Find(ctx context.Context, chatId int64) (*model.ChatSettings, error)
	FindOrDefaults(ctx context.Context, chatId int64) (*model.ChatSettings, error)
}

type service struct {
	chatSettingsRepository chat_settings.RepositoryInterface
	txManager              db.TxManager
}

func NewService(chatSettingsRepository chat_settings.RepositoryInterface, txManager db.TxManager) ServiceInterface {
	return &service{
		chatSettingsRepository: chatSettingsRepository,
		txManager:              txManager,
	}
}
