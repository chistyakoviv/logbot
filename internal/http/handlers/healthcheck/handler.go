package healthcheck

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/chistyakoviv/logbot/internal/lib/slogger"
)

func New(
	ctx context.Context,
	logger *slog.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("alive")); err != nil {
			// optional: log or handle the error
			logger.Error("failed to write response: %v", slogger.Err(err))
		}
	}
}
