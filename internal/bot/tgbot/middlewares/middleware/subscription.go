package middleware

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/service/subscriptions"
	errs "github.com/pkg/errors"
)

func Subscription(
	logger *slog.Logger,
	i18n I18n.I18nInterface,
	subscriptions subscriptions.ServiceInterface,
) middlewares.TgMiddleware {
	fn := func(next middlewares.TgMiddlewareHandler) middlewares.TgMiddlewareHandler {
		return func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) error {
			msg := ectx.EffectiveMessage

			logger.Debug(
				"subscription middleware",
				slog.Int64("chat_id", msg.Chat.Id),
				slog.String("from", msg.From.Username),
			)

			lang, ok := ctx.Value(LangKey).(string)
			if !ok {
				return ErrMissingLangMiddleware
			}

			subs, err := subscriptions.FindByChatId(ctx, msg.Chat.Id)
			if err != nil {
				return err
			}

			if len(subs) == 0 {
				_, err := b.SendMessage(
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
						T(lang, "subscribe_subscription_required").
						String(),
					&gotgbot.SendMessageOpts{
						ParseMode: "html",
					},
				)
				if err != nil {
					return err
				}
				return errs.Wrap(middlewares.ErrMiddlewareCanceled, "subscription required")
			}
			return next(ctx, b, ectx)
		}
	}
	return fn
}
