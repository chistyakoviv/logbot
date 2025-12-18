package help

import (
	"context"
	"log/slog"
	"sort"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares/middleware"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
)

func begin(
	logger *slog.Logger,
	i18n I18n.I18nInterface,
	tgCommands command.TgCommands,
) middlewares.TgMiddlewareHandler {
	// Output commands in lexicographic order
	commandNames := make([]string, 0, len(tgCommands)+1)
	for name := range tgCommands {
		commandNames = append(commandNames, name)
	}
	// Trick to add itself to the list of commands
	commandNames = append(commandNames, CommandName)
	sort.Strings(commandNames)
	return func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) error {
		msg := ectx.EffectiveMessage

		logger.Debug(
			"help command",
			slog.Int64("chat_id", msg.Chat.Id),
			slog.String("from", msg.From.Username),
		)

		lang, ok := ctx.Value(middleware.LangKey).(string)
		if !ok {
			return middleware.ErrMissingLangMiddleware
		}

		isSilenced, ok := ctx.Value(middleware.SilenceKey).(bool)
		if !ok {
			return middleware.ErrMissingSilenceMiddleware
		}

		messageBuilder := i18n.
			Chain().
			T(lang, "help_title")

		for _, cmdName := range commandNames {
			tgCommands[cmdName].ApplyDescription(lang, messageBuilder)
		}

		_, err := b.SendMessage(
			msg.Chat.Id,
			messageBuilder.String(),
			&gotgbot.SendMessageOpts{
				DisableNotification: isSilenced,
				ParseMode:           "html",
			},
		)
		return err
	}
}
