package user_settings

const userSettingsTable = "user_settings"

const (
	userSettingsTableColumnUserId    = "user_id"
	userSettingsTableColumnUsername  = "username"
	userSettingsTableColumnLang      = "lang"
	userSettingsTableColumnUpdatedAt = "updated_at"
)

var userSettingsTableColumns = []string{
	userSettingsTableColumnUserId,
	userSettingsTableColumnUsername,
	userSettingsTableColumnLang,
	userSettingsTableColumnUpdatedAt,
}
