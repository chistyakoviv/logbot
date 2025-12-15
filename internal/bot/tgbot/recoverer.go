package tgbot

import (
	"io"
	"log/slog"
	"runtime/debug"

	"github.com/chistyakoviv/logbot/internal/lib/slogger"
	"github.com/chistyakoviv/logbot/internal/parser"
)

func TgRecoverer(
	panicWriter io.Writer,
	stackParser parser.StackParser,
	logger *slog.Logger,
) {
	if rvr := recover(); rvr != nil {
		debugStack := debug.Stack()
		out, err := stackParser.Parse(debugStack, rvr)
		if err != nil {
			logger.Error("failed to parse stack", slogger.Err(err))
			panicWriter.Write(debugStack)
		} else {
			panicWriter.Write(out)
		}
	}
}
