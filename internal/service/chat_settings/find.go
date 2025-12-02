package chat_settings

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) Find(ctx context.Context, chatId int64) (*model.ChatSettings, error) {
	return s.chatSettingsRepository.Find(ctx, chatId)
}
