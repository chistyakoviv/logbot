package logs

import (
	"context"
	"time"

	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) Create(ctx context.Context, log *model.LogInfo) (*model.Log, error) {
	newLog := &model.Log{
		Data:      log.Data,
		Label:     log.Label,
		Token:     log.Token,
		Hash:      s.logHasher.Hash(log.Data),
		CreatedAt: time.Now(),
	}

	return s.logsRepository.Create(ctx, newLog)
}
