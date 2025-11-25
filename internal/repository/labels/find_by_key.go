package labels

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *repository) FindByKey(ctx context.Context, in *model.LabelKey) (*model.Label, error) {
	q := db.Query{
		Name: "repository.labels.find_by_key",
		Sqlizer: r.sq.Select(labelsTableColumns...).
			From(labelsTable).
			Where(sq.Eq{
				labelsTableColumnUserId: in.UserId,
				labelsTableColumnChatId: in.ChatId,
			}),
	}

	var row LabelsRow
	if err := r.db.DB().Getx(ctx, &row, q); err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	return ToModel(&row), nil
}
