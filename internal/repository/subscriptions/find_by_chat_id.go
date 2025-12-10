package subscriptions

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *repository) FindByChatId(ctx context.Context, chatId int64) ([]*model.Subscription, error) {
	q := db.Query{
		Name: "repository.subscriptions.find_by_chat_id",
		Sqlizer: r.sq.Select(subscriptionsTableColumns...).
			From(subscriptionsTable).
			Where(sq.Eq{
				subscriptionsTableColumnChatId: chatId,
			}),
	}

	var rows []SubscriptionRow
	if err := r.db.DB().Selectx(ctx, &rows, q); err != nil {
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	out := make([]*model.Subscription, 0, len(rows))
	for _, v := range rows {
		out = append(out, ToModel(&v))
	}

	return out, nil
}
