package start

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares/middleware"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
)

func begin(
	logger *slog.Logger,
	i18n I18n.I18nInterface,
) middlewares.TgMiddlewareHandler {
	return func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) (context.Context, error) {
		msg := ectx.EffectiveMessage

		logger.Debug(
			"message received",
			slog.Int64("chat_id", msg.Chat.Id),
			slog.String("from", msg.From.Username),
			slog.String("message", msg.Text),
		)

		lang, ok := ctx.Value(middleware.LangKey).(string)
		if !ok {
			return ctx, middleware.ErrMissingLangMiddleware
		}

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
		return ctx, err
	}
}
