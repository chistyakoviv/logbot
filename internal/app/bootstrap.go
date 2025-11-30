package app

import (
	"context"
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
	tgMiddleware "github.com/chistyakoviv/logbot/internal/bot/tgbot/middleware"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middleware/middlewares"
	"github.com/chistyakoviv/logbot/internal/config"
	"github.com/chistyakoviv/logbot/internal/constants"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/db/pg"
	"github.com/chistyakoviv/logbot/internal/db/transaction"
	"github.com/chistyakoviv/logbot/internal/deferredq"
	"github.com/chistyakoviv/logbot/internal/di"
	mwLogger "github.com/chistyakoviv/logbot/internal/http/middleware/logger"
	"github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
	"github.com/chistyakoviv/logbot/internal/rbac"
	"github.com/chistyakoviv/logbot/internal/repository/commands"
	"github.com/chistyakoviv/logbot/internal/repository/labels"
	"github.com/chistyakoviv/logbot/internal/repository/subscriptions"
	"github.com/chistyakoviv/logbot/internal/repository/user_settings"
	srvCommands "github.com/chistyakoviv/logbot/internal/service/commands"
	srvLabels "github.com/chistyakoviv/logbot/internal/service/labels"
	srvSubscriptions "github.com/chistyakoviv/logbot/internal/service/subscriptions"
	srvUserSettings "github.com/chistyakoviv/logbot/internal/service/user_settings"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
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

	c.RegisterSingleton("router", func(c di.Container) *chi.Mux {
		router := chi.NewRouter()
		logger := resolveLogger(c)

		router.Use(middleware.RequestID)
		// Replace middleware.Logger with custom logger middleware to keep logs consistent with the rest of the application
		// router.Use(middleware.Logger)
		router.Use(mwLogger.New(logger))
		// router.Use(middleware.Heartbeat("/ping"))
		router.Use(middleware.Recoverer)
		router.Use(middleware.URLFormat)
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

	c.RegisterSingleton("txManager", func(c di.Container) db.TxManager {
		return transaction.NewTransactionManager(resolveDbClient(c).DB())
	})

	c.RegisterSingleton("tgCmdstage", func(c di.Container) handlers.Response {
		logger := resolveLogger(c)
		i18n := resolveI18n(c)
		tgCommands := resolveTgCommands(c)
		return handler.NewCommandStage(ctx, logger, i18n, resolveCommandsService(c), tgCommands)
	})

	c.RegisterSingleton("tgJoin", func(c di.Container) handlers.Response {
		logger := resolveLogger(c)
		i18n := resolveI18n(c)
		return handler.NewJoin(logger, i18n)
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

	c.RegisterSingleton("tgMiddleware", func(c di.Container) tgMiddleware.TgMiddlewareInterface {
		return tgMiddleware.NewMiddleware()
	})

	// Middlewares
	c.RegisterSingleton("tgLangMiddleware", func(c di.Container) tgMiddleware.TgMiddlewareHandler {
		return middlewares.Lang(resolveLogger(c), resolveUserSettingsService(c))
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

	c.RegisterSingleton("labelsRepository", func(c di.Container) labels.RepositoryInterface {
		return labels.NewRepository(resolveDbClient(c), resolveStatementBuilder(c))
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

	c.RegisterSingleton("labelsService", func(c di.Container) srvLabels.ServiceInterface {
		return srvLabels.NewService(resolveLabelsRepository(c), resolveTxManager(c))
	})
}
