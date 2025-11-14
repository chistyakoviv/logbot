package groups

import (
	"context"
	"fmt"
	"strings"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/models"
)

func (r *repository) Create(ctx context.Context, in *models.Group) (*models.Group, error) {
	row := FromModel(in)

	q := db.Query{
		Name: "repository.groups.create",
		Sqlizer: r.sq.Insert(groupsTable).
			Columns(groupsTableColumns...).
			Values(row.Values()...).
			Suffix("RETURNING " + strings.Join(groupsTableColumns, ",")),
	}

	var out GroupRow
	if err := r.db.DB().Getx(ctx, &out, q); err != nil {
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	return ToModel(&out), nil
}
