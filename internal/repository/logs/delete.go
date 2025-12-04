package logs

import (
	"context"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
)

func (r *repository) Delete(ctx context.Context, id int) error {
	q := db.Query{
		Name: "repository.logs.delete",
		Sqlizer: r.sq.Delete(logsTable).
			Where(sq.Eq{logsTableColumnId: id}).
			Suffix("RETURNING " + strings.Join(logsTableColumns, ",")),
	}

	return r.db.DB().Getx(ctx, nil, q)
}
