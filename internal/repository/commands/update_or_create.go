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
	newRow := FromModel(in)

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
				Values(newRow.Values()...).
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
			commandsTableColumnUserId: newRow.UserId,
			commandsTableColumnChatId: newRow.ChatId,
		}).
		Set(commandsTableColumnStage, newRow.Stage).
		Set(commandsTableColumnUpdatedAt, newRow.UpdatedAt).
		Suffix("RETURNING " + strings.Join(commandsTableColumns, ","))

	if newRow.Data != nil {
		builder = builder.Set(commandsTableColumnData, newRow.Data)
	}

	if newRow.Name != "" {
		builder = builder.Set(commandsTableColumnName, newRow.Name)
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
