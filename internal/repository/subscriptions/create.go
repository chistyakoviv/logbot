package subscriptions

import (
	"context"
	"fmt"
	"strings"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *repository) Create(ctx context.Context, in *model.Subscription) (*model.Subscription, error) {
	row := FromModel(in)

	q := db.Query{
		Name: "repository.subscriptions.create",
		Sqlizer: r.sq.Insert(subscriptionsTable).
			Columns(subscriptionsTableColumns...).
			Values(row.Values()...).
			Suffix("RETURNING " + strings.Join(subscriptionsTableColumns, ",")),
	}

	var out SubscriptionRow
	if err := r.db.DB().Getx(ctx, &out, q); err != nil {
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	return ToModel(&out), nil
}
