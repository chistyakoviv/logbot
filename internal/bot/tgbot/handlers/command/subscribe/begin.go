package subscribe

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/chistyakoviv/logbot/internal/service/commands"
)

func begin(
	ctx context.Context,
	logger *slog.Logger,
	i18n *I18n.I18n,
	commands commands.IService,
) handlers.Response {
	return func(b *gotgbot.Bot, ectx *ext.Context) error {
		msg := ectx.EffectiveMessage

		logger.Debug(
			"subscribe command: initiate",
			slog.Int64("chat_id", msg.Chat.Id),
			slog.String("from", msg.From.Username),
		)

		_, err := commands.ResetByKey(
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
						"en",
						"mention",
						I18n.WithArgs([]any{
							msg.From.Id,
							msg.From.Username,
						}),
					).
					Append("\n").
					T("en", "subscribe_error").
					String(),
				&gotgbot.SendMessageOpts{
					ParseMode: "html",
				},
			)
			return err
		}

		// TODO: implement user settings service to obtain current language
		_, err = b.SendMessage(
			msg.Chat.Id,
			i18n.
				Chain().
				T(
					"en",
					"mention",
					I18n.WithArgs([]any{
						msg.From.Id,
						msg.From.Username,
					}),
				).
				Append("\n").
				T("en", "subscribe_begin").
				String(),
			&gotgbot.SendMessageOpts{
				ParseMode: "html",
			},
		)
		return err
	}
}
