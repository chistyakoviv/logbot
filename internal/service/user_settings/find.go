package user_settings

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) Find(ctx context.Context, userId int64) (*model.UserSettings, error) {
	return s.userSettingsRepository.Find(ctx, userId)
}
