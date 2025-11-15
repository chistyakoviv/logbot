package subscribe

import (
	"context"
	"log/slog"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/chistyakoviv/logbot/internal/service/commands"
	"github.com/chistyakoviv/logbot/internal/service/subscriptions"
)

func stage0(
	ctx context.Context,
	logger *slog.Logger,
	i18n *i18n.I18n,
	subscriptions subscriptions.IService,
	commands commands.IService,
) handlers.Response {
	return func(b *gotgbot.Bot, ectx *ext.Context) error {
		msg := ectx.EffectiveMessage
		token := strings.Trim(msg.Text, " ")

		logger.Debug(
			"subscribe command: retrieve token",
			slog.Int64("chat_id", msg.Chat.Id),
			slog.String("user", msg.From.Username),
			slog.String("token", token),
		)

		if token == "" {
			_, err := b.SendMessage(msg.Chat.Id, i18n.T("en", "subscribe_empty_token"), &gotgbot.SendMessageOpts{
				ParseMode: "html",
			})
			return err
		}

		_, err := commands.CompleteByKey(
			ctx,
			&model.CommandKey{
				ChatId: msg.Chat.Id,
				UserId: msg.From.Id,
			},
		)
		if err != nil {
			logger.Error("error occurred while subscribing", slogger.Err(err))
			_, err = b.SendMessage(msg.Chat.Id, i18n.T("en", "subscribe_error"), &gotgbot.SendMessageOpts{
				ParseMode: "html",
			})
			return err
		}

		_, err = subscriptions.Subscribe(ctx, &model.SubscriptionInfo{
			ChatId: msg.Chat.Id,
			Token:  token,
		})
		if err != nil {
			logger.Error("error occurred while subscribing", slogger.Err(err))
			_, err = b.SendMessage(msg.Chat.Id, i18n.T("en", "subscribe_error"), &gotgbot.SendMessageOpts{
				ParseMode: "html",
			})
			return err
		}

		_, err = b.SendMessage(msg.Chat.Id, i18n.T("en", "subscribe_complete", token), &gotgbot.SendMessageOpts{
			ParseMode: "html",
		})
		return err
	}
}
