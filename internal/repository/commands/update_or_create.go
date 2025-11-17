package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

// TODO: execute in transaction to avoid concurrent creation
func (r *repository) UpdateOrCreate(ctx context.Context, in *model.Command) (*model.Command, error) {
	row := FromModel(in)

	// Try to find the command
	_, err := r.FindByKey(ctx, &model.CommandKey{
		ChatId: in.ChatId,
		UserId: in.UserId,
	})

	if err != nil {
		// If an error is something other than "not found", return it
		if !errors.Is(err, db.ErrNotFound) {
			return nil, err
		}

		// Otherwise, the command doesn't exist, so create it
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

	// The command already exists, so update it
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
