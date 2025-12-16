package gentoken

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares/middleware"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/google/uuid"
)

func begin(
	logger *slog.Logger,
	i18n I18n.I18nInterface,
) middlewares.TgMiddlewareHandler {
	return func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) error {
		msg := ectx.EffectiveMessage

		logger.Debug(
			"generate UUID token command",
			slog.Int64("chat_id", msg.Chat.Id),
			slog.String("from", msg.From.Username),
			slog.String("message", msg.Text),
		)

		lang, ok := ctx.Value(middleware.LangKey).(string)
		if !ok {
			return middleware.ErrMissingLangMiddleware
		}

		isSilenced, ok := ctx.Value(middleware.SilenceKey).(bool)
		if !ok {
			return middleware.ErrMissingSilenceMiddleware
		}

		token := uuid.New()
		message := i18n.
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
			T(
				lang,
				"gentoken_generated",
				I18n.WithArgs([]any{
					token.String(),
				}),
			).
			String()

		_, err := b.SendMessage(
			msg.Chat.Id,
			message,
			&gotgbot.SendMessageOpts{
				DisableNotification: isSilenced,
				ParseMode:           "html",
			},
		)
		return err
	}
}
