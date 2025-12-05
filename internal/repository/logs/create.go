package logs

import (
	"context"
	"fmt"
	"strings"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *repository) Create(ctx context.Context, log *model.Log) (*model.Log, error) {
	row := FromModel(log)

	q := db.Query{
		Name: "repository.logs.create",
		Sqlizer: r.sq.Insert(logsTable).
			Columns(logsTableInsertableColumns...).
			Values(row.Values()...).
			Suffix("RETURNING " + strings.Join(logsTableColumns, ",")),
	}

	var out LogsRow
	if err := r.db.DB().Getx(ctx, &out, q); err != nil {
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	return ToModel(&out), nil
}
