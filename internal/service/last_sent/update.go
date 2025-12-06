package last_sent

import (
	"context"
	"errors"
	"time"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) Update(ctx context.Context, lastSentKey *model.LastSentKey) (*model.LastSent, error) {
	newLastSent := &model.LastSent{
		ChatId:    lastSentKey.ChatId,
		Token:     lastSentKey.Token,
		Hash:      lastSentKey.Hash,
		UpdatedAt: time.Now(),
	}

	_, err := s.lastSentRepository.FindByKey(ctx, lastSentKey)
	if err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return nil, err
		}
		return s.lastSentRepository.Create(ctx, newLastSent)
	}

	return s.lastSentRepository.Update(ctx, newLastSent)
}
