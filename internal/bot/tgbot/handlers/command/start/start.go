package start

import (
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
)

const CommandName = "start"

func New(logger *slog.Logger, i18n *I18n.I18n) *command.TgCommand {
	return &command.TgCommand{
		Handler: begin(logger, i18n),
		Stages:  []handlers.Response{},
	}
}

func begin(logger *slog.Logger, i18n *I18n.I18n) handlers.Response {
	lang := i18n.DefaultLang()
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		msg := ctx.EffectiveMessage

		logger.Debug(
			"message received",
			slog.Int64("chat_id", msg.Chat.Id),
			slog.String("from", msg.From.Username),
			slog.String("message", msg.Text),
		)

		message := i18n.
			Chain().
			T(lang, "greeting").
			Append("\n\n").
			T(lang, "description").
			Append("\n\n").
			T(lang, "intro").
			Append("\n\n").
			T(lang, "help").
			String()
		// fmt.Fprintf(&message, "%s", "<pre language=\"typescript\">console.log('Hello, world!')</pre>")
		_, err := b.SendMessage(msg.Chat.Id, message, &gotgbot.SendMessageOpts{
			ParseMode: "html",
		})
		return err
	}
}
