package chat_settings

import (
	"context"
	"fmt"
	"strings"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *repository) Create(ctx context.Context, in *model.ChatSettings) (*model.ChatSettings, error) {
	row := FromModel(in)

	q := db.Query{
		Name: "repository.chat_settings.create",
		Sqlizer: r.sq.Insert(chatSettingsTable).
			Columns(chatSettingsTableColumns...).
			Values(row.Values()...).
			Suffix("RETURNING " + strings.Join(chatSettingsTableColumns, ",")),
	}

	var out ChatSettingsRow
	if err := r.db.DB().Getx(ctx, &out, q); err != nil {
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	return ToModel(&out), nil
}
