package commands

import (
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

type TgCommand struct {
	Name    string
	Handler handlers.Response
}

type TgCommandHandler func(logger *slog.Logger) handlers.Response

func NewTgCommand(name string, handler TgCommandHandler, logger *slog.Logger) *TgCommand {
	return &TgCommand{
		Name:    name,
		Handler: handler(logger),
	}
}
