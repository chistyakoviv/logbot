package unsubscribe

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/service/commands"
	"github.com/chistyakoviv/logbot/internal/service/subscriptions"
)

const CommandName = "unsubscribe"

type unsubscribeCommand struct {
	command.TgCommand
}

func New(
	ctx context.Context,
	mw middlewares.TgMiddlewareChainInterface,
	logger *slog.Logger,
	i18n I18n.I18nInterface,
	subscriptions subscriptions.ServiceInterface,
	commands commands.ServiceInterface,
) command.TgCommandInterface {
	return &unsubscribeCommand{
		TgCommand: command.TgCommand{
			StageHandlers: []handlers.Response{
				mw.Handler(ctx, requestToken(logger, i18n, commands)),
				mw.Handler(ctx, unsubscribe(logger, i18n, subscriptions, commands)),
			},
		},
	}
}

func (c *unsubscribeCommand) ApplyDescription(lang string, i18n I18n.I18nChainInterface) {
	i18n.
		Appendf("\n\n/%s - ", CommandName).
		T(lang, "unsubscribe_description")
}
