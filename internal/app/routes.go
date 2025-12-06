package app

import (
	"context"
	"net/http"

	"github.com/chistyakoviv/logbot/internal/di"
	"github.com/chistyakoviv/logbot/internal/http/handlers/log"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
)

func initRoutes(ctx context.Context, c di.Container) {
	router := resolveRouter(c)
	tgBot := resolveTgBot(c)
	logger := resolveLogger(c)
	validator := resolveValidator(c)
	logs := resolveLogsService(c)
	subscriptions := resolveSubscriptionsService(c)
	chatSettings := resolveChatSettingsService(c)
	labels := resolveLabelsService(c)
	loghasher := resolveLogHasher(c)

	router.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("alive")); err != nil {
			// optional: log or handle the error
			logger.Error("failed to write response: %v", slogger.Err(err))
		}
	})

	router.Post(tgBot.WebhookPath()+"{token}", tgBot.HandlerFunc())

	router.Post("/log", log.New(
		ctx,
		logger,
		validator,
		tgBot,
		loghasher,
		logs,
		subscriptions,
		chatSettings,
		labels,
	))
}
