package labels

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	"github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/service/labels"
)

const CommandName = "labels"

func New(
	ctx context.Context,
	mw middlewares.TgMiddlewareInterface,
	logger *slog.Logger,
	i18n i18n.I18nInterface,
	labels labels.ServiceInterface,
) *command.TgCommand {
	return &command.TgCommand{
		Handler: mw.Pipe(begin(logger, i18n, labels)).Handler(ctx),
		Stages:  []handlers.Response{},
	}
}
