package subscriptions

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *repository) Find(ctx context.Context, token string, chatId int64) (*model.Subscription, error) {
	q := db.Query{
		Name: "repository.subscriptions.find_by_token_and_chat",
		Sqlizer: r.sq.Select(subscriptionsTableColumns...).
			From(subscriptionsTable).
			Where(sq.Eq{
				subscriptionsTableColumnToken:  token,
				subscriptionsTableColumnChatId: chatId,
			}),
	}

	var row SubscriptionRow
	if err := r.db.DB().Getx(ctx, &row, q); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	return ToModel(&row), nil
}
