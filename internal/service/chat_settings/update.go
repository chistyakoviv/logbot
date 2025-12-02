package chat_settings

import (
	"context"
	"errors"
	"time"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) Update(ctx context.Context, chatId int64, in *model.ChatSettingsInfo) (*model.ChatSettings, error) {
	chatSettings := &model.ChatSettings{
		ChatId:         chatId,
		CollapsePeriod: in.CollapsePeriod,
		SilenceUntil:   in.SilenceUntil,
		UpdatedAt:      time.Now(),
	}

	_, err := s.chatSettingsRepository.Find(ctx, chatId)
	if err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return nil, err
		}

		return s.chatSettingsRepository.Create(ctx, chatSettings)
	}

	return s.chatSettingsRepository.Update(ctx, chatSettings)
}
