package subscriptions

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
)

func (r *repository) Exists(ctx context.Context, token string) (bool, error) {
	q := db.Query{
		Name: "repository.subscriptions.exists",
		Sqlizer: r.sq.Select("1").
			From(subscriptionsTable).
			Where(sq.Eq{subscriptionsTableColumnToken: token}),
	}

	var row SubscriptionRow
	if err := r.db.DB().Getx(ctx, &row, q); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("%s: %w", q.Name, err)
	}

	return true, nil
}
