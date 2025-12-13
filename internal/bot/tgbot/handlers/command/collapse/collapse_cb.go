package collapse

import (
	"context"
	"log/slog"
	"net/url"
	"strconv"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares/middleware"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
	"github.com/chistyakoviv/logbot/internal/service/chat_settings"
)

func collapseCb(
	logger *slog.Logger,
	i18n I18n.I18nInterface,
	chatSettings chat_settings.ServiceInterface,
) middlewares.TgMiddlewareHandler {
	return func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) (context.Context, error) {
		cb := ectx.CallbackQuery

		logger.Debug(
			"collapse command: button clicked",
			slog.Int64("chat_id", cb.Message.GetChat().Id),
			slog.String("from", cb.From.Username),
		)

		lang, ok := ctx.Value(middleware.LangKey).(string)
		if !ok {
			return ctx, middleware.ErrMissingLangMiddleware
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
		rawPeriod := queryParams.Get(collapsePeriodParam)
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
		var periodArg any
		if period == 0 {
			periodArg = "none"
		} else {
			periodArg = period
		}

		_, err = chatSettings.UpdateCollapsePeriod(ctx, cb.Message.GetChat().Id, period)
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
				"collapse_period_set",
				I18n.WithArgs([]any{
					periodArg,
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
					"collapse_period_set",
					I18n.WithArgs([]any{
						periodArg,
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
