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
	stageMentions = iota
	stageLabels
)

type rmlabelsCommand struct {
	command.TgCommand
}

func New(
	ctx context.Context,
	mw middlewares.TgMiddlewareInterface,
	mwSubscription middlewares.TgMiddlewareHandler,
	logger *slog.Logger,
	i18n I18n.I18nInterface,
	labels labels.ServiceInterface,
	commands commands.ServiceInterface,
) command.TgCommandInterface {
	mw = mw.Pipe(mwSubscription)
	return &rmlabelsCommand{
		TgCommand: command.TgCommand{
			Handler: mw.Pipe(begin(logger, i18n, commands)).Handler(ctx),
			Stages: []handlers.Response{
				mw.Pipe(stage0(logger, i18n, commands)).Handler(ctx),
				mw.Pipe(stage1(logger, i18n, labels, commands)).Handler(ctx),
			},
		},
	}
}

func (c *rmlabelsCommand) ApplyDescription(lang string, i18n I18n.I18nChainInterface) {
	i18n.
		Appendf("\n\n/%s - ", CommandName).
		T(lang, "rmlabels_description")
}
