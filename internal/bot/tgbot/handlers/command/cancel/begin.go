package cancel

import (
	"context"
	"errors"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/db"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/chistyakoviv/logbot/internal/service/commands"
	"github.com/chistyakoviv/logbot/internal/service/user_settings"
)

func begin(
	ctx context.Context,
	logger *slog.Logger,
	i18n *I18n.I18n,
	commands commands.IService,
	userSettings user_settings.IService,
) handlers.Response {
	return func(b *gotgbot.Bot, ectx *ext.Context) error {
		msg := ectx.EffectiveMessage

		logger.Debug(
			"cancel current command",
			slog.Int64("chat_id", msg.Chat.Id),
			slog.String("from", msg.From.Username),
			slog.String("message", msg.Text),
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
			return err
		}

		currCommand, err := commands.FindByKey(
			ctx,
			&model.CommandKey{
				ChatId: msg.Chat.Id,
				UserId: msg.From.Id,
			},
		)
		if err != nil && !errors.Is(err, db.ErrNotFound) {
			logger.Error("error occurred while finding a command", slogger.Err(err))
			_, err := b.SendMessage(
				msg.Chat.Id,
				i18n.
					Chain().
					T(lang, "mention", I18n.WithArgs([]any{
						msg.From.Id,
						msg.From.Username,
					})).
					Append("\n").
					T(lang, "cancel_command_error").
					String(),
				&gotgbot.SendMessageOpts{
					ParseMode: "html",
				},
			)
			return err
		}
		if errors.Is(err, db.ErrNotFound) || currCommand.Stage == model.NoStage {
			_, err := b.SendMessage(
				msg.Chat.Id,
				i18n.
					Chain().
					T(lang, "mention", I18n.WithArgs([]any{
						msg.From.Id,
						msg.From.Username,
					})).
					Append("\n").
					T(lang, "cancel_no_current_command").
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
		if err != nil {
			logger.Error("error occurred while canceling command", slogger.Err(err))
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
					T(lang, "cancel_command_error").
					String(),
				&gotgbot.SendMessageOpts{
					ParseMode: "html",
				},
			)
			return err
		}

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
			T(lang, "cancel_command").
			String()
		_, err = b.SendMessage(msg.Chat.Id, message, &gotgbot.SendMessageOpts{
			ParseMode: "html",
		})
		return err
	}
}
