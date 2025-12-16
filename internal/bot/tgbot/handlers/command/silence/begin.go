package silence

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares/middleware"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/utils"
)

const columns = 2

func begin(
	logger *slog.Logger,
	i18n I18n.I18nInterface,
) middlewares.TgMiddlewareHandler {
	return func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) error {
		msg := ectx.EffectiveMessage

		logger.Debug(
			"silence command: initiate",
			slog.Int64("chat_id", msg.Chat.Id),
			slog.String("from", msg.From.Username),
		)

		lang, ok := ctx.Value(middleware.LangKey).(string)
		if !ok {
			return middleware.ErrMissingLangMiddleware
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
				// For the silence command notifications are always disabled
				DisableNotification: true,
				ParseMode:           "html",
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{
					InlineKeyboard: utils.Chunk(buttons, columns),
				},
			},
		)
		return err
	}
}
