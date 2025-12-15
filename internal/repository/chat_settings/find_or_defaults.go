package chat_settings

import (
	"context"
	"errors"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *repository) FindOrDefaults(ctx context.Context, chatId int64) (*model.ChatSettings, error) {
	settings, err := r.Find(ctx, chatId)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return &model.ChatSettings{ChatId: chatId}, nil
		}
		return nil, err
	}
	return settings, nil
}
