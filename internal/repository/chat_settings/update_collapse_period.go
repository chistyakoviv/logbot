package chat_settings

import (
	"context"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *repository) UpdateCollapsePeriod(ctx context.Context, chatId int64, period time.Duration) (*model.ChatSettings, error) {
	q := db.Query{
		Name: "repository.chat_settings.update_collapse_period",
		Sqlizer: r.sq.Update(chatSettingsTable).
			Set(chatSettingsTableColumnCollapsePeriod, period).
			Set(chatSettingsTableColumnUpdatedAt, time.Now()).
			Where(sq.Eq{
				chatSettingsTableColumnChatId: chatId,
			}).
			Suffix("RETURNING " + strings.Join(chatSettingsTableColumns, ",")),
	}

	var out ChatSettingsRow
	if err := r.db.DB().Getx(ctx, &out, q); err != nil {
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	return ToModel(&out), nil
}
