package user_settings

import (
	"context"
	"fmt"
	"strings"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *respository) Create(ctx context.Context, in *model.UserSettings) (*model.UserSettings, error) {
	row := FromModel(in)

	q := db.Query{
		Name: "repository.user_settings.create",
		Sqlizer: r.sq.Insert(userSettingsTable).
			Columns(userSettingsTableColumns...).
			Values(row.Values()...).
			Suffix("RETURNING " + strings.Join(userSettingsTableColumns, ",")),
	}

	var out UserSettingsRow
	if err := r.db.DB().Getx(ctx, &out, q); err != nil {
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	return ToModel(&out), nil
}
