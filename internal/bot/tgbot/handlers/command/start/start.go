package start

import (
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/i18n"
)

func New(logger *slog.Logger, i18n *i18n.I18n) *command.TgCommand {
	return &command.TgCommand{
		Handler: stage0(logger, i18n),
		Stages:  []handlers.Response{},
	}
}

func stage0(logger *slog.Logger, i18n *i18n.I18n) handlers.Response {
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
		_, err := b.SendMessage(msg.Chat.Id, i18n.T("en", "greeting"), nil)
		return err
	}
}
