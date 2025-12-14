package silence

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

const CommandName = "silence"
const silenceCbName = "silence"
const silencePeriodParam = "period"

type silencePeriod struct {
	Label    string
	Duration time.Duration
}

var periods = []silencePeriod{
	{"none", time.Second * 0},
	{"5 minutes", time.Minute * 5},
	{"10 minutes", time.Minute * 10},
	{"30 minutes", time.Minute * 30},
	{"1 hour", time.Hour * 1},
	{"3 hours", time.Hour * 3},
	{"12 hours", time.Hour * 12},
	{"24 hours", time.Hour * 24},
}

type silenceCommand struct {
	command.TgCommand
}

func New(
	ctx context.Context,
	mw middlewares.TgMiddlewareInterface,
	mwSubscription middlewares.TgMiddlewareHandler,
	logger *slog.Logger,
	i18n I18n.I18nInterface,
	commands commands.ServiceInterface,
	chatSettings chat_settings.ServiceInterface,
) command.TgCommandInterface {
	mw = mw.Pipe(mwSubscription)
	return &silenceCommand{
		TgCommand: command.TgCommand{
			Handler: mw.Pipe(begin(logger, i18n)).Handler(ctx),
			Callbacks: map[string]handlers.Response{
				silenceCbName: mw.Pipe(silenceCb(logger, i18n, chatSettings)).Handler(ctx),
			},
		},
	}
}

func (c *silenceCommand) ApplyDescription(lang string, i18n I18n.I18nChainInterface) {
	i18n.
		Appendf("\n\n/%s - ", CommandName).
		T(lang, "silence_description")
}
