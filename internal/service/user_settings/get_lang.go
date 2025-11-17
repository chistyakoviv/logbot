package user_settings

import (
	"context"
)

func (s *service) GetLang(ctx context.Context, id int64) (string, error) {
	settings, err := s.userSettingsRepository.Find(ctx, id)
	if err != nil {
		return "", err
	}

	return settings.Language(), nil
}
