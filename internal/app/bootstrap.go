package app

import (
	"context"
	"log/slog"
	"os"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/bot"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot"
	"github.com/chistyakoviv/logbot/internal/commands/tgcommand"
	"github.com/chistyakoviv/logbot/internal/config"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/db/pg"
	"github.com/chistyakoviv/logbot/internal/db/transaction"
	"github.com/chistyakoviv/logbot/internal/deferredq"
	"github.com/chistyakoviv/logbot/internal/di"
	"github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
)

func bootstrap(ctx context.Context, c di.Container) {
	c.RegisterSingleton("config", func(c di.Container) *config.Config {
		cfg := config.MustLoad(nil)
		return cfg
	})

	c.RegisterSingleton("logger", func(c di.Container) *slog.Logger {
		cfg := resolveConfig(c)

		var logger *slog.Logger

		switch cfg.Env {
		case config.EnvProd:
			logger = slog.New(
				slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
			)
		case config.EnvDev:
			logger = slog.New(
				slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			)
		default:
			logger = slog.New(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			)
		}

		logger = logger.With(
			slog.String("service", "logbot"),
		)

		return logger
	})

	c.RegisterSingleton("db", func(c di.Container) db.Client {
		cfg := resolveConfig(c)
		logger := resolveLogger(c)
		dq := resolveDeferredQ(c)

		client, err := pg.NewClient(ctx, cfg.Postgres.Dsn, logger)

		// Close db connections
		dq.Add(func() error {
			defer logger.Info("db connections closed")
			return client.Close()
		})

		if err != nil {
			logger.Error("failed to create db client", slogger.Err(err))
			os.Exit(1)
		}

		return client
	})

	c.RegisterSingleton("sq", func(c di.Container) sq.StatementBuilderType {
		return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	})

	c.RegisterSingleton("dq", func(c di.Container) deferredq.DQueue {
		return deferredq.New(resolveLogger(c))
	})

	c.RegisterSingleton("txManager", func(c di.Container) db.TxManager {
		return transaction.NewTransactionManager(resolveDbClient(c).DB())
	})

	c.RegisterSingleton("commands", func(c di.Container) tgcommand.TgCommands {
		return BuildTgCommands(ctx, c)
	})

	c.RegisterSingleton("tgBot", func(c di.Container) bot.Bot {
		return tgbot.New(resolveConfig(c), resolveCommands(c))
	})

	c.RegisterSingleton("i18n", func(c di.Container) *i18n.I18n {
		return i18n.New()
	})
}
