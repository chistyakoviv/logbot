package setlang

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middleware"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/service/commands"
	"github.com/chistyakoviv/logbot/internal/service/user_settings"
)

const CommandName = "setlang"
const SetLangCbName = "setlang"
const langParam = "lang"

func New(
	ctx context.Context,
	mw middleware.TgMiddlewareInterface,
	logger *slog.Logger,
	i18n *I18n.I18n,
	commands commands.IService,
	userSettings user_settings.IService,
) *command.TgCommand {
	return &command.TgCommand{
		Handler: mw.Pipe(begin(ctx, logger, i18n, commands)).Handler(ctx),
		Stages:  []handlers.Response{},
		Callbacks: map[string]handlers.Response{
			SetLangCbName: mw.Pipe(setlangCb(ctx, logger, i18n, commands, userSettings)).Handler(ctx),
		},
	}
}
