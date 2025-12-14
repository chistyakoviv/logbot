package logs

import (
	"context"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
)

func (r *repository) DeleteOlderThan(ctx context.Context, timestamp time.Time) error {
	q := db.Query{
		Name: "repository.logs.delete_older_than",
		Sqlizer: r.sq.Delete(logsTable).
			Where(sq.Lt{logsTableColumnCreatedAt: timestamp}).
			Suffix("RETURNING " + strings.Join(logsTableColumns, ",")),
	}

	return r.db.DB().Getx(ctx, nil, q)
}
