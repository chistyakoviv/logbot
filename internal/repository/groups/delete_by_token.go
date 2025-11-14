package groups

import (
	"context"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/models"
)

func (r *repository) DeleteByToken(ctx context.Context, token string) (*models.Group, error) {
	q := db.Query{
		Name: "repository.groups.delete_by_token",
		Sqlizer: r.sq.Delete(groupsTable).
			Where(sq.Eq{groupsTableColumnToken: token}).
			Suffix("RETURNING " + strings.Join(groupsTableColumns, ",")),
	}

	var row GroupRow
	if err := r.db.DB().Getx(ctx, &row, q); err != nil {
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	return ToModel(&row), nil
}
