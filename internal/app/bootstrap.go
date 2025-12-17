package app

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/bot"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/handler"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/messages"
	tgMiddlewares "github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	tgMiddleware "github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares/middleware"
	"github.com/chistyakoviv/logbot/internal/config"
	"github.com/chistyakoviv/logbot/internal/constants"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/db/pg"
	"github.com/chistyakoviv/logbot/internal/db/transaction"
	"github.com/chistyakoviv/logbot/internal/di"
	httpMiddleware "github.com/chistyakoviv/logbot/internal/http/middleware"
	"github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/lib/chrono"
	"github.com/chistyakoviv/logbot/internal/lib/deferredq"
	"github.com/chistyakoviv/logbot/internal/lib/loghasher"
	"github.com/chistyakoviv/logbot/internal/lib/markdown"
	"github.com/chistyakoviv/logbot/internal/lib/panic_writer"
	"github.com/chistyakoviv/logbot/internal/lib/rbac"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
	"github.com/chistyakoviv/logbot/internal/lib/stack_parser"
	"github.com/chistyakoviv/logbot/internal/repository/chat_settings"
	"github.com/chistyakoviv/logbot/internal/repository/commands"
	"github.com/chistyakoviv/logbot/internal/repository/labels"
	"github.com/chistyakoviv/logbot/internal/repository/last_sent"
	"github.com/chistyakoviv/logbot/internal/repository/logs"
	"github.com/chistyakoviv/logbot/internal/repository/subscriptions"
	"github.com/chistyakoviv/logbot/internal/repository/user_settings"
	srvChatSettings "github.com/chistyakoviv/logbot/internal/service/chat_settings"
	srvCommands "github.com/chistyakoviv/logbot/internal/service/commands"
	srvLabels "github.com/chistyakoviv/logbot/internal/service/labels"
	srvLogs "github.com/chistyakoviv/logbot/internal/service/logs"
	srvSubscriptions "github.com/chistyakoviv/logbot/internal/service/subscriptions"
	srvUserSettings "github.com/chistyakoviv/logbot/internal/service/user_settings"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
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

	c.RegisterSingleton("stackParser", func(c di.Container) stack_parser.StackParser {
		cfg := resolveConfig(c)
		var stackParser stack_parser.StackParser

		switch cfg.Env {
		case config.EnvProd:
			// Do not use colors in production
			stackParser = stack_parser.NewSimpleStackParser()
		case config.EnvDev:
			stackParser = stack_parser.NewSimpleStackParser()
		default:
			// Use colors in development
			stackParser = stack_parser.NewPrettyStackParser()
		}

		return stackParser
	})

	c.RegisterSingleton("panicWriter", func(c di.Container) io.Writer {
		cfg := resolveConfig(c)
		logger := resolveLogger(c)

		var panicWriter io.Writer

		switch cfg.Env {
		case config.EnvProd:
			// Write panics to logger in production
			panicWriter = panic_writer.NewLoggerPanicWriter(logger)
		case config.EnvDev:
			panicWriter = panic_writer.NewLoggerPanicWriter(logger)
		default:
			// Write panics to stderr in development
			panicWriter = os.Stderr
		}

		return panicWriter
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

	c.RegisterSingleton("router", func(c di.Container) *chi.Mux {
		router := chi.NewRouter()
		logger := resolveLogger(c)
		panicWriter := resolvePanicWriter(c)
		stackParser := resolveStackParser(c)

		router.Use(middleware.RequestID)
		router.Use(httpMiddleware.NewLogger(logger))
		router.Use(httpMiddleware.NewRecoverer(panicWriter, stackParser, logger))
		router.Use(middleware.NoCache)

		return router
	})

	c.RegisterSingleton("httpServer", func(c di.Container) *http.Server {
		cfg := resolveConfig(c)
		router := resolveRouter(c)

		return &http.Server{
			Addr:         cfg.HTTPServer.Host + ":" + cfg.HTTPServer.Port,
			Handler:      router,
			ReadTimeout:  cfg.HTTPServer.ReadTimeout,
			WriteTimeout: cfg.HTTPServer.WriteTimeout,
			IdleTimeout:  cfg.HTTPServer.IdleTimeout,
		}
	})

	c.RegisterSingleton("dq", func(c di.Container) deferredq.DQueue {
		return deferredq.New(resolveLogger(c))
	})

	c.RegisterSingleton("validator", func(c di.Container) *validator.Validate {
		return validator.New()
	})

	c.RegisterSingleton("txManager", func(c di.Container) db.TxManager {
		return transaction.NewTransactionManager(resolveDbClient(c).DB())
	})

	c.RegisterSingleton("tgCmdstage", func(c di.Container) handlers.Response {
		logger := resolveLogger(c)
		i18n := resolveI18n(c)
		tgCommands := resolveTgCommands(c)
		srvCommands := resolveCommandsService(c)
		return command.NewCommandStage(
			ctx,
			logger,
			i18n,
			srvCommands,
			tgCommands,
		)
	})

	c.RegisterSingleton("tgJoin", func(c di.Container) handlers.Response {
		logger := resolveLogger(c)
		i18n := resolveI18n(c)
		mw := resolveTgMiddleware(c)

		// Middlewares
		mwRecoverer := resolveTgRecovererMiddleware(c)
		mwLang := resolveTgLangMiddleware(c)

		mw = mw.Pipe(mwRecoverer).Pipe(mwLang)
		return handler.NewJoin(ctx, mw, logger, i18n)
	})

	c.RegisterSingleton("tgCommands", func(c di.Container) command.TgCommands {
		return BuildTgCommands(ctx, c)
	})

	c.RegisterSingleton("tgBot", func(c di.Container) bot.Bot {
		return tgbot.New(&tgbot.TgBotSpec{
			Cfg:      resolveConfig(c),
			Commands: resolveTgCommands(c),
			Cmdstage: resolveTgCommandStage(c),
			Join:     resolveTgJoin(c),
			Logger:   resolveLogger(c),
		})
	})

	c.RegisterSingleton("i18n", func(c di.Container) i18n.I18nInterface {
		return i18n.New()
	})

	c.RegisterSingleton("rbac", func(c di.Container) rbac.ManagerInterface {
		cfg := resolveConfig(c)

		ruleFactory := rbac.NewRuleFactory()
		// ruleFactory.Add("superuser", func() rbac.RuleInterface {
		// 	return NewSuperuserRule()
		// })

		itemsStorage := rbac.NewItemsStorageInMemory()
		assignmentsStorage := rbac.NewAssignmentsStorageInMemory()
		manager := rbac.NewManager(ruleFactory, itemsStorage, assignmentsStorage, nil)

		// Set guest role name
		manager.SetGuestRoleName("guest")

		// RBAC hierarchy
		// Permissions
		managePermission := rbac.NewPermission(constants.PermissionManage)
		// managePermission = managePermission.WithRuleName("superuser")
		err := manager.AddPermission(managePermission)
		if err != nil {
			log.Fatalf("couldn't add permission: %v", err)
		}

		// Roles
		superuser := rbac.NewRole("superuser")
		err = manager.AddRole(superuser)
		if err != nil {
			log.Fatalf("couldn't add role: %v", err)
		}

		// Attach permissions to roles
		err = manager.AddChild(superuser.GetName(), managePermission.GetName())
		if err != nil {
			log.Fatalf("couldn't add permission to role: %v", err)
		}

		// Assignments (store assignments in database)
		// Assign the superuser role to a user id (assigning permissions directly is disabled by default)
		superuserId, err := strconv.ParseInt(cfg.Superuser, 10, 64)
		if err != nil {
			log.Fatalf("couldn't parse superuser id from config")
		}
		err = manager.Assign(superuserId, superuser.GetName(), time.Now())
		if err != nil {
			log.Fatalf("couldn't assign superuser role: %v", err)
		}

		return manager
	})

	c.RegisterSingleton("logHasher", func(c di.Container) loghasher.HasherInterface {
		return loghasher.NewHasher()
	})

	c.RegisterSingleton("markdowner", func(c di.Container) markdown.MarkdownerInterface {
		return markdown.NewMarkdowner()
	})

	c.RegisterSingleton("chrono", func(c di.Container) chrono.Chrono {
		return chrono.New()
	})

	// Report messages
	c.RegisterSingleton("errorReportMessage", func(c di.Container) messages.ErrorReportMessageInterface {
		return messages.NewErrorReportMessage(resolveMarkdowner(c))
	})

	// Tg middleware
	c.RegisterSingleton("tgMiddleware", func(c di.Container) tgMiddlewares.TgMiddlewareChainInterface {
		return tgMiddlewares.NewMiddleware()
	})

	// Tg middlewares
	c.RegisterSingleton("tgLangMiddleware", func(c di.Container) tgMiddlewares.TgMiddleware {
		return tgMiddleware.Lang(resolveLogger(c), resolveUserSettingsService(c))
	})

	c.RegisterSingleton("tgSubscriptionMiddleware", func(c di.Container) tgMiddlewares.TgMiddleware {
		return tgMiddleware.Subscription(resolveLogger(c), resolveI18n(c), resolveSubscriptionsService(c))
	})

	c.RegisterSingleton("tgSuperuserMiddleware", func(c di.Container) tgMiddlewares.TgMiddleware {
		return tgMiddleware.Superuser(resolveLogger(c), resolveI18n(c), resolveRbac(c))
	})

	c.RegisterSingleton("tgSilenceMiddleware", func(c di.Container) tgMiddlewares.TgMiddleware {
		return tgMiddleware.Silence(resolveLogger(c), resolveI18n(c), resolveChatSettingsService(c))
	})

	c.RegisterSingleton("tgRecovererMiddleware", func(c di.Container) tgMiddlewares.TgMiddleware {
		return tgMiddleware.Recoverer(resolveLogger(c), resolvePanicWriter(c), resolveStackParser(c))
	})

	// Repositories
	c.RegisterSingleton("subscriptionsRepository", func(c di.Container) subscriptions.RepositoryInterface {
		return subscriptions.NewRepository(resolveDbClient(c), resolveStatementBuilder(c))
	})

	c.RegisterSingleton("commandsRepository", func(c di.Container) commands.RepositoryInterface {
		return commands.NewRepository(resolveDbClient(c), resolveStatementBuilder(c))
	})

	c.RegisterSingleton("userSettingsRepository", func(c di.Container) user_settings.RepositoryInterface {
		return user_settings.NewRepository(resolveDbClient(c), resolveStatementBuilder(c))
	})

	c.RegisterSingleton("chatSettingsRepository", func(c di.Container) chat_settings.RepositoryInterface {
		return chat_settings.NewRepository(resolveDbClient(c), resolveStatementBuilder(c))
	})

	c.RegisterSingleton("labelsRepository", func(c di.Container) labels.RepositoryInterface {
		return labels.NewRepository(resolveDbClient(c), resolveStatementBuilder(c))
	})

	c.RegisterSingleton("logsRepository", func(c di.Container) logs.RepositoryInterface {
		return logs.NewRepository(resolveDbClient(c), resolveStatementBuilder(c))
	})

	c.RegisterSingleton("lastSentRepository", func(c di.Container) last_sent.RepositoryInterface {
		return last_sent.NewRepository(resolveDbClient(c), resolveStatementBuilder(c))
	})

	// Services
	c.RegisterSingleton("subscriptionsService", func(c di.Container) srvSubscriptions.ServiceInterface {
		return srvSubscriptions.NewService(resolveSubscriptionsRepository(c), resolveTxManager(c))
	})

	c.RegisterSingleton("commandsService", func(c di.Container) srvCommands.ServiceInterface {
		return srvCommands.NewService(resolveCommandsRepository(c), resolveTxManager(c))
	})

	c.RegisterSingleton("userSettingsService", func(c di.Container) srvUserSettings.ServiceInterface {
		return srvUserSettings.NewService(resolveUserSettingsRepository(c), resolveTxManager(c))
	})

	c.RegisterSingleton("chatSettingsService", func(c di.Container) srvChatSettings.ServiceInterface {
		return srvChatSettings.NewService(resolveChatSettingsRepository(c), resolveTxManager(c))
	})

	c.RegisterSingleton("labelsService", func(c di.Container) srvLabels.ServiceInterface {
		return srvLabels.NewService(resolveLabelsRepository(c), resolveTxManager(c))
	})

	c.RegisterSingleton("logsService", func(c di.Container) srvLogs.ServiceInterface {
		return srvLogs.NewService(
			resolveLogsRepository(c),
			resolveSubscriptionsRepository(c),
			resolveChatSettingsRepository(c),
			resolveLabelsRepository(c),
			resolveLastSentRepository(c),
			resolveLogHasher(c),
			resolveChrono(c),
			resolveTgBot(c),
			resolveErrorReportMessage(c),
			resolveLogger(c),
			resolveTxManager(c),
		)
	})
}
