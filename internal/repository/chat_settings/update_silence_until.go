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

func (r *repository) UpdateSilenceUntil(ctx context.Context, chatId int64, silenceUntil time.Time) (*model.ChatSettings, error) {
	q := db.Query{
		Name: "repository.chat_settings.update_silence_until",
		Sqlizer: r.sq.Update(chatSettingsTable).
			Set(chatSettingsTableColumnSilenceUntil, silenceUntil).
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
