package last_sent

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *repository) FindByKey(ctx context.Context, lastSentKey *model.LastSentKey) (*model.LastSent, error) {
	q := db.Query{
		Name: "repository.last_sent.find_by_key",
		Sqlizer: s.sq.Select(lastSentTableColumns...).
			From(lastSentTable).
			Where(squirrel.Eq{
				lastSentTableColumnChatId: lastSentKey.ChatId,
				lastSentTableColumnToken:  lastSentKey.Token,
				lastSentTableColumnHash:   lastSentKey.Hash,
			}),
	}

	var out LastSentRow
	if err := s.db.DB().Getx(ctx, &out, q); err != nil {
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	return ToModel(&out), nil
}
