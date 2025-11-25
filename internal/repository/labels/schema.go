package labels

const labelsTable = "labels"

const (
	labelsTableColumnChatId    = "chat_id"
	labelsTableColumnUserId    = "user_id"
	labelsTableColumnLabels    = "labels"
	labelsTableColumnUpdatedAt = "updated_at"
)

var labelsTableColumns = []string{
	labelsTableColumnChatId,
	labelsTableColumnUserId,
	labelsTableColumnLabels,
	labelsTableColumnUpdatedAt,
}
