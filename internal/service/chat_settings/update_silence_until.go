package chat_settings

import (
	"context"
	"errors"
	"time"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) UpdateSilenceUntil(ctx context.Context, chatId int64, silenceUntil time.Time) (*model.ChatSettings, error) {
	_, err := s.chatSettingsRepository.Find(ctx, chatId)
	if err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return nil, err
		}

		return s.chatSettingsRepository.Create(ctx, &model.ChatSettings{
			ChatId:       chatId,
			SilenceUntil: silenceUntil,
			UpdatedAt:    time.Now(),
		})
	}
	return s.chatSettingsRepository.UpdateSilenceUntil(ctx, chatId, silenceUntil)
}
