package writer

import (
	"io"
	"log/slog"
)

type panicLoggerWriter struct {
	logger *slog.Logger
}

func NewPanicLoggerWriter(logger *slog.Logger) io.Writer {
	return &panicLoggerWriter{
		logger: logger,
	}
}

func (l *panicLoggerWriter) Write(p []byte) (n int, err error) {
	l.logger.Error("panic", slog.String("stack", string(p)))
	return len(p), nil
}
