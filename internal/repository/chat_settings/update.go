package chat_settings

import (
	"context"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *repository) Update(ctx context.Context, in *model.ChatSettings) (*model.ChatSettings, error) {
	row := FromModel(in)

	builder := r.sq.Update(chatSettingsTable).
		Where(sq.Eq{
			chatSettingsTableColumnChatId: row.ChatId,
		}).
		Set(chatSettingsTableColumnUpdatedAt, row.UpdatedAt).
		Suffix("RETURNING " + strings.Join(chatSettingsTableColumns, ","))

	if row.CollapsePeriod != 0 {
		builder = builder.Set(chatSettingsTableColumnCollapsePeriod, row.CollapsePeriod)
	}

	if !row.MuteUntil.IsZero() {
		builder = builder.Set(chatSettingsTableColumnMuteUntil, row.MuteUntil)
	}

	q := db.Query{
		Name:    "repository.chat_settings.update",
		Sqlizer: builder,
	}

	var out ChatSettingsRow
	if err := r.db.DB().Getx(ctx, &out, q); err != nil {
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	return ToModel(&out), nil
}
