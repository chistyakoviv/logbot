package user_settings

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *respository) Find(ctx context.Context, id int64) (*model.UserSettings, error) {
	q := db.Query{
		Name: "repository.user_settings.find",
		Sqlizer: r.sq.Select(userSettingsTableColumns...).
			From(userSettingsTable).
			Where(sq.Eq{
				userSettingsTableColumnUserId: id,
			}),
	}

	var row UserSettingsRow
	if err := r.db.DB().Getx(ctx, &row, q); err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return nil, err
		}
	}

	return ToModel(&row), nil
}
