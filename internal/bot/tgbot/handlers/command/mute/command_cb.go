package mute

import (
	"context"
	"log/slog"
	"net/url"
	"strconv"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middleware"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middleware/middlewares"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/chistyakoviv/logbot/internal/service/chat_settings"
)

func muteCb(
	logger *slog.Logger,
	i18n I18n.I18nInterface,
	chatSettings chat_settings.ServiceInterface,
) middleware.TgMiddlewareHandler {
	return func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) (context.Context, error) {
		cb := ectx.CallbackQuery

		logger.Debug(
			"mute command: button clicked",
			slog.Int64("chat_id", cb.Message.GetChat().Id),
			slog.String("from", cb.From.Username),
		)

		lang, ok := ctx.Value(middlewares.LangKey).(string)
		if !ok {
			return ctx, middlewares.ErrMissingLangMiddleware
		}

		query, err := url.Parse(cb.Data)
		if err != nil {
			logger.Error("error occurred while parsing the callback data", slogger.Err(err))
			_, err := b.SendMessage(
				cb.Message.GetChat().Id,
				i18n.T(lang, "callback_data_parse_error"),
				&gotgbot.SendMessageOpts{
					ParseMode: "html",
				},
			)
			return ctx, err
		}
		queryParams := query.Query()
		rawPeriod := queryParams.Get(mutePeriodParam)
		periodIdx, err := strconv.Atoi(rawPeriod)
		if err != nil {
			logger.Error("error occurred while parsing the callback data", slogger.Err(err))
			_, err := b.SendMessage(
				cb.Message.GetChat().Id,
				i18n.T(lang, "callback_data_parse_error"),
				&gotgbot.SendMessageOpts{
					ParseMode: "html",
				},
			)
			return ctx, err
		}
		if periodIdx < 0 || periodIdx >= len(periods) {
			logger.Error("period out of range error", slog.Attr{Key: "index", Value: slog.IntValue(periodIdx)})
			_, err := b.SendMessage(
				cb.Message.GetChat().Id,
				i18n.T(lang, "callback_data_parse_error"),
				&gotgbot.SendMessageOpts{
					ParseMode: "html",
				},
			)
			return ctx, err
		}

		period := periods[periodIdx].Duration

		_, err = chatSettings.Update(ctx, cb.Message.GetChat().Id, &model.ChatSettingsInfo{
			MuteUntil: time.Now().Add(period),
		})
		if err != nil {
			logger.Error("error occurred while setting the collapse period", slogger.Err(err))
			_, err := b.SendMessage(
				cb.Message.GetChat().Id,
				i18n.T(lang, "callback_data_parse_error"),
				&gotgbot.SendMessageOpts{
					ParseMode: "html",
				},
			)
			return ctx, err
		}

		_, err = cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Text: i18n.T(
				lang,
				"mute_period_set",
				I18n.WithArgs([]any{
					period,
				}),
			),
		})
		if err != nil {
			logger.Error("failed to answer callback", slogger.Err(err))
			_, err := b.SendMessage(
				cb.Message.GetChat().Id,
				i18n.T(lang, "callback_failed_to_answer"),
				&gotgbot.SendMessageOpts{
					ParseMode: "html",
				},
			)
			return ctx, err
		}

		_, err = b.SendMessage(
			cb.Message.GetChat().Id,
			i18n.
				Chain().
				T(
					lang,
					"mention",
					I18n.WithArgs([]any{
						cb.From.Id,
						cb.From.Username,
					}),
				).
				Append("\n").
				T(
					lang,
					"mute_period_set",
					I18n.WithArgs([]any{
						period,
					}),
				).
				String(),
			&gotgbot.SendMessageOpts{
				ParseMode: "html",
			},
		)
		return ctx, err
	}
}
