package setlang

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/service/commands"
	"github.com/chistyakoviv/logbot/internal/service/user_settings"
)

const CommandName = "setlang"
const SetLangCbName = "setlang"
const langParam = "lang"

type setlangCommand struct {
	command.TgCommand
}

func New(
	ctx context.Context,
	mw middlewares.TgMiddlewareInterface,
	logger *slog.Logger,
	i18n I18n.I18nInterface,
	commands commands.ServiceInterface,
	userSettings user_settings.ServiceInterface,
) command.TgCommandInterface {
	return &setlangCommand{
		TgCommand: command.TgCommand{
			Handler: mw.Pipe(begin(logger, i18n)).Handler(ctx),
			Callbacks: map[string]handlers.Response{
				SetLangCbName: mw.Pipe(setlangCb(logger, i18n, userSettings)).Handler(ctx),
			},
		},
	}
}

func (c *setlangCommand) ApplyDescription(lang string, i18n I18n.I18nChainInterface) {
	i18n.
		Appendf("\n\n/%s - ", CommandName).
		T(lang, "setlang_description")
}
