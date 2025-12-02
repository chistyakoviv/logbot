package user_settings

import (
	"context"
	"errors"
	"time"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

// TODO: execute in transaction to avoid concurrent creation
func (s *service) Update(ctx context.Context, userId int64, in *model.UserSettingsInfo) (*model.UserSettings, error) {
	userSettings := &model.UserSettings{
		UserId:    userId,
		Username:  in.Username,
		Lang:      in.Lang,
		UpdatedAt: time.Now(),
	}

	_, err := s.userSettingsRepository.Find(ctx, userId)
	if err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return nil, err
		}

		return s.userSettingsRepository.Create(ctx, userSettings)
	}

	return s.userSettingsRepository.Update(ctx, userSettings)
}
