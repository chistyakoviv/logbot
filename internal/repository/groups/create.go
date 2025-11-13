package groups

import (
	"context"
	"strings"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/models"
)

func (r *Repository) Create(ctx context.Context, in *models.Group) (*models.Group, error) {
	row := FromModel(in)

	builder := r.sq.Insert(groupsTable).
		Columns(groupsTableColumns...).
		Values(row.Values()...).
		Suffix("RETURNING " + strings.Join(groupsTableColumns, ","))

	q := db.Query{
		Name:    "repository.groups.create",
		Sqlizer: builder,
	}

	var out GroupRow
	if err := r.db.Getx(ctx, &out, q); err != nil {
		return nil, err
	}

	return ToModel(&out), nil
}
