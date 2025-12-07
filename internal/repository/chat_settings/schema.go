package chat_settings

const chatSettingsTable = "chat_settings"

const (
	chatSettingsTableColumnChatId         = "chat_id"
	chatSettingsTableColumnCollapsePeriod = "collapse_period"
	chatSettingsTableColumnMuteUntil      = "mute_until"
	chatSettingsTableColumnUpdatedAt      = "updated_at"
)

var chatSettingsTableColumns = []string{
	chatSettingsTableColumnChatId,
	chatSettingsTableColumnCollapsePeriod,
	chatSettingsTableColumnMuteUntil,
	chatSettingsTableColumnUpdatedAt,
}
