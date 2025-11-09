package app

import (
	"log"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/bot"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/cmdstage"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/config"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/deferredq"
	"github.com/chistyakoviv/logbot/internal/di"
	"github.com/chistyakoviv/logbot/internal/i18n"
)

// Retrieves the application configuration from the dependency injection container,
// centralizing error handling to avoid repetitive error checks across the codebase.
// Logs a fatal error and terminates the program if the configuration cannot be resolved.
func resolveConfig(c di.Container) *config.Config {
	cfg, err := di.Resolve[*config.Config](c, "config")

	if err != nil {
		log.Fatalf("Couldn't resolve config definition: %v", err)
	}

	return cfg
}

func resolveLogger(c di.Container) *slog.Logger {
	logger, err := di.Resolve[*slog.Logger](c, "logger")

	if err != nil {
		log.Fatalf("Couldn't resolve logger definition: %v", err)
	}

	return logger
}

func resolveDbClient(c di.Container) db.Client {
	client, err := di.Resolve[db.Client](c, "db")

	if err != nil {
		log.Fatalf("Couldn't resolve db client definition: %v", err)
	}

	return client
}

func resolveStatementBuilder(c di.Container) sq.StatementBuilderType {
	sq, err := di.Resolve[sq.StatementBuilderType](c, "sq")

	if err != nil {
		log.Fatalf("Couldn't resolve statement builder definition: %v", err)
	}

	return sq
}

func resolveDeferredQ(c di.Container) deferredq.DQueue {
	dq, err := di.Resolve[deferredq.DQueue](c, "dq")

	if err != nil {
		log.Fatalf("Couldn't resolve deferred queue definition: %v", err)
	}

	return dq
}

func resolveTxManager(c di.Container) db.TxManager {
	txManager, err := di.Resolve[db.TxManager](c, "txManager")

	if err != nil {
		log.Fatalf("Couldn't resolve tx manager definition: %v", err)
	}

	return txManager
}

func resolveTgBot(c di.Container) bot.Bot {
	bot, err := di.Resolve[bot.Bot](c, "tgBot")

	if err != nil {
		log.Fatalf("Couldn't resolve telegram bot definition: %v", err)
	}

	return bot
}

func resolveCommands(c di.Container) command.TgCommands {
	commands, err := di.Resolve[command.TgCommands](c, "tgcommands")

	if err != nil {
		log.Fatalf("Couldn't resolve commands definition: %v", err)
	}

	return commands
}

func resolveCmdstage(c di.Container) *cmdstage.TgCmdstage {
	cmdstage, err := di.Resolve[*cmdstage.TgCmdstage](c, "tgcmdstage")

	if err != nil {
		log.Fatalf("Couldn't resolve command stage definition: %v", err)
	}

	return cmdstage
}

func resolveI18n(c di.Container) *i18n.I18n {
	i18n, err := di.Resolve[*i18n.I18n](c, "i18n")

	if err != nil {
		log.Fatalf("Couldn't resolve i18n definition: %v", err)
	}

	return i18n
}
