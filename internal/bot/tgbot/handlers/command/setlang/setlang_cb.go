package setlang

import (
	"context"
	"errors"
	"log/slog"
	"net/url"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares/middleware"
	"github.com/chistyakoviv/logbot/internal/db"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/i18n/language"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/chistyakoviv/logbot/internal/service/user_settings"
)

func setlangCb(
	logger *slog.Logger,
	i18n I18n.I18nInterface,
	userSettings user_settings.ServiceInterface,
) middlewares.TgMiddlewareHandler {
	return func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) (context.Context, error) {
		cb := ectx.CallbackQuery

		logger.Debug(
			"set language command: button clicked",
			slog.Int64("chat_id", cb.Message.GetChat().Id),
			slog.String("from", cb.From.Username),
		)

		isSilenced, ok := ctx.Value(middleware.SilenceKey).(bool)
		if !ok {
			return ctx, middleware.ErrMissingSilenceMiddleware
		}

		lang, currLangErr := userSettings.GetLang(ctx, cb.From.Id)
		if currLangErr != nil && !errors.Is(currLangErr, db.ErrNotFound) {
			logger.Error("error occurred while getting the user's language", slogger.Err(currLangErr))
			_, err := b.SendMessage(
				cb.Message.GetChat().Id,
				i18n.T(lang, "callback_data_parse_error"),
				&gotgbot.SendMessageOpts{
					DisableNotification: isSilenced,
					ParseMode:           "html",
				},
			)
			return ctx, err
		}

		query, err := url.Parse(cb.Data)
		if err != nil {
			logger.Error("error occurred while parsing the callback data", slogger.Err(err))
			_, err := b.SendMessage(
				cb.Message.GetChat().Id,
				i18n.T(lang, "callback_data_parse_error"),
				&gotgbot.SendMessageOpts{
					DisableNotification: isSilenced,
					ParseMode:           "html",
				},
			)
			return ctx, err
		}
		queryParams := query.Query()
		newLang := queryParams.Get(langParam)
		// Check if the selected language is the same as the current language
		if newLang == lang && !errors.Is(currLangErr, db.ErrNotFound) {
			_, err = cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
				Text: i18n.T(lang, "setlang_same_language"),
			})
			if err != nil {
				logger.Error("failed to answer callback", slogger.Err(err))
				_, err := b.SendMessage(
					cb.Message.GetChat().Id,
					i18n.T(lang, "callback_failed_to_answer"),
					&gotgbot.SendMessageOpts{
						DisableNotification: isSilenced,
						ParseMode:           "html",
					},
				)
				return ctx, err
			}
			return ctx, err
		}
		newLangCode := i18n.GetLangCode(newLang)
		// Check if the selected language is supported
		if newLangCode == language.UnknownLanguage {
			_, err = cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
				Text: i18n.T(lang, "setlang_unknown_language"),
			})
			if err != nil {
				logger.Error("failed to answer callback", slogger.Err(err))
				_, err := b.SendMessage(
					cb.Message.GetChat().Id,
					i18n.T(lang, "callback_failed_to_answer"),
					&gotgbot.SendMessageOpts{
						DisableNotification: isSilenced,
						ParseMode:           "html",
					},
				)
				return ctx, err
			}
			return ctx, err
		}
		_, err = userSettings.Update(ctx, cb.From.Id, &model.UserSettingsInfo{
			Username: cb.From.Username,
			Lang:     newLangCode,
		})
		if err != nil {
			logger.Error("error occurred while setting the user's language", slogger.Err(err))
			_, err := b.SendMessage(
				cb.Message.GetChat().Id,
				i18n.T(lang, "setlang_error"),
				&gotgbot.SendMessageOpts{
					DisableNotification: isSilenced,
					ParseMode:           "html",
				},
			)
			return ctx, err
		}

		_, err = cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Text: i18n.T(lang, "setlang_success"),
		})
		if err != nil {
			logger.Error("failed to answer callback", slogger.Err(err))
			_, err := b.SendMessage(
				cb.Message.GetChat().Id,
				i18n.T(lang, "callback_failed_to_answer"),
				&gotgbot.SendMessageOpts{
					DisableNotification: isSilenced,
					ParseMode:           "html",
				},
			)
			return ctx, err
		}
		return ctx, err
	}
}
