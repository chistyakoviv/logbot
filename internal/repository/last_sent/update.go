package last_sent

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *repository) Update(ctx context.Context, lastSent *model.LastSent) (*model.LastSent, error) {
	row := FromModel(lastSent)

	q := db.Query{
		Name: "repository.last_sent.update",
		Sqlizer: r.sq.Update(lastSentTable).
			Where(squirrel.Eq{
				lastSentTableColumnChatId: row.ChatId,
				lastSentTableColumnToken:  row.Token,
				lastSentTableColumnHash:   row.Hash,
			}).
			Set(lastSentTableColumnUpdatedAt, row.UpdatedAt).
			Suffix("RETURNING " + strings.Join(lastSentTableColumns, ",")),
	}

	var out LastSentRow
	if err := r.db.DB().Getx(ctx, &out, q); err != nil {
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	return ToModel(&out), nil
}
