package mute

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
	return func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) (context.Context, error) {
		msg := ectx.EffectiveMessage

		logger.Debug(
			"mute command: initiate",
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

		var buttons []gotgbot.InlineKeyboardButton
		for idx, period := range periods {
			queryParams := url.Values{}
			queryParams.Add(mutePeriodParam, strconv.Itoa(idx))
			buttons = append(buttons, gotgbot.InlineKeyboardButton{
				Text:         period.Label,
				CallbackData: fmt.Sprintf("%s?%s", muteCbName, queryParams.Encode()),
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
				T(lang, "mute_select_period").
				String(),
			&gotgbot.SendMessageOpts{
				DisableNotification: isSilenced,
				ParseMode:           "html",
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{
					InlineKeyboard: utils.Chunk(buttons, columns),
				},
			},
		)
		return ctx, err
	}
}
