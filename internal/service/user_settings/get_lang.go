package user_settings

import (
	"context"
	"errors"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) GetLang(ctx context.Context, userId int64) (string, error) {
	settings, err := s.userSettingsRepository.Find(ctx, userId)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			// Create a model with default values and return the default language
			settings := &model.UserSettings{}
			return settings.Language(), db.ErrNotFound
		}
		return "", err
	}

	return settings.Language(), nil
}
