package subscriptions

import (
	"context"
	"errors"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/jackc/pgx/v5"
)

func (r *repository) FindByToken(ctx context.Context, token string) (*model.Subscription, error) {
	if token == "" {
		return nil, errors.New("token is empty")
	}

	q := db.Query{
		Name: "repository.subscriptions.find_by_token",
		Sqlizer: r.sq.Select(subscriptionsTableColumns...).
			From(subscriptionsTable).
			Where(sq.Eq{subscriptionsTableColumnToken: token}).
			Suffix("RETURNING " + strings.Join(subscriptionsTableColumns, ",")),
	}

	var row SubscriptionRow
	if err := r.db.DB().Getx(ctx, &row, q); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, db.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	return ToModel(&row), nil
}
