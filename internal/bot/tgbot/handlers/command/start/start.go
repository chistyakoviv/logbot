package start

import (
	"bytes"
	"fmt"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/i18n"
)

func New(logger *slog.Logger, i18n *i18n.I18n) *command.TgCommand {
	return &command.TgCommand{
		Handler: begin(logger, i18n),
		Stages:  []handlers.Response{},
	}
}

func begin(logger *slog.Logger, i18n *i18n.I18n) handlers.Response {
	lang := i18n.DefaultLang()
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		// _, err := ctx.EffectiveMessage.Reply(b, "👋 Welcome! I’m your Go webhook bot.\nUse /help for commands.", nil)
		msg := ctx.EffectiveMessage

		logger.Debug(
			"message received",
			slog.Int64("chat_id", msg.Chat.Id),
			slog.String("from", msg.From.Username),
			slog.String("message", msg.Text),
		)

		var message bytes.Buffer
		fmt.Fprintf(&message, "%s\n\n", i18n.T(lang, "greeting"))
		fmt.Fprintf(&message, "%s\n\n", i18n.T(lang, "description"))
		fmt.Fprintf(&message, "%s\n\n", i18n.T(lang, "intro"))
		fmt.Fprintf(&message, "%s", i18n.T(lang, "help"))
		// fmt.Fprintf(&message, "%s", "<pre language=\"typescript\">console.log('Hello, world!')</pre>")
		_, err := b.SendMessage(msg.Chat.Id, message.String(), &gotgbot.SendMessageOpts{
			ParseMode: "html",
		})
		return err
	}
}
