package last_sent

import (
	"context"
	"fmt"
	"strings"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *repository) Create(ctx context.Context, lastSent *model.LastSent) (*model.LastSent, error) {
	row := FromModel(lastSent)

	q := db.Query{
		Name: "repository.last_sent.create",
		Sqlizer: r.sq.Insert(lastSentTable).
			Columns(lastSentTableColumns...).
			Values(row.Values()...).
			Suffix("RETURNING " + strings.Join(lastSentTableColumns, ",")),
	}

	var out LastSentRow
	if err := r.db.DB().Getx(ctx, &out, q); err != nil {
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	return ToModel(&out), nil
}
