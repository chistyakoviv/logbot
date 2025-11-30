package app

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/addlabels"
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
	srvLabels := resolveLabelsService(c)
	rbac := resolveRbac(c)
	mw := resolveTgMiddleware(c)

	// Middlewares
	mwLang := resolveTgLangMiddleware(c)

	// Lang middleware must be the first
	mw = mw.Pipe(mwLang)
	return command.TgCommands{
		start.CommandName: start.New(
			ctx,
			mw,
			logger,
			i18n,
		),
		cancel.CommandName: cancel.New(
			ctx,
			mw,
			logger,
			i18n,
			srvCommands,
		),
		subscribe.CommandName: subscribe.New(
			ctx,
			mw,
			logger,
			i18n,
			rbac,
			srvSubscriptions,
			srvCommands,
		),
		unsubscribe.CommandName: unsubscribe.New(
			ctx,
			mw,
			logger,
			i18n,
			rbac,
			srvSubscriptions,
			srvCommands,
		),
		setlang.CommandName: setlang.New(
			ctx,
			mw,
			logger,
			i18n,
			srvCommands,
			srvUserSettings,
		),
		addlabels.CommandName: addlabels.New(
			ctx,
			mw,
			logger,
			i18n,
			srvLabels,
			srvCommands,
		),
	}
}
