package logs

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/google/uuid"
)

func (s *service) FindAllByToken(ctx context.Context, token uuid.UUID) ([]*model.Log, error) {
	return s.logsRepository.FindAllByToken(ctx, token)
}
