package groups

const groupsTable = "groups"

const (
	groupsTableColumnID        = "id"
	groupsTableColumnChatID    = "chat_id"
	groupsTableColumnToken     = "token"
	groupsTableColumnCreatedAt = "created_at"
)

var groupsTableColumns = []string{
	groupsTableColumnID,
	groupsTableColumnChatID,
	groupsTableColumnToken,
	groupsTableColumnCreatedAt,
}
