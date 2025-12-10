package subscribe

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares/middleware"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/chistyakoviv/logbot/internal/service/commands"
)

func begin(
	logger *slog.Logger,
	i18n I18n.I18nInterface,
	commands commands.ServiceInterface,
) middlewares.TgMiddlewareHandler {
	return func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) (context.Context, error) {
		msg := ectx.EffectiveMessage

		logger.Debug(
			"subscribe command: initiate",
			slog.Int64("chat_id", msg.Chat.Id),
			slog.String("from", msg.From.Username),
		)

		lang, ok := ctx.Value(middleware.LangKey).(string)
		if !ok {
			return ctx, middleware.ErrMissingLangMiddleware
		}

		var err error

		_, err = commands.ResetByKey(
			ctx,
			&model.CommandKey{
				ChatId: msg.Chat.Id,
				UserId: msg.From.Id,
			},
			CommandName,
			nil,
		)
		if err != nil {
			logger.Error("error occurred while subscribing", slogger.Err(err))
			_, err = b.SendMessage(
				msg.Chat.Id,
				i18n.
					Chain().
					T(
						lang,
						"mention",
						I18n.WithArgs([]any{
							msg.From.Id,
							msg.From.Username,
						}),
					).
					Append("\n").
					T(lang, "subscribe_error").
					String(),
				&gotgbot.SendMessageOpts{
					ParseMode: "html",
				},
			)
			return ctx, err
		}

		_, err = b.SendMessage(
			msg.Chat.Id,
			i18n.
				Chain().
				T(
					lang,
					"mention",
					I18n.WithArgs([]any{
						msg.From.Id,
						msg.From.Username,
					}),
				).
				Append("\n").
				T(lang, "subscribe_begin").
				String(),
			&gotgbot.SendMessageOpts{
				ParseMode: "html",
			},
		)
		return ctx, err
	}
}
