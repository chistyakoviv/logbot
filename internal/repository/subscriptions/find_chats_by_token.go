package subscriptions

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
)

func (r *repository) FindChatsByToken(ctx context.Context, token string) ([]int64, error) {
	q := db.Query{
		Name: "repository.subscriptions.find_chats_by_token",
		Sqlizer: r.sq.Select(subscriptionsTableColumnChatId).
			From(subscriptionsTable).
			Where(sq.Eq{subscriptionsTableColumnToken: token}),
	}

	var rows []int64
	if err := r.db.DB().Selectx(ctx, &rows, q); err != nil {
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	return rows, nil
}
