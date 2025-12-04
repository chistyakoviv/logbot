package logs

import (
	"context"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/google/uuid"
)

func (r *repository) DeleteByToken(ctx context.Context, token uuid.UUID) error {
	q := db.Query{
		Name: "repository.logs.delete_by_token",
		Sqlizer: r.sq.Delete(logsTable).
			Where(sq.Eq{logsTableColumnToken: token}).
			Suffix("RETURNING " + strings.Join(logsTableColumns, ",")),
	}

	return r.db.DB().Getx(ctx, nil, q)
}
