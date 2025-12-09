package logs

import (
	"context"
	"time"

	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) Create(ctx context.Context, log *model.LogInfo) (*model.Log, error) {
	newLog := &model.Log{
		Data:          log.Data,
		Service:       log.Service,
		ContainerName: log.ContainerName,
		ContainerId:   log.ContainerId,
		NodeId:        log.NodeId,
		Token:         log.Token,
		Hash:          log.Hash,
		CreatedAt:     time.Now(),
	}

	return s.logsRepository.Create(ctx, newLog)
}
