package start

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
)

const CommandName = "start"

type startCommand struct {
	command.TgCommand
}

func New(
	ctx context.Context,
	mw middlewares.TgMiddlewareChainInterface,
	logger *slog.Logger,
	i18n I18n.I18nInterface,
) command.TgCommandInterface {
	return &startCommand{
		TgCommand: command.TgCommand{
			StageHandlers: []handlers.Response{
				mw.Handler(ctx, start(logger, i18n)),
			},
		},
	}
}
