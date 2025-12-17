package messages

import (
	"bytes"
	"encoding/json"
	"sort"

	"github.com/chistyakoviv/logbot/internal/lib/markdown"
	"github.com/chistyakoviv/logbot/internal/model"
)

const stackTraceKey = "stack"

type ErrorReportMessageInterface interface {
	Create(*model.Log, *model.Subscription, []*model.Label) string
}

type errorReportMessage struct {
	markdowner markdown.MarkdownerInterface
}

func NewErrorReportMessage(
	markdowner markdown.MarkdownerInterface,
) ErrorReportMessageInterface {
	return &errorReportMessage{
		markdowner: markdowner,
	}
}

func (m *errorReportMessage) Create(
	log *model.Log,
	subscription *model.Subscription,
	subscribers []*model.Label,
) string {
	var message bytes.Buffer
	for i, subscriber := range subscribers {
		message.WriteString("@")
		message.WriteString(subscriber.Username)
		if i < len(subscribers)-1 {
			message.WriteString(", ")
		}
	}

	message.WriteString("\n\n*Project ")
	message.WriteString(subscription.ProjectName)
	message.WriteString("*\n")

	message.WriteString("\n*Info*\n")

	message.WriteString("service: `")
	message.WriteString(log.Service)
	message.WriteString("`\n")

	if len(log.ContainerId) > 0 {
		message.WriteString("container id: `")
		message.WriteString(log.ContainerId)
		message.WriteString("`\n")
	}

	if len(log.NodeId) > 0 {
		message.WriteString("node id: `")
		message.WriteString(log.NodeId)
		message.WriteString("`\n")
	}
	message.WriteString("\n*Data*\n")

	var decodedData map[string]string
	err := json.Unmarshal([]byte(log.Data), &decodedData)

	if err == nil {
		keys := make([]string, 0, len(decodedData))
		for key := range decodedData {
			if key != stackTraceKey {
				keys = append(keys, key)
			}
		}
		// Always print entries in the same order
		sort.Strings(keys)
		for _, key := range keys {
			message.WriteString(m.markdowner.Escape(key))
			message.WriteString(": ")
			message.WriteString(m.markdowner.Escape(decodedData[key]))
			message.WriteString("\n")
		}

		// Print stack trace in code block
		if code, ok := decodedData[stackTraceKey]; ok {
			message.WriteString("```\n")
			message.WriteString(code)
			message.WriteString("\n```")
		}
	} else {
		// Log data is in an unexpected format, print it in code block
		message.WriteString("```\n")
		message.WriteString(log.Data)
		message.WriteString("\n```")
	}

	return message.String()
}
