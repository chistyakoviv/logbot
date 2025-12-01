package labels

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *repository) FindAllByChat(ctx context.Context, chatId int64) ([]*model.Label, error) {
	q := db.Query{
		Name: "repository.labels.find_all_by_chat",
		Sqlizer: r.sq.Select(labelsTableColumns...).
			From(labelsTable).
			Where(sq.Eq{
				labelsTableColumnChatId: chatId,
			}),
	}

	var rows []LabelsRow
	if err := r.db.DB().Selectx(ctx, &rows, q); err != nil {
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	out := make([]*model.Label, 0, len(rows))
	for _, v := range rows {
		out = append(out, ToModel(&v))
	}

	return out, nil
}
