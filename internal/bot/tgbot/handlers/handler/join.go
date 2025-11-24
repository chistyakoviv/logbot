package handler

import (
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/i18n"
)

func NewJoin(logger *slog.Logger, i18n i18n.I18nInterface) handlers.Response {
	return joinHandler(logger, i18n)
}

func joinHandler(logger *slog.Logger, i18n i18n.I18nInterface) handlers.Response {
	lang := i18n.DefaultLang()
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		msg := ctx.EffectiveMessage

		for _, member := range msg.NewChatMembers {
			if member.Id == b.Id {
				logger.Debug(
					"Bot joined a new chat",
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
				_, err := b.SendMessage(msg.Chat.Id, message, &gotgbot.SendMessageOpts{
					ParseMode: "html",
				})
				return err
			}
		}
		return nil
	}
}
