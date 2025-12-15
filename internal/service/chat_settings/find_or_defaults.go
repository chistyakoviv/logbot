package chat_settings

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) FindOrDefaults(ctx context.Context, chatId int64) (*model.ChatSettings, error) {
	return s.chatSettingsRepository.FindOrDefaults(ctx, chatId)
}
