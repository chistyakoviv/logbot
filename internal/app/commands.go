package app

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/cancel"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/setlang"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/start"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/subscribe"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/unsubscribe"
	"github.com/chistyakoviv/logbot/internal/di"
)

func BuildTgCommands(
	ctx context.Context,
	c di.Container,
) command.TgCommands {
	logger := resolveLogger(c)
	i18n := resolveI18n(c)
	srvCommands := resolveCommandsService(c)
	srvSubscriptions := resolveSubscriptionsService(c)
	srvUserSettings := resolveUserSettingsService(c)
	rbac := resolveRbac(c)
	return command.TgCommands{
		start.CommandName: start.New(logger, i18n),
		cancel.CommandName: cancel.New(
			ctx,
			logger,
			i18n,
			srvCommands,
			srvUserSettings,
		),
		subscribe.CommandName: subscribe.New(
			ctx,
			logger,
			i18n,
			rbac,
			srvSubscriptions,
			srvCommands,
			srvUserSettings,
		),
		unsubscribe.CommandName: unsubscribe.New(
			ctx,
			logger,
			i18n,
			srvSubscriptions,
			srvCommands,
			srvUserSettings,
		),
		setlang.CommandName: setlang.New(
			ctx,
			logger,
			i18n,
			srvCommands,
			srvUserSettings,
		),
	}
}
