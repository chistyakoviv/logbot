package labels

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares/middleware"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
	"github.com/chistyakoviv/logbot/internal/service/labels"
)

func begin(
	logger *slog.Logger,
	i18n I18n.I18nInterface,
	labelsService labels.ServiceInterface,
) middlewares.TgMiddlewareHandler {
	return func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) (context.Context, error) {
		msg := ectx.EffectiveMessage

		logger.Debug(
			"Show labels command",
			slog.Int64("chat_id", msg.Chat.Id),
			slog.String("from", msg.From.Username),
		)

		lang, ok := ctx.Value(middleware.LangKey).(string)
		if !ok {
			return ctx, middleware.ErrMissingLangMiddleware
		}

		isSilenced, ok := ctx.Value(middleware.SilenceKey).(bool)
		if !ok {
			return ctx, middleware.ErrMissingSilenceMiddleware
		}

		entries, err := labelsService.FindAllByChat(ctx, msg.Chat.Id)
		if err != nil {
			logger.Error("error occurred while retrieving labels", slogger.Err(err))
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
					T(lang, "labels_error").
					String(),
				&gotgbot.SendMessageOpts{
					DisableNotification: isSilenced,
					ParseMode:           "html",
				},
			)
			return ctx, err
		}

		if len(entries) == 0 {
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
					T(lang, "labels_empty").
					String(),
				&gotgbot.SendMessageOpts{
					DisableNotification: isSilenced,
					ParseMode:           "html",
				},
			)
			return ctx, err
		}

		messageBuilder := i18n.
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
			T(lang, "labels_assigned").
			Append("\n\n")

		for _, row := range entries {
			if len(row.Labels) == 0 {
				continue
			}
			messageBuilder.Appendf("%s: ", row.Username)
			for i, label := range row.Labels {
				messageBuilder.Appendf("<code>%s</code>", label)
				if i < len(row.Labels)-1 {
					messageBuilder.Append(", ")
				}
			}
			messageBuilder.Append("\n\n")
		}

		_, err = b.SendMessage(
			msg.Chat.Id,
			messageBuilder.String(),
			&gotgbot.SendMessageOpts{
				DisableNotification: isSilenced,
				ParseMode:           "html",
			},
		)

		return ctx, err
	}
}
