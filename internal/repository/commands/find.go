package commands

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *repository) FindByKey(ctx context.Context, in *model.CommandKey) (*model.Command, error) {
	q := db.Query{
		Name: "repository.commands.find",
		Sqlizer: r.sq.Select(commandsTableColumns...).
			From(commandsTable).
			Where(sq.Eq{
				commandsTableColumnUserId: in.UserId,
				commandsTableColumnChatId: in.ChatId,
			}),
	}

	var row CommandRow
	if err := r.db.DB().Getx(ctx, &row, q); err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	return ToModel(&row), nil
}
