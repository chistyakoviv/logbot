package user_settings

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/chistyakoviv/logbot/internal/repository/user_settings"
)

type IService interface {
	Find(ctx context.Context, id int64) (*model.UserSettings, error)
	Update(ctx context.Context, id int64, in *model.UserSettingsInfo) (*model.UserSettings, error)
}

type service struct {
	userSettingsRepository user_settings.IRepository
	txManager              db.TxManager
}

func NewService(userSettingsRepository user_settings.IRepository, txManager db.TxManager) IService {
	return &service{
		userSettingsRepository: userSettingsRepository,
		txManager:              txManager,
	}
}
