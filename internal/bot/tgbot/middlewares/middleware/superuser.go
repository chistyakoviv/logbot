package middleware

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	"github.com/chistyakoviv/logbot/internal/constants"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/lib/rbac"
	errs "github.com/pkg/errors"
)

func Superuser(logger *slog.Logger, i18n I18n.I18nInterface, rbac rbac.ManagerInterface) middlewares.TgMiddleware {
	fn := func(next middlewares.TgMiddlewareHandler) middlewares.TgMiddlewareHandler {
		return func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) error {
			msg := ectx.EffectiveMessage

			logger.Debug(
				"superuser middleware",
				slog.Int64("chat_id", msg.Chat.Id),
				slog.String("from", msg.From.Username),
			)

			lang, ok := ctx.Value(LangKey).(string)
			if !ok {
				return ErrMissingLangMiddleware
			}

			if !rbac.UserHasPermission(msg.From.Id, constants.PermissionManage, nil) {
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
						T(lang, "access_denied").
						String(),
					&gotgbot.SendMessageOpts{
						ParseMode: "html",
					},
				)
				if err != nil {
					return err
				}
				return errs.Wrap(middlewares.ErrMiddlewareCanceled, "access denied")
			}
			return next(ctx, b, ectx)
		}
	}
	return fn
}
