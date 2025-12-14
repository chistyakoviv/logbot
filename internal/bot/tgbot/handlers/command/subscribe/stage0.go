package subscribe

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares/middleware"
	"github.com/chistyakoviv/logbot/internal/db"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/chistyakoviv/logbot/internal/service/commands"
	"github.com/chistyakoviv/logbot/internal/service/subscriptions"
	"github.com/google/uuid"
)

func stage0(
	logger *slog.Logger,
	i18n I18n.I18nInterface,
	subscriptions subscriptions.ServiceInterface,
	commands commands.ServiceInterface,
) middlewares.TgMiddlewareHandler {
	return func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) (context.Context, error) {
		msg := ectx.EffectiveMessage
		token := strings.Trim(msg.Text, " ")

		logger.Debug(
			"subscribe command: retrieve token",
			slog.Int64("chat_id", msg.Chat.Id),
			slog.String("user", msg.From.Username),
			slog.String("token", token),
		)

		lang, ok := ctx.Value(middleware.LangKey).(string)
		if !ok {
			return ctx, middleware.ErrMissingLangMiddleware
		}

		isSilenced, ok := ctx.Value(middleware.SilenceKey).(bool)
		if !ok {
			return ctx, middleware.ErrMissingSilenceMiddleware
		}

		var err error
		if err := uuid.Validate(token); err != nil {
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
					T(lang, "subscribe_invalid_token").
					String(),
				&gotgbot.SendMessageOpts{
					DisableNotification: isSilenced,
					ParseMode:           "html",
				},
			)
			return ctx, err
		}
		_, subErr := subscriptions.Find(ctx, token, msg.Chat.Id)
		if subErr == nil {
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
					T(lang, "subscribe_token_exists").
					String(),
				&gotgbot.SendMessageOpts{
					DisableNotification: isSilenced,
					ParseMode:           "html",
				},
			)
			return ctx, err
		}

		_, err = commands.UpdateByKey(
			ctx,
			&model.CommandKey{
				ChatId: msg.Chat.Id,
				UserId: msg.From.Id,
			},
			stageProjectName,
			map[string]any{
				"token": token,
			},
		)
		if err != nil || !errors.Is(subErr, db.ErrNotFound) {
			if err == nil {
				err = subErr
			}
			logger.Error("error occurred while subscribing", slogger.Err(err))
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
					T(lang, "subscribe_error").
					String(),
				&gotgbot.SendMessageOpts{
					DisableNotification: isSilenced,
					ParseMode:           "html",
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
		return ctx, err
	}
}
