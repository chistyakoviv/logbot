package labels

import (
	"context"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *repository) Update(ctx context.Context, in *model.Label) (*model.Label, error) {
	row := FromModel(in)

	q := db.Query{
		Name: "repository.labels.update",
		Sqlizer: r.sq.Update(labelsTable).
			Where(sq.Eq{
				labelsTableColumnUserId: row.UserId,
				labelsTableColumnChatId: row.ChatId,
			}).
			SetMap(
				map[string]interface{}{
					labelsTableColumnLabels:    row.Labels,
					labelsTableColumnUpdatedAt: row.UpdatedAt,
				},
			).
			Suffix("RETURNING " + strings.Join(labelsTableColumns, ",")),
	}

	var out LabelsRow
	if err := r.db.DB().Getx(ctx, &out, q); err != nil {
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	return ToModel(&out), nil
}
