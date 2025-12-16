package mute

import (
	"context"
	"log/slog"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/service/chat_settings"
	"github.com/chistyakoviv/logbot/internal/service/commands"
)

const CommandName = "mute"
const muteCbName = "mute"
const mutePeriodParam = "period"

type mutePeriod struct {
	Label    string
	Duration time.Duration
}

var periods = []mutePeriod{
	{"none", time.Second * 0},
	{"5 minutes", time.Minute * 5},
	{"10 minutes", time.Minute * 10},
	{"30 minutes", time.Minute * 30},
	{"1 hour", time.Hour * 1},
	{"3 hours", time.Hour * 3},
	{"6 hours", time.Hour * 6},
	{"12 hours", time.Hour * 12},
	{"24 hours", time.Hour * 24},
	{"3 days", time.Hour * 24 * 3},
}

type muteCommand struct {
	command.TgCommand
}

func New(
	ctx context.Context,
	mw middlewares.TgMiddlewareChainInterface,
	logger *slog.Logger,
	i18n I18n.I18nInterface,
	commands commands.ServiceInterface,
	chatSettings chat_settings.ServiceInterface,
) command.TgCommandInterface {
	return &muteCommand{
		TgCommand: command.TgCommand{
			StartHandler: mw.Handler(ctx, begin(logger, i18n)),
			CallbackHandlers: map[string]handlers.Response{
				muteCbName: mw.Handler(ctx, muteCb(logger, i18n, chatSettings)),
			},
		},
	}
}

func (c *muteCommand) ApplyDescription(lang string, i18n I18n.I18nChainInterface) {
	i18n.
		Appendf("\n\n/%s - ", CommandName).
		T(lang, "mute_description")
}
