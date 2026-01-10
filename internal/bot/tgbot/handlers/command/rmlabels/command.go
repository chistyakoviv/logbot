package rmlabels

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/service/commands"
	"github.com/chistyakoviv/logbot/internal/service/labels"
)

const CommandName = "rmlabels"

const (
	stageStart = iota
	stageAcceptUsers
	stageRemoveLabels
)

type rmlabelsCommand struct {
	command.TgCommand
}

func New(
	ctx context.Context,
	mw middlewares.TgMiddlewareChainInterface,
	logger *slog.Logger,
	i18n I18n.I18nInterface,
	labels labels.ServiceInterface,
	commands commands.ServiceInterface,
) command.TgCommandInterface {
	return &rmlabelsCommand{
		TgCommand: command.TgCommand{
			StageHandlers: []handlers.Response{
				mw.Handler(ctx, requestUsers(logger, i18n, commands)),
				mw.Handler(
					ctx,
					command.
						NewCommandChain().
						Add(acceptUsers(logger, i18n, commands)).
						Add(requestLabels(logger, i18n)).
						Build(),
				),
				mw.Handler(ctx, removeLabels(logger, i18n, labels, commands)),
			},
		},
	}
}

func (c *rmlabelsCommand) ApplyDescription(lang string, i18n I18n.I18nChainInterface) {
	i18n.
		Appendf("\n\n/%s - ", CommandName).
		T(lang, "rmlabels_description")
}
