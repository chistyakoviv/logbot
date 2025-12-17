package app

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/di"
	"github.com/chistyakoviv/logbot/internal/http/handlers/healthcheck"
	"github.com/chistyakoviv/logbot/internal/http/handlers/log"
)

func initRoutes(ctx context.Context, c di.Container) {
	router := resolveRouter(c)
	tgBot := resolveTgBot(c)
	logger := resolveLogger(c)
	validator := resolveValidator(c)
	logs := resolveLogsService(c)

	router.Get("/healthcheck", healthcheck.New(ctx, logger))

	router.Post(tgBot.WebhookPath()+"{token}", tgBot.HandlerFunc())

	router.Post("/log", log.New(
		ctx,
		logger,
		validator,
		logs,
	))
}
