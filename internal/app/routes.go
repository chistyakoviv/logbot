package app

import (
	"context"
	"net/http"

	"github.com/chistyakoviv/logbot/internal/di"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
)

func initRoutes(ctx context.Context, c di.Container) {
	router := resolveRouter(c)
	tgBot := resolveTgBot(c)

	router.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		logger := resolveLogger(c)
		if _, err := w.Write([]byte("alive")); err != nil {
			// optional: log or handle the error
			logger.Error("failed to write response: %v", slogger.Err(err))
		}
	})

	// router.Post("/convert", convert.New(
	// 	ctx,
	// 	resolveLogger(c),
	// 	resolveValidator(c),
	// 	resolveConversionQueueService(c),
	// 	resolveTaskService(c),
	// ))

	router.Post(tgBot.WebhookPath(), tgBot.HandlerFunc())
}
