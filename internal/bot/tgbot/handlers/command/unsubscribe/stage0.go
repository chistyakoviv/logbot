package unsubscribe

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/db"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/chistyakoviv/logbot/internal/service/commands"
	"github.com/chistyakoviv/logbot/internal/service/subscriptions"
	"github.com/chistyakoviv/logbot/internal/service/user_settings"
	"github.com/google/uuid"
)

func stage0(
	ctx context.Context,
	logger *slog.Logger,
	i18n *I18n.I18n,
	subscriptions subscriptions.IService,
	commands commands.IService,
	userSettings user_settings.IService,
) handlers.Response {
	return func(b *gotgbot.Bot, ectx *ext.Context) error {
		msg := ectx.EffectiveMessage
		token := strings.Trim(msg.Text, " ")

		logger.Debug(
			"unsubscribe command: retrieve token",
			slog.Int64("chat_id", msg.Chat.Id),
			slog.String("user", msg.From.Username),
			slog.String("token", token),
		)

		lang, err := userSettings.GetLang(ctx, msg.From.Id)
		if err != nil {
			logger.Error("error occurred while getting the user's language", slogger.Err(err))
			_, err := b.SendMessage(
				msg.Chat.Id,
				"Failed to get the user's language. Please check the log for more information.",
				&gotgbot.SendMessageOpts{
					ParseMode: "html",
				},
			)
			return err
		}

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
					T(lang, "unsubscribe_invalid_token").
					String(),
				&gotgbot.SendMessageOpts{
					ParseMode: "html",
				},
			)
			return err
		}
		_, unsubErr := subscriptions.Find(ctx, token, msg.Chat.Id)
		if errors.Is(unsubErr, db.ErrNotFound) {
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
					T(lang, "unsubscribe_token_not_exists").
					String(),
				&gotgbot.SendMessageOpts{
					ParseMode: "html",
				},
			)
			return err
		}

		_, err = commands.CompleteByKey(
			ctx,
			&model.CommandKey{
				ChatId: msg.Chat.Id,
				UserId: msg.From.Id,
			},
		)
		if err != nil || unsubErr != nil {
			logger.Error("error occurred while unsubscribing", slogger.Err(err))
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
					T(lang, "unsubscribe_error").
					String(),
				&gotgbot.SendMessageOpts{
					ParseMode: "html",
				},
			)
			return err
		}

		_, err = subscriptions.Unsubscribe(ctx, token, msg.Chat.Id)
		if err != nil {
			logger.Error("error occurred while unsubscribing", slogger.Err(err))
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
					T(lang, "unsubscribe_error").
					String(),
				&gotgbot.SendMessageOpts{
					ParseMode: "html",
				},
			)
			return err
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
					"unsubscribe_complete",
					I18n.WithArgs([]any{
						token,
					}),
				).
				String(),
			&gotgbot.SendMessageOpts{
				ParseMode: "html",
			},
		)
		return err
	}
}
