package subscribe

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	"github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/service/commands"
	"github.com/chistyakoviv/logbot/internal/service/subscriptions"
)

const CommandName = "subscribe"

const (
	stageStart = iota
	stageAcceptToken
	stageAcceptSubscription
)

type subscribeCommand struct {
	command.TgCommand
}

func New(
	ctx context.Context,
	mw middlewares.TgMiddlewareChainInterface,
	logger *slog.Logger,
	i18n i18n.I18nInterface,
	subscriptions subscriptions.ServiceInterface,
	commands commands.ServiceInterface,
) command.TgCommandInterface {
	return &subscribeCommand{
		TgCommand: command.TgCommand{
			StageHandlers: []handlers.Response{
				mw.Handler(ctx, requestToken(logger, i18n, commands)),
				mw.Handler(
					ctx,
					command.
						NewCommandChain().
						Add(acceptToken(logger, i18n, subscriptions, commands)).
						Add(requestProjectName(logger, i18n)).
						Build(),
				),
				mw.Handler(ctx, acceptSubscription(logger, i18n, subscriptions, commands)),
			},
		},
	}
}

func (c *subscribeCommand) ApplyDescription(lang string, i18n i18n.I18nChainInterface) {
	i18n.
		Appendf("\n\n/%s - ", CommandName).
		T(lang, "subscribe_description")
}
