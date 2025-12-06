package last_sent

import (
	"context"
	"time"

	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) Get(ctx context.Context, lastSentKey *model.LastSentKey) time.Time {
	t, _ := s.lastSentRepository.LastSent(ctx, lastSentKey)
	return t
}
