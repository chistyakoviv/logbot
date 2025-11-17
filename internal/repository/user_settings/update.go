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
	builder := r.sq.Update(userSettingsTable).
		Where(sq.Eq{
			userSettingsTableColumnUserId: in.UserId,
		}).
		Suffix("RETURNING " + strings.Join(userSettingsTableColumns, ","))

	if in.Lang != 0 {
		builder = builder.Set(userSettingsTableColumnLang, in.Lang)
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
