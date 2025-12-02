package user_settings

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/chistyakoviv/logbot/internal/repository/user_settings"
)

type ServiceInterface interface {
	Find(ctx context.Context, userId int64) (*model.UserSettings, error)
	Update(ctx context.Context, userId int64, in *model.UserSettingsInfo) (*model.UserSettings, error)
	GetLang(ctx context.Context, userId int64) (string, error)
}

type service struct {
	userSettingsRepository user_settings.RepositoryInterface
	txManager              db.TxManager
}

func NewService(userSettingsRepository user_settings.RepositoryInterface, txManager db.TxManager) ServiceInterface {
	return &service{
		userSettingsRepository: userSettingsRepository,
		txManager:              txManager,
	}
}
