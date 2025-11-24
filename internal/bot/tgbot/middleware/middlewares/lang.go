package middlewares

import (
	"context"
	"errors"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middleware"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
	"github.com/chistyakoviv/logbot/internal/service/user_settings"
)

const (
	LangKey key = "lang"
)

var (
	ErrMissingLangMiddleware error = errors.New("missing lang middleware")
)

func Lang(logger *slog.Logger, userSettings user_settings.IService) middleware.TgMiddlewareHandler {
	return func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) (context.Context, error) {
		msg := ectx.EffectiveMessage

		logger.Debug(
			"lang middleware",
			slog.Int64("chat_id", msg.Chat.Id),
			slog.String("from", msg.From.Username),
		)

		lang, err := userSettings.GetLang(ctx, msg.From.Id)
		if err != nil && !errors.Is(err, db.ErrNotFound) {
			logger.Error("error occurred while getting the user's language", slogger.Err(err))
			_, err := b.SendMessage(
				msg.Chat.Id,
				"Failed to get the user's language. Please check the log for more information.",
				&gotgbot.SendMessageOpts{
					ParseMode: "html",
				},
			)
			return ctx, err
		}

		ctx = context.WithValue(ctx, LangKey, lang)
		return ctx, nil
	}
}
