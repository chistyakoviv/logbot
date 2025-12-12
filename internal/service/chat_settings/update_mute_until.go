package chat_settings

import (
	"context"
	"errors"
	"time"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) UpdateMuteUntil(ctx context.Context, chatId int64, muteUntil time.Time) (*model.ChatSettings, error) {
	_, err := s.chatSettingsRepository.Find(ctx, chatId)
	if err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return nil, err
		}

		return s.chatSettingsRepository.Create(ctx, &model.ChatSettings{
			ChatId:    chatId,
			MuteUntil: muteUntil,
			UpdatedAt: time.Now(),
		})
	}

	return s.chatSettingsRepository.UpdateMuteUntil(ctx, chatId, muteUntil)
}
