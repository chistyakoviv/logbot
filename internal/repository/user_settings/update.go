package user_settings

import (
	"context"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *respository) Update(ctx context.Context, in *model.UserSettings) (*model.UserSettings, error) {
	row := FromModel(in)

	builder := r.sq.Update(userSettingsTable).
		Where(sq.Eq{
			userSettingsTableColumnUserId: row.UserId,
		}).
		Set(userSettingsTableColumnUpdatedAt, row.UpdatedAt).
		Suffix("RETURNING " + strings.Join(userSettingsTableColumns, ","))

	if in.Lang != 0 {
		builder = builder.Set(userSettingsTableColumnLang, row.Lang)
	}

	q := db.Query{
		Name:    "repository.user_settings.update",
		Sqlizer: builder,
	}

	var out UserSettingsRow
	if err := r.db.DB().Getx(ctx, &out, q); err != nil {
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	return ToModel(&out), nil
}
