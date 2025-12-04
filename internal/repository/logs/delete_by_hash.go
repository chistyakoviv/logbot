package logs

import (
	"context"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
)

func (r *repository) DeleteByHash(ctx context.Context, hash string) error {
	q := db.Query{
		Name: "repository.logs.delete_by_hash",
		Sqlizer: r.sq.Delete(logsTable).
			Where(sq.Eq{logsTableColumnHash: hash}).
			Suffix("RETURNING " + strings.Join(logsTableColumns, ",")),
	}

	return r.db.DB().Getx(ctx, nil, q)
}
