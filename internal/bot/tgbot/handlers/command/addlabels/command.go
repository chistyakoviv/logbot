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
	stageMentions = iota
	stageLabels
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
			StartHandler: mw.Handler(ctx, begin(logger, i18n, commands)),
			StageHandlers: []handlers.Response{
				mw.Handler(ctx, stage0(logger, i18n, commands)),
				mw.Handler(ctx, stage1(logger, i18n, labels, commands)),
			},
		},
	}
}

func (c *addlabelsCommand) ApplyDescription(lang string, i18n I18n.I18nChainInterface) {
	i18n.
		Appendf("\n\n/%s - ", CommandName).
		T(lang, "addlabels_description")
}
