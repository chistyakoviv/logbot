package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *repository) Create(ctx context.Context, in *model.Command) (*model.Command, error) {
	row := FromModel(in)

	q := db.Query{
		Name: "repository.commands.create",
		Sqlizer: r.sq.Insert(commandsTable).
			Columns(commandsTableColumns...).
			Values(row.Values()...).
			Suffix("RETURNING " + strings.Join(commandsTableColumns, ",")),
	}

	var out CommandRow
	if err := r.db.DB().Getx(ctx, &out, q); err != nil {
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	return ToModel(&out), nil
}
