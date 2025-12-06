package logs

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
)

func (r *repository) FindLastTimestampByTokenAndHashBefore(
	ctx context.Context,
	token string,
	hash string,
	before time.Time,
) (time.Time, error) {
	q := db.Query{
		Name: "repository.logs.find_last_timestamp_by_token_and_hash_before",
		Sqlizer: r.sq.Select(logsTableColumnCreatedAt).
			From(logsTable).
			Where(sq.Eq{logsTableColumnToken: token}).
			Where(sq.Eq{logsTableColumnHash: hash}).
			Where(sq.Lt{logsTableColumnCreatedAt: before}).
			OrderBy(logsTableColumnCreatedAt + " DESC").
			Limit(1),
	}

	var row LogsRow
	if err := r.db.DB().Getx(ctx, &row, q); err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return time.Time{}, err
		}
		return time.Time{}, fmt.Errorf("%s: %w", q.Name, err)
	}

	return row.CreatedAt, nil
}
