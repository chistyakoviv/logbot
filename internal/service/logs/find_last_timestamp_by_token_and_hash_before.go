package logs

import (
	"context"
	"time"
)

func (s *service) FindLastTimestampByTokenAndHashBefore(
	ctx context.Context,
	token string,
	hash string,
	before time.Time,
) (time.Time, error) {
	return s.logsRepository.FindLastTimestampByTokenAndHashBefore(ctx, token, hash, before)
}
