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
	stageToken = iota
	stageProjectName
)

type subscribeCommand struct {
	command.TgCommand
}

func New(
	ctx context.Context,
	mw middlewares.TgMiddlewareInterface,
	mwSuperuser middlewares.TgMiddlewareHandler,
	logger *slog.Logger,
	i18n i18n.I18nInterface,
	subscriptions subscriptions.ServiceInterface,
	commands commands.ServiceInterface,
) command.TgCommandInterface {
	mw = mw.Pipe(mwSuperuser)
	return &subscribeCommand{
		TgCommand: command.TgCommand{
			Handler: mw.Pipe(begin(logger, i18n, commands)).Handler(ctx),
			Stages: []handlers.Response{
				mw.Pipe(stage0(logger, i18n, subscriptions, commands)).Handler(ctx),
				mw.Pipe(stage1(logger, i18n, subscriptions, commands)).Handler(ctx),
			},
		},
	}
}

func (c *subscribeCommand) ApplyDescription(lang string, i18n i18n.I18nChainInterface) {
	i18n.
		Appendf("\n\n/%s - ", CommandName).
		T(lang, "subscribe_description")
}
