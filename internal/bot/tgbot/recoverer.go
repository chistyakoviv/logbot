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
}
