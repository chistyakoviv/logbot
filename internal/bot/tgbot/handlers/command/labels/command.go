package labels

import (
	"context"
	"log/slog"

	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	"github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/service/labels"
)

const CommandName = "labels"

type labelsCommand struct {
	command.TgCommand
}

func New(
	ctx context.Context,
	mw middlewares.TgMiddlewareInterface,
	logger *slog.Logger,
	i18n i18n.I18nInterface,
	labels labels.ServiceInterface,
) command.TgCommandInterface {
	return &labelsCommand{
		TgCommand: command.TgCommand{
			Handler: mw.Pipe(begin(logger, i18n, labels)).Handler(ctx),
		},
	}
}

func (c *labelsCommand) ApplyDescription(lang string, i18n i18n.I18nChainInterface) {
	i18n.
		Appendf("\n\n/%s - ", CommandName).
		T(lang, "labels_description")
}
