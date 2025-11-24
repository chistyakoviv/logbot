package app

import (
	"log"
	"log/slog"
	"net/http"

	sq "github.com/Masterminds/squirrel"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/bot"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	tgMiddleware "github.com/chistyakoviv/logbot/internal/bot/tgbot/middleware"
	"github.com/chistyakoviv/logbot/internal/config"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/deferredq"
	"github.com/chistyakoviv/logbot/internal/di"
	"github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/rbac"
	"github.com/chistyakoviv/logbot/internal/repository/commands"
	"github.com/chistyakoviv/logbot/internal/repository/subscriptions"
	"github.com/chistyakoviv/logbot/internal/repository/user_settings"
	srvCommands "github.com/chistyakoviv/logbot/internal/service/commands"
	srvSubscriptions "github.com/chistyakoviv/logbot/internal/service/subscriptions"
	srvUserSettings "github.com/chistyakoviv/logbot/internal/service/user_settings"
	"github.com/go-chi/chi/v5"
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

func resolveRouter(c di.Container) *chi.Mux {
	router, err := di.Resolve[*chi.Mux](c, "router")

	if err != nil {
		log.Fatalf("Couldn't resolve router definition: %v", err)
	}

	return router
}

func resolveHttpServer(c di.Container) *http.Server {
	srv, err := di.Resolve[*http.Server](c, "httpServer")

	if err != nil {
		log.Fatalf("Couldn't resolve http server definition: %v", err)
	}

	return srv
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

func resolveTgCommands(c di.Container) command.TgCommands {
	commands, err := di.Resolve[command.TgCommands](c, "tgCommands")

	if err != nil {
		log.Fatalf("Couldn't resolve commands definition: %v", err)
	}

	return commands
}

func resolveTgCommandStage(c di.Container) handlers.Response {
	cmdstage, err := di.Resolve[handlers.Response](c, "tgCmdstage")

	if err != nil {
		log.Fatalf("Couldn't resolve command stage definition: %v", err)
	}

	return cmdstage
}

func resolveTgJoin(c di.Container) handlers.Response {
	join, err := di.Resolve[handlers.Response](c, "tgJoin")

	if err != nil {
		log.Fatalf("Couldn't resolve join definition: %v", err)
	}

	return join
}

func resolveI18n(c di.Container) *i18n.I18n {
	i18n, err := di.Resolve[*i18n.I18n](c, "i18n")

	if err != nil {
		log.Fatalf("Couldn't resolve i18n definition: %v", err)
	}

	return i18n
}

func resolveRbac(c di.Container) rbac.ManagerInterface {
	rbac, err := di.Resolve[rbac.ManagerInterface](c, "rbac")

	if err != nil {
		log.Fatalf("Couldn't resolve rbac definition: %v", err)
	}

	return rbac
}

func resolveTgMiddleware(c di.Container) tgMiddleware.TgMiddlewareInterface {
	middleware, err := di.Resolve[tgMiddleware.TgMiddlewareInterface](c, "tgMiddleware")

	if err != nil {
		log.Fatalf("Couldn't resolve middleware definition: %v", err)
	}

	return middleware
}

// MIddlewares
func resolveTgLangMiddleware(c di.Container) tgMiddleware.TgMiddlewareHandler {
	middleware, err := di.Resolve[tgMiddleware.TgMiddlewareHandler](c, "tgLangMiddleware")

	if err != nil {
		log.Fatalf("Couldn't resolve lang middleware definition: %v", err)
	}

	return middleware
}

// Repositories
func resolveSubscriptionsRepository(c di.Container) subscriptions.IRepository {
	repo, err := di.Resolve[subscriptions.IRepository](c, "subscriptionsRepository")

	if err != nil {
		log.Fatalf("Couldn't resolve subscriptions repository definition: %v", err)
	}

	return repo
}

func resolveCommandsRepository(c di.Container) commands.IRepository {
	repo, err := di.Resolve[commands.IRepository](c, "commandsRepository")

	if err != nil {
		log.Fatalf("Couldn't resolve commands repository definition: %v", err)
	}

	return repo
}

func resolveUserSettingsRepository(c di.Container) user_settings.IRepository {
	repo, err := di.Resolve[user_settings.IRepository](c, "userSettingsRepository")

	if err != nil {
		log.Fatalf("Couldn't resolve user settings repository definition: %v", err)
	}

	return repo
}

// Services
func resolveSubscriptionsService(c di.Container) srvSubscriptions.IService {
	service, err := di.Resolve[srvSubscriptions.IService](c, "subscriptionsService")

	if err != nil {
		log.Fatalf("Couldn't resolve subscriptions service definition: %v", err)
	}

	return service
}

func resolveCommandsService(c di.Container) srvCommands.IService {
	service, err := di.Resolve[srvCommands.IService](c, "commandsService")

	if err != nil {
		log.Fatalf("Couldn't resolve commands service definition: %v", err)
	}

	return service
}

func resolveUserSettingsService(c di.Container) srvUserSettings.IService {
	service, err := di.Resolve[srvUserSettings.IService](c, "userSettingsService")

	if err != nil {
		log.Fatalf("Couldn't resolve user settings service definition: %v", err)
	}

	return service
}
