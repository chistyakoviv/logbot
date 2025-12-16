package middleware

import (
	"context"
	"io"
	"log/slog"
	"runtime/debug"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
	"github.com/chistyakoviv/logbot/internal/lib/stack_parser"
)

func Recoverer(
	logger *slog.Logger,
	panicWriter io.Writer,
	stackParser stack_parser.StackParser,
) middlewares.TgMiddleware {
	fn := func(next middlewares.TgMiddlewareHandler) middlewares.TgMiddlewareHandler {
		return func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) error {
			logger.Debug("recoverer middleware")

			defer func() {
				if rvr := recover(); rvr != nil {
					var writeErr error
					debugStack := debug.Stack()
					out, err := stackParser.Parse(debugStack, rvr)
					if err != nil {
						logger.Error("failed to parse stack", slogger.Err(err))
						// Just write raw stack in case of parse error
						_, writeErr = panicWriter.Write(debugStack)
					} else {
						_, writeErr = panicWriter.Write(out)
					}
					if writeErr != nil {
						logger.Error("failed to write stack", slogger.Err(err))
					}
				}
			}()
			return next(ctx, b, ectx)
		}
	}
	return fn
}
