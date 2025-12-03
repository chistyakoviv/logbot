package silence

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middleware"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middleware/middlewares"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/utils"
)

func begin(
	logger *slog.Logger,
	i18n I18n.I18nInterface,
) middleware.TgMiddlewareHandler {
	return func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) (context.Context, error) {
		msg := ectx.EffectiveMessage

		logger.Debug(
			"silence command: initiate",
			slog.Int64("chat_id", msg.Chat.Id),
			slog.String("from", msg.From.Username),
		)

		lang, ok := ctx.Value(middlewares.LangKey).(string)
		if !ok {
			return ctx, middlewares.ErrMissingLangMiddleware
		}

		var buttons []gotgbot.InlineKeyboardButton
		for idx, period := range periods {
			queryParams := url.Values{}
			queryParams.Add(silencePeriodParam, strconv.Itoa(idx))
			buttons = append(buttons, gotgbot.InlineKeyboardButton{
				Text:         period.Label,
				CallbackData: fmt.Sprintf("%s?%s", silenceCbName, queryParams.Encode()),
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
				T(lang, "silence_select_period").
				String(),
			&gotgbot.SendMessageOpts{
				ParseMode: "html",
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{
					InlineKeyboard: utils.Chunk(buttons, 2),
				},
			},
		)
		return ctx, err
	}
}
