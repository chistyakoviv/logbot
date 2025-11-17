package user_settings

const userSettingsTable = "user_settings"

const (
	userSettingsTableColumnUserId    = "user_id"
	userSettingsTableColumnLang      = "lang"
	userSettingsTableColumnUpdatedAt = "updated_at"
)

var userSettingsTableColumns = []string{
	userSettingsTableColumnUserId,
	userSettingsTableColumnLang,
	userSettingsTableColumnUpdatedAt,
}
