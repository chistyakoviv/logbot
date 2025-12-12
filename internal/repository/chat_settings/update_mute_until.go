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

func (r *repository) UpdateMuteUntil(ctx context.Context, chatId int64, muteUntil time.Time) (*model.ChatSettings, error) {
	q := db.Query{
		Name: "repository.chat_settings.update_mute_until",
		Sqlizer: r.sq.Update(chatSettingsTable).
			Set(chatSettingsTableColumnMuteUntil, muteUntil).
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
