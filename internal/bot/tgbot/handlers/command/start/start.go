package start

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middleware"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middleware/middlewares"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
)

const CommandName = "start"

func New(
	ctx context.Context,
	mw middleware.TgMiddlewareInterface,
	logger *slog.Logger,
	i18n *I18n.I18n,
) *command.TgCommand {
	return &command.TgCommand{
		Handler: mw.Pipe(begin(logger, i18n)).Handler(ctx),
		Stages:  []handlers.Response{},
	}
}

func begin(
	logger *slog.Logger,
	i18n *I18n.I18n,
) middleware.TgMiddlewareHandler {
	return func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) (context.Context, error) {
		msg := ectx.EffectiveMessage

		logger.Debug(
			"message received",
			slog.Int64("chat_id", msg.Chat.Id),
			slog.String("from", msg.From.Username),
			slog.String("message", msg.Text),
		)

		lang, ok := ctx.Value(middlewares.LangKey).(string)
		if !ok {
			return ctx, middlewares.ErrMissingLangMiddleware
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
