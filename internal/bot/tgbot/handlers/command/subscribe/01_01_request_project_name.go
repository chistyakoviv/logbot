package subscribe

import (
	"context"
	"log/slog"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares/middleware"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/model"
)

func requestProjectName(
	logger *slog.Logger,
	i18n I18n.I18nInterface,
) middlewares.TgMiddlewareHandler {
	return func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) error {
		msg := ectx.EffectiveMessage
		token := strings.Trim(msg.Text, " ")

		logger.Debug(
			"subscribe command: request project name",
			slog.Int64("chat_id", msg.Chat.Id),
			slog.String("user", msg.From.Username),
			slog.String("token", token),
		)

		lang, ok := ctx.Value(middleware.LangKey).(string)
		if !ok {
			return middleware.ErrMissingLangMiddleware
		}

		isSilenced, ok := ctx.Value(middleware.SilenceKey).(bool)
		if !ok {
			return middleware.ErrMissingSilenceMiddleware
		}

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
				T(
					lang,
					"subscribe_enter_project_name",
					I18n.WithArgs([]any{
						model.MaxProjectNameLength,
					}),
				).
				String(),
			&gotgbot.SendMessageOpts{
				DisableNotification: isSilenced,
				ParseMode:           "html",
			},
		)
		return err
	}
}
