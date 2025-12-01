package labels

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) FindAllByChat(ctx context.Context, chatId int64) ([]*model.Label, error) {
	return s.labelsRepository.FindAllByChat(ctx, chatId)
}
