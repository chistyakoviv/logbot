package rmlabels

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middleware"
	"github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/service/commands"
	"github.com/chistyakoviv/logbot/internal/service/labels"
)

const CommandName = "rmlabels"

const (
	stageMentions = iota
	stageLabels
)

func New(
	ctx context.Context,
	mw middleware.TgMiddlewareInterface,
	logger *slog.Logger,
	i18n i18n.I18nInterface,
	labels labels.ServiceInterface,
	commands commands.ServiceInterface,
) *command.TgCommand {
	return &command.TgCommand{
		Handler: mw.Pipe(begin(logger, i18n, commands)).Handler(ctx),
		Stages: []handlers.Response{
			mw.Pipe(stage0(logger, i18n, commands)).Handler(ctx),
			mw.Pipe(stage1(logger, i18n, labels, commands)).Handler(ctx),
		},
	}
}
