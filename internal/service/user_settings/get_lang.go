package user_settings

import (
	"context"
	"errors"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) GetLang(ctx context.Context, id int64) (string, error) {
	settings, err := s.userSettingsRepository.Find(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			// Create a model with default values and return the default language
			settings := &model.UserSettings{}
			return settings.Language(), nil
		}
		return "", err
	}

	return settings.Language(), nil
}
