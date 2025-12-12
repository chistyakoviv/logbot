package chat_settings

import (
	"context"
	"errors"
	"time"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) UpdateCollapsePeriod(ctx context.Context, chatId int64, period time.Duration) (*model.ChatSettings, error) {
	_, err := s.chatSettingsRepository.Find(ctx, chatId)
	if err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return nil, err
		}

		return s.chatSettingsRepository.Create(ctx, &model.ChatSettings{
			ChatId:         chatId,
			CollapsePeriod: period,
			UpdatedAt:      time.Now(),
		})
	}

	return s.chatSettingsRepository.UpdateCollapsePeriod(ctx, chatId, period)
}
