package commands

const commandsTable = "commands"

const (
	commandsTableColumnName      = "name"
	commandsTableColumnUserId    = "user_id"
	commandsTableColumnChatId    = "chat_id"
	commandsTableColumnStage     = "stage"
	commandsTableColumnData      = "data"
	commandsTableColumnUpdatedAt = "updated_at"
)

var commandsTableColumns = []string{
	commandsTableColumnName,
	commandsTableColumnUserId,
	commandsTableColumnChatId,
	commandsTableColumnStage,
	commandsTableColumnData,
	commandsTableColumnUpdatedAt,
}
