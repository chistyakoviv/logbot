package subscriptions

const subscriptionsTable = "subscriptions"

const (
	subscriptionsTableColumnID        = "id"
	subscriptionsTableColumnChatID    = "chat_id"
	subscriptionsTableColumnToken     = "token"
	subscriptionsTableColumnCreatedAt = "created_at"
)

var subscriptionsTableColumns = []string{
	subscriptionsTableColumnID,
	subscriptionsTableColumnChatID,
	subscriptionsTableColumnToken,
	subscriptionsTableColumnCreatedAt,
}
