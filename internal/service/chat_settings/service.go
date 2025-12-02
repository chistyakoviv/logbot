package chat_settings

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/chistyakoviv/logbot/internal/repository/chat_settings"
)

type ServiceInterface interface {
	Update(ctx context.Context, chatId int64, in *model.ChatSettingsInfo) (*model.ChatSettings, error)
	Find(ctx context.Context, chatId int64) (*model.ChatSettings, error)
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
