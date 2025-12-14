package logs

import (
	"context"
	"time"
)

func (s *service) DeleteOlderThan(ctx context.Context, timestamp time.Time) error {
	return s.logsRepository.DeleteOlderThan(ctx, timestamp)
}
