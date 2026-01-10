package addlabels

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

const CommandName = "addlabels"

const (
	stageStart = iota
	stageAcceptUsers
	stageApplyLabels
)

type addlabelsCommand struct {
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
	return &addlabelsCommand{
		TgCommand: command.TgCommand{
			StageHandlers: []handlers.Response{
				mw.Handler(ctx, reqeustUsers(logger, i18n, commands)),
				mw.Handler(
					ctx,
					command.
						NewCommandChain().
						Add(acceptUsers(logger, i18n, commands)).
						Add(requestLabels(logger, i18n)).
						Build(),
				),
				mw.Handler(ctx, applyLabels(logger, i18n, labels, commands)),
			},
		},
	}
}

func (c *addlabelsCommand) ApplyDescription(lang string, i18n I18n.I18nChainInterface) {
	i18n.
		Appendf("\n\n/%s - ", CommandName).
		T(lang, "addlabels_description")
}
