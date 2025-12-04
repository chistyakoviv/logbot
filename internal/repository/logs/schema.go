package logs

const logsTable = "logs"

const (
	logsTableColumnId         = "id"
	logsTableColumnToken      = "token"
	logsTableColumnData       = "data"
	logsTableColumnLabel      = "label"
	logsTableColumnHash       = "hash"
	logsTableColumnCreateddAt = "created_at"
)

var logsTableColumns = []string{
	logsTableColumnId,
	logsTableColumnToken,
	logsTableColumnData,
	logsTableColumnLabel,
	logsTableColumnHash,
	logsTableColumnCreateddAt,
}
