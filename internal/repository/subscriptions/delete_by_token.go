package subscriptions

import (
	"context"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *repository) DeleteByToken(ctx context.Context, token string) (*model.Subscription, error) {
	q := db.Query{
		Name: "repository.subscriptions.delete_by_token",
		Sqlizer: r.sq.Delete(subscriptionsTable).
			Where(sq.Eq{subscriptionsTableColumnToken: token}).
			Suffix("RETURNING " + strings.Join(subscriptionsTableColumns, ",")),
	}

	var row SubscriptionRow
	if err := r.db.DB().Getx(ctx, &row, q); err != nil {
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	return ToModel(&row), nil
}
