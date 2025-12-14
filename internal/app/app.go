package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	// _ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/chistyakoviv/logbot/internal/di"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
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
	bot := resolveTgBot(a.container)
	logs := resolveLogsService(a.container)

	logger.Debug("Application is running in DEBUG mode")

	initRoutes(ctx, a.container)

	// Exec the command to empty the memory buffer: echo 3 | sudo tee /proc/sys/vm/drop_caches
	// see https://medium.com/@bobzsj87/demist-the-memory-ghost-d6b7cf45dd2a
	// pprof
	// if cfg.Env == config.EnvLocal {
	// 	go func() {
	// 		logger.Info("pprof server started")
	// 		http.ListenAndServe("0.0.0.0:6060", nil)
	// 	}()
	// }

	// http server
	go func() {
		logger.Info(
			"starting http server",
			slog.String("address", cfg.HTTPServer.Host+":"+cfg.HTTPServer.Port),
			slog.String("env", cfg.Env),
		)

		srv := resolveHttpServer(a.container)
		dq.Add(func() error {
			return srv.Shutdown(ctx)
		})

		// ListenAndServe always returns a non-nil error. After [Server.Shutdown] or [Server.Close], the returned error is [ErrServerClosed]
		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server error", slogger.Err(err))
		}
		logger.Info("http server stopped")
	}()

	// tg bot
	go func() {
		logger.Info(
			"starting telegram bot",
			slog.String("address", cfg.HTTPServer.Host+":"+cfg.HTTPServer.Port),
			slog.String("domain", cfg.Webhook.Domain),
			slog.String("env", cfg.Env),
		)

		dq.Add(func() error {
			return bot.Shutdown(ctx)
		})

		// Wait for the server to be ready
		url := "http://127.0.0.1" + ":" + cfg.HTTPServer.Port + "/healthcheck"
		interval := 1 * time.Second // how often to poll
		timeout := 1 * time.Second  // per-request timeout
		startTime := time.Now()

		logger.Info(
			"waiting for the http server to be ready",
			slog.String("url", url),
			slog.String("interval", interval.String()),
			slog.String("timeout", timeout.String()),
		)

		client := &http.Client{Timeout: timeout}

		for {
			// Check if the application is terminated
			select {
			case <-ctx.Done():
				return
			default:
			}
			resp, err := client.Get(url)
			if err != nil {
				// logger.Info("healthcheck failed", slogger.Err(err))
			} else {
				err = resp.Body.Close()
				if err != nil {
					logger.Info("failed to close response body", slogger.Err(err))
				}
				if resp.StatusCode == http.StatusOK {
					break
				}
				logger.Info("healthcheck failed", slog.Int("status", resp.StatusCode))
			}

			time.Sleep(interval)
		}

		logger.Info(
			"http server is ready",
			slog.String("delay", time.Since(startTime).String()),
		)

		bot.Start(ctx)
	}()

	// Log cleaner
	go func() {
		logger.Info(
			"starting log cleaner",
			slog.String("interval", cfg.LogCleaner.Interval.String()),
			slog.String("retain", cfg.LogCleaner.Retain.String()),
		)

		ticker := time.NewTicker(cfg.LogCleaner.Interval)
		dq.Add(func() error {
			// Free resources to avoid leaks
			ticker.Stop()
			return nil
		})

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				now := time.Now().UTC()
				logs.DeleteOlderThan(ctx, now.Add(-cfg.LogCleaner.Retain))
			}
		}
	}()

	// Graceful Shutdown
	select {
	case <-ctx.Done():
		logger.Info("terminating: context canceled")
	case <-waitSignal():
		logger.Info("terminating: via signal")
	}

	cancel()

	// Call all deferred functions and wait them to be done
	dq.Release()
	dq.Wait()
}

func waitSignal() chan os.Signal {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	return stop
}
