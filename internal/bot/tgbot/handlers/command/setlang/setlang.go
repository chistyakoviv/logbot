package setlang

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/service/commands"
	"github.com/chistyakoviv/logbot/internal/service/user_settings"
)

const CommandName = "setlang"
const SetLangCbName = "setlang"
const langParam = "lang"

func New(
	ctx context.Context,
	logger *slog.Logger,
	i18n *I18n.I18n,
	commands commands.IService,
	userSettings user_settings.IService,
) *command.TgCommand {
	return &command.TgCommand{
		Handler: begin(ctx, logger, i18n, commands, userSettings),
		Stages:  []handlers.Response{},
		Callbacks: map[string]handlers.Response{
			SetLangCbName: setlangCb(ctx, logger, i18n, commands, userSettings),
		},
	}
}
