package bot

import (
	"context"
	"log/slog"
	"net/http"
)

type Bot interface {
	Start(ctx context.Context, logger *slog.Logger) error
	Shutdown(ctx context.Context) error
	HandlerFunc() http.HandlerFunc
	WebhookPath() string
}
