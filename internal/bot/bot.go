package bot

import (
	"context"
	"log/slog"
)

type Bot interface {
	Start(ctx context.Context, logger *slog.Logger) error
	Shutdown(ctx context.Context) error
}
