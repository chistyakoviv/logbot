package cancel

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/service/commands"
)

const CommandName = "cancel"

type cancelCommand struct {
	command.TgCommand
}

func New(
	ctx context.Context,
	mw middlewares.TgMiddlewareInterface,
	logger *slog.Logger,
	i18n I18n.I18nInterface,
	commands commands.ServiceInterface,
) command.TgCommandInterface {
	return &cancelCommand{
		TgCommand: command.TgCommand{
			Handler: mw.Pipe(begin(logger, i18n, commands)).Handler(ctx),
			Stages:  []handlers.Response{},
		},
	}
}

func (c *cancelCommand) ApplyDescription(lang string, i18n I18n.I18nChainInterface) {
	i18n.
		Appendf("\n\n/%s - ", CommandName).
		T(lang, "cancel_description")
}
