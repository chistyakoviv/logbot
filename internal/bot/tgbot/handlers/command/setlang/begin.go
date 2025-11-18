package setlang

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/db"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
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
			"set language command: initiate",
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
			return err
		}

		var langs []gotgbot.InlineKeyboardButton
		for _, lang := range i18n.GetLangs() {
			queryParams := url.Values{}
			queryParams.Add(langParam, lang)
			langs = append(langs, gotgbot.InlineKeyboardButton{
				Text:         lang,
				CallbackData: fmt.Sprintf("%s?%s", CommandName, queryParams.Encode()),
			})
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
				T(lang, "setlang_select_language").
				String(),
			&gotgbot.SendMessageOpts{
				ParseMode: "html",
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{
					InlineKeyboard: [][]gotgbot.InlineKeyboardButton{langs},
				},
			},
		)
		return err
	}
}
