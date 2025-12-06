package last_sent

const lastSentTable = "last_sent"

const (
	lastSentTableColumnChatId    = "chat_id"
	lastSentTableColumnToken     = "token"
	lastSentTableColumnHash      = "hash"
	lastSentTableColumnUpdatedAt = "updated_at"
)

var lastSentTableColumns = []string{
	lastSentTableColumnChatId,
	lastSentTableColumnToken,
	lastSentTableColumnHash,
	lastSentTableColumnUpdatedAt,
}
