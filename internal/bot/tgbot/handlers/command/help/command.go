package help

import (
	"context"
	"log/slog"

	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
)

const CommandName = "help"

type helpCommand struct {
	command.TgCommand
}

func New(
	ctx context.Context,
	mw middlewares.TgMiddlewareInterface,
	logger *slog.Logger,
	i18n I18n.I18nInterface,
	tgCommands command.TgCommands,
) command.TgCommandInterface {
	return &helpCommand{
		TgCommand: command.TgCommand{
			Handler: mw.Pipe(begin(logger, i18n, tgCommands)).Handler(ctx),
		},
	}
}

func (c *helpCommand) ApplyDescription(lang string, i18n I18n.I18nChainInterface) {
	i18n.
		Appendf("\n\n/%s - ", CommandName).
		T(lang, "help_description")
}
