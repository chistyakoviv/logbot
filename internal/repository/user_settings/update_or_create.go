package user_settings

import (
	"context"
	"errors"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *respository) UpdateOrCreate(ctx context.Context, in *model.UserSettings) (*model.UserSettings, error) {
	row := FromModel(in)

	_, err := r.Find(ctx, in.UserId)
	if err != nil {
		// If an error is something other than "not found", return it
		if !errors.Is(err, db.ErrNotFound) {
			return nil, err
		}

		// Otherwise, the settings don't exist, so create it
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

	// Otherwise, the settings exist, so update it
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
