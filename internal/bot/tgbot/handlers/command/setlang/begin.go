package setlang

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares/middleware"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
)

func begin(
	logger *slog.Logger,
	i18n I18n.I18nInterface,
) middlewares.TgMiddlewareHandler {
	return func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) (context.Context, error) {
		msg := ectx.EffectiveMessage

		logger.Debug(
			"set language command: initiate",
			slog.Int64("chat_id", msg.Chat.Id),
			slog.String("from", msg.From.Username),
		)

		lang, ok := ctx.Value(middleware.LangKey).(string)
		if !ok {
			return ctx, middleware.ErrMissingLangMiddleware
		}

		var langs []gotgbot.InlineKeyboardButton
		for _, lang := range i18n.GetLangs() {
			queryParams := url.Values{}
			queryParams.Add(langParam, lang)
			langs = append(langs, gotgbot.InlineKeyboardButton{
				Text:         lang,
				CallbackData: fmt.Sprintf("%s?%s", SetLangCbName, queryParams.Encode()),
			})
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
				T(lang, "setlang_select_language").
				String(),
			&gotgbot.SendMessageOpts{
				ParseMode: "html",
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{
					InlineKeyboard: [][]gotgbot.InlineKeyboardButton{langs},
				},
			},
		)
		return ctx, err
	}
}
