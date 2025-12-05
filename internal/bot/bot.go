package bot

import (
	"context"
	"net/http"
)

type Bot interface {
	Start(ctx context.Context)
	Shutdown(ctx context.Context) error
	HandlerFunc() http.HandlerFunc
	WebhookPath() string
	SendMessage(chatId int64, text string, opts interface{}) error
}
