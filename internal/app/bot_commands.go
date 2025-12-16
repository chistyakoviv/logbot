package app

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/addlabels"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/cancel"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/collapse"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/gentoken"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/help"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/labels"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/mute"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/rmlabels"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/setlang"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/silence"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/start"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/subscribe"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/subscriptions"
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
	srvChatSettings := resolveChatSettingsService(c)
	srvLabels := resolveLabelsService(c)
	mw := resolveTgMiddleware(c)

	// Middlewares
	mwLang := resolveTgLangMiddleware(c)
	mwSuperuser := resolveTgSuperuserMiddleware(c)
	mwSubscription := resolveTgSubscriptionMiddleware(c)
	mwSilence := resolveTgSilenceMiddleware(c)
	// gotgbot intercepts panics before the router recoverer,
	// so a custom recoverer is required.
	mwRecoverer := resolveTgRecovererMiddleware(c)

	// Lang middleware must be the first, silence must be the second
	mw = mw.
		Pipe(mwRecoverer).
		Pipe(mwLang).
		Pipe(mwSilence)
	tgCommands := command.TgCommands{
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
			mw.Pipe(mwSuperuser),
			logger,
			i18n,
			srvSubscriptions,
			srvCommands,
		),
		unsubscribe.CommandName: unsubscribe.New(
			ctx,
			mw.Pipe(mwSuperuser),
			logger,
			i18n,
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
			mw.Pipe(mwSubscription),
			logger,
			i18n,
			srvLabels,
			srvCommands,
		),
		rmlabels.CommandName: rmlabels.New(
			ctx,
			mw.Pipe(mwSubscription),
			logger,
			i18n,
			srvLabels,
			srvCommands,
		),
		labels.CommandName: labels.New(
			ctx,
			mw,
			logger,
			i18n,
			srvLabels,
		),
		collapse.CommandName: collapse.New(
			ctx,
			mw.Pipe(mwSubscription),
			logger,
			i18n,
			srvCommands,
			srvChatSettings,
		),
		mute.CommandName: mute.New(
			ctx,
			mw.Pipe(mwSubscription),
			logger,
			i18n,
			srvCommands,
			srvChatSettings,
		),
		gentoken.CommandName: gentoken.New(
			ctx,
			mw,
			logger,
			i18n,
		),
		subscriptions.CommandName: subscriptions.New(
			ctx,
			mw,
			logger,
			i18n,
			srvSubscriptions,
		),
		silence.CommandName: silence.New(
			ctx,
			mw.Pipe(mwSubscription),
			logger,
			i18n,
			srvCommands,
			srvChatSettings,
		),
	}

	// The help command requires all commands as a dependency,
	// but the help command itself is command, so commands cannot be
	// constructed before the help command.
	// Use a little trick to overcome this problem by constructing
	// the help command after all commands are constructed.
	tgCommands[help.CommandName] = help.New(
		ctx,
		mw,
		logger,
		i18n,
		tgCommands,
	)

	return tgCommands
}
