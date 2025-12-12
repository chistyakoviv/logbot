package labels

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) FindByChatIdAndLabel(ctx context.Context, chatId int64, label string) ([]*model.Label, error) {
	return s.labelsRepository.FindByChatIdAndLabel(ctx, chatId, label)
}
