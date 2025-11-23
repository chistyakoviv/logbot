package subscribe

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/rbac"
	"github.com/chistyakoviv/logbot/internal/service/commands"
	"github.com/chistyakoviv/logbot/internal/service/subscriptions"
	"github.com/chistyakoviv/logbot/internal/service/user_settings"
)

const CommandName = "subscribe"

func New(
	ctx context.Context,
	logger *slog.Logger,
	i18n *i18n.I18n,
	rbac rbac.ManagerInterface,
	subscriptions subscriptions.IService,
	commands commands.IService,
	userSettings user_settings.IService,
) *command.TgCommand {
	return &command.TgCommand{
		Handler: begin(ctx, logger, i18n, rbac, commands, userSettings),
		Stages: []handlers.Response{
			stage0(ctx, logger, i18n, subscriptions, commands, userSettings),
		},
	}
}
