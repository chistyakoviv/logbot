package last_sent

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *repository) LastSentOrDefault(ctx context.Context, lastSentKey *model.LastSentKey) (time.Time, error) {
	q := db.Query{
		Name: "repository.last_sent.last_sent",
		Sqlizer: r.sq.Select(lastSentTableColumnUpdatedAt).
			From(lastSentTable).
			Where(squirrel.Eq{
				lastSentTableColumnChatId: lastSentKey.ChatId,
				lastSentTableColumnToken:  lastSentKey.Token,
				lastSentTableColumnHash:   lastSentKey.Hash,
			}),
	}

	var out LastSentRow
	if err := r.db.DB().Getx(ctx, &out, q); err != nil {
		return time.Time{}, fmt.Errorf("%s: %w", q.Name, err)
	}

	return out.UpdatedAt, nil
}
