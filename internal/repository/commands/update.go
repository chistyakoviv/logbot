package commands

import (
	"context"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *repository) Update(ctx context.Context, in *model.Command) (*model.Command, error) {
	row := FromModel(in)

	builder := r.sq.Update(commandsTable).
		Where(sq.Eq{
			commandsTableColumnUserId: row.UserId,
			commandsTableColumnChatId: row.ChatId,
		}).
		Set(commandsTableColumnStage, row.Stage).
		Set(commandsTableColumnUpdatedAt, row.UpdatedAt).
		Suffix("RETURNING " + strings.Join(commandsTableColumns, ","))

	if row.Data != nil {
		builder = builder.Set(commandsTableColumnData, row.Data)
	}

	if row.Name != "" {
		builder = builder.Set(commandsTableColumnName, row.Name)
	}

	q := db.Query{
		Name:    "repository.commands.update",
		Sqlizer: builder,
	}

	var out CommandRow
	if err := r.db.DB().Getx(ctx, &out, q); err != nil {
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	return ToModel(&out), nil
}
