package subscriptions

import (
	"context"
	"log/slog"

	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/service/subscriptions"
)

const CommandName = "subscriptions"

type subscriptionsCommand struct {
	command.TgCommand
}

func New(
	ctx context.Context,
	mw middlewares.TgMiddlewareInterface,
	logger *slog.Logger,
	i18n I18n.I18nInterface,
	subscriptions subscriptions.ServiceInterface,
) command.TgCommandInterface {
	return &subscriptionsCommand{
		TgCommand: command.TgCommand{
			Handler: mw.Pipe(begin(logger, i18n, subscriptions)).Handler(ctx),
		},
	}
}

func (c *subscriptionsCommand) ApplyDescription(lang string, i18n I18n.I18nChainInterface) {
	i18n.
		Appendf("\n\n/%s - ", CommandName).
		T(lang, "subscriptions_description")
}
