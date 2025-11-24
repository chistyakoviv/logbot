package unsubscribe

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middleware"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middleware/middlewares"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/rbac"
	"github.com/chistyakoviv/logbot/internal/service/commands"
	"github.com/chistyakoviv/logbot/internal/service/subscriptions"
)

const CommandName = "unsubscribe"

func New(
	ctx context.Context,
	mw middleware.TgMiddlewareInterface,
	logger *slog.Logger,
	i18n I18n.I18nInterface,
	rbac rbac.ManagerInterface,
	subscriptions subscriptions.ServiceInterface,
	commands commands.ServiceInterface,
) *command.TgCommand {
	mw = mw.Pipe(middlewares.Superuser(logger, i18n, rbac))
	return &command.TgCommand{
		Handler: mw.Pipe(begin(logger, i18n, commands)).Handler(ctx),
		Stages: []handlers.Response{
			mw.Pipe(stage0(logger, i18n, subscriptions, commands)).Handler(ctx),
		},
	}
}
