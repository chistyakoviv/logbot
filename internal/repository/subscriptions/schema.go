package subscriptions

const subscriptionsTable = "subscriptions"

const (
	subscriptionsTableColumnId        = "id"
	subscriptionsTableColumnChatId    = "chat_id"
	subscriptionsTableColumnToken     = "token"
	subscriptionsTableColumnCreatedAt = "created_at"
)

var subscriptionsTableColumns = []string{
	subscriptionsTableColumnId,
	subscriptionsTableColumnChatId,
	subscriptionsTableColumnToken,
	subscriptionsTableColumnCreatedAt,
}
