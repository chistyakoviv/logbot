package gentoken

import (
	"context"
	"log/slog"

	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
)

const CommandName = "gentoken"

type gentokenCommand struct {
	command.TgCommand
}

func New(
	ctx context.Context,
	mw middlewares.TgMiddlewareInterface,
	logger *slog.Logger,
	i18n I18n.I18nInterface,
) command.TgCommandInterface {
	return &gentokenCommand{
		TgCommand: command.TgCommand{
			Handler: mw.Pipe(begin(logger, i18n)).Handler(ctx),
		},
	}
}

func (c *gentokenCommand) ApplyDescription(lang string, i18n I18n.I18nChainInterface) {
	i18n.
		Appendf("\n\n/%s - ", CommandName).
		T(lang, "gentoken_description")
}
