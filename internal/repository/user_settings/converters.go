package user_settings

import (
	"time"

	"github.com/chistyakoviv/logbot/internal/model"
)

type UserSettingsRow struct {
	UserId    int64     `db:"user_id"`
	Lang      int       `db:"lang"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r *UserSettingsRow) Values() []any {
	return []any{
		r.UserId,
		r.Lang,
		r.UpdatedAt,
	}
}

func ToModel(r *UserSettingsRow) *model.UserSettings {
	if r == nil {
		return nil
	}
	return &model.UserSettings{
		UserId:    r.UserId,
		Lang:      r.Lang,
		UpdatedAt: r.UpdatedAt,
	}
}

func FromModel(m *model.UserSettings) UserSettingsRow {
	if m == nil {
		return UserSettingsRow{}
	}
	return UserSettingsRow{
		UserId:    m.UserId,
		Lang:      m.Lang,
		UpdatedAt: m.UpdatedAt,
	}
}
