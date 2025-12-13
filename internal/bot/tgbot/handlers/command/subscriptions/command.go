package subscriptions

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/service/subscriptions"
)

const CommandName = "subscriptions"

func New(
	ctx context.Context,
	mw middlewares.TgMiddlewareInterface,
	logger *slog.Logger,
	i18n I18n.I18nInterface,
	subscriptions subscriptions.ServiceInterface,
) *command.TgCommand {
	return &command.TgCommand{
		Handler:   mw.Pipe(begin(logger, i18n, subscriptions)).Handler(ctx),
		Stages:    []handlers.Response{},
		Callbacks: map[string]handlers.Response{},
	}
}
