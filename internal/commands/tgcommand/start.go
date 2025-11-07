package tgcommand

import (
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func Start(logger *slog.Logger) handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		// _, err := ctx.EffectiveMessage.Reply(b, "👋 Welcome! I’m your Go webhook bot.\nUse /help for commands.", nil)
		msg := ctx.EffectiveMessage

		logger.Info(
			"message received",
			slog.Int64("chat_id", msg.Chat.Id),
			slog.String("from", msg.From.Username),
			slog.String("message", msg.Text),
		)

		// Send a new message instead of replying
		_, err := b.SendMessage(msg.Chat.Id, "👋 Welcome! I’m your Go webhook bot.\nUse /help for commands.", nil)
		return err
	}
}
