package subscribe

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middleware"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middleware/middlewares"
	"github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/rbac"
	"github.com/chistyakoviv/logbot/internal/service/commands"
	"github.com/chistyakoviv/logbot/internal/service/subscriptions"
)

const CommandName = "subscribe"

func New(
	ctx context.Context,
	mw middleware.TgMiddlewareInterface,
	logger *slog.Logger,
	i18n *i18n.I18n,
	rbac rbac.ManagerInterface,
	subscriptions subscriptions.IService,
	commands commands.IService,
) *command.TgCommand {
	mw = mw.Pipe(middlewares.Superuser(logger, i18n, rbac))
	return &command.TgCommand{
		Handler: mw.Pipe(begin(logger, i18n, commands)).Handler(ctx),
		Stages: []handlers.Response{
			mw.Pipe(stage0(logger, i18n, subscriptions, commands)).Handler(ctx),
		},
	}
}
