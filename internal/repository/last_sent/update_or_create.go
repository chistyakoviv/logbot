package last_sent

import (
	"context"
	"errors"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *repository) UpdateOrCreate(ctx context.Context, lastSent *model.LastSent) (*model.LastSent, error) {
	lastSentKey := &model.LastSentKey{
		ChatId: lastSent.ChatId,
		Token:  lastSent.Token,
		Hash:   lastSent.Hash,
	}

	_, err := r.FindByKey(ctx, lastSentKey)
	if err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return nil, err
		}

		lastSent, err = r.Create(ctx, lastSent)
		if err != nil {
			return nil, err
		}
		return lastSent, nil
	}

	lastSent, err = r.Update(ctx, lastSent)
	if err != nil {
		return nil, err
	}
	return lastSent, nil
}
