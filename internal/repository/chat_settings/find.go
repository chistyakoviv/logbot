package chat_settings

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *repository) Find(ctx context.Context, chatId int64) (*model.ChatSettings, error) {
	q := db.Query{
		Name: "repository.chat_settings.find",
		Sqlizer: r.sq.Select(chatSettingsTableColumns...).
			From(chatSettingsTable).
			Where(sq.Eq{
				chatSettingsTableColumnChatId: chatId,
			}),
	}

	var row ChatSettingsRow
	if err := r.db.DB().Getx(ctx, &row, q); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	return ToModel(&row), nil
}
