package user_settings

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) Update(ctx context.Context, id int64, in *model.UserSettingsInfo) (*model.UserSettings, error) {
	userSettings := &model.UserSettings{
		UserId: id,
		Lang:   in.Lang,
	}
	return s.userSettingsRepository.UpdateOrCreate(ctx, userSettings)
}
