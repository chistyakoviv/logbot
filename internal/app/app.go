package app

import (
	"context"
	"log/slog"

	// _ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/chistyakoviv/logbot/internal/di"
)

type Application interface {
	Run(ctx context.Context)
	Container() di.Container
}

type app struct {
	container di.Container
}

func NewApp(ctx context.Context) Application {
	container := di.NewContainer()
	a := &app{
		container: container,
	}
	a.init(ctx)
	return a
}

func (a *app) Container() di.Container {
	return a.container
}

func (a *app) init(ctx context.Context) {
	bootstrap(ctx, a.container)
}

func (a *app) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)

	cfg := resolveConfig(a.container)
	logger := resolveLogger(a.container)
	dq := resolveDeferredQ(a.container)

	logger.Debug("Application is running in DEBUG mode")

	// Exec the command to empty the memory buffer: echo 3 | sudo tee /proc/sys/vm/drop_caches
	// see https://medium.com/@bobzsj87/demist-the-memory-ghost-d6b7cf45dd2a
	// pprof
	// if cfg.Env == config.EnvLocal {
	// 	go func() {
	// 		logger.Info("pprof server started")
	// 		http.ListenAndServe("0.0.0.0:6060", nil)
	// 	}()
	// }

	// tg bot
	go func() {
		logger.Info(
			"starting telegram bot",
			slog.String("address", cfg.ListenAddr),
			slog.String("domain", cfg.Webhook.Domain),
			slog.String("env", cfg.Env),
		)

		bot := resolveTgBot(a.container)
		dq.Add(func() error {
			return bot.Shutdown(ctx)
		})

		_ = bot.Start(ctx, logger)
		logger.Info("telegram bot stopped")
	}()

	// Graceful Shutdown
	select {
	case <-ctx.Done():
		logger.Info("terminating: context canceled")
	// No need for a wait group until the application is blocked, waiting for an OS signal
	case <-waitSignal():
		logger.Info("terminating: via signal")
	}

	// Call all deferred functions and wait them to be done
	dq.Release()
	dq.Wait()

	cancel()
}

func waitSignal() chan os.Signal {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	return stop
}
