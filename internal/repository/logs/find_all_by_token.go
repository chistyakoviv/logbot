package logs

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/google/uuid"
)

func (r *repository) FindAllByToken(ctx context.Context, token uuid.UUID) ([]*model.Log, error) {
	q := db.Query{
		Name: "repository.logs.find_all_by_token",
		Sqlizer: r.sq.Select(logsTableColumns...).
			From(logsTable).
			Where(sq.Eq{logsTableColumnToken: token}),
	}

	var rows []LogsRow
	if err := r.db.DB().Selectx(ctx, &rows, q); err != nil {
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	out := make([]*model.Log, 0, len(rows))
	for _, v := range rows {
		out = append(out, ToModel(&v))
	}

	return out, nil
}
