package addlabels

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares/middleware"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/chistyakoviv/logbot/internal/service/commands"
	"github.com/chistyakoviv/logbot/internal/service/labels"
)

func applyLabels(
	logger *slog.Logger,
	i18n I18n.I18nInterface,
	labelsService labels.ServiceInterface,
	commands commands.ServiceInterface,
) middlewares.TgMiddlewareHandler {
	return func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) error {
		msg := ectx.EffectiveMessage

		logger.Debug(
			"add label command: apply labels to users",
			slog.Int64("chat_id", msg.Chat.Id),
			slog.String("user", msg.From.Username),
			slog.String("label", msg.Text),
		)

		lang, ok := ctx.Value(middleware.LangKey).(string)
		if !ok {
			return middleware.ErrMissingLangMiddleware
		}

		isSilenced, ok := ctx.Value(middleware.SilenceKey).(bool)
		if !ok {
			return middleware.ErrMissingSilenceMiddleware
		}

		labels := make([]string, 0)
		for _, label := range strings.Split(msg.Text, ",") {
			trimmedLabel := strings.TrimSpace(label)
			if len(trimmedLabel) > 0 && !strings.Contains(trimmedLabel, " ") {
				labels = append(labels, trimmedLabel)
			}
		}

		if len(labels) == 0 {
			_, _ = b.SendMessage(
				msg.Chat.Id,
				i18n.T(lang, "addlabels_no_labels_error"),
				&gotgbot.SendMessageOpts{
					DisableNotification: isSilenced,
					ParseMode:           "html",
				},
			)
			return errors.New("no labels provided")
		}

		var err error
		_, err = commands.CompleteByKey(
			ctx,
			&model.CommandKey{
				ChatId: msg.Chat.Id,
				UserId: msg.From.Id,
			},
		)
		if err != nil {
			logger.Error("error occurred while adding labels: failed to complete command", slogger.Err(err))
			_, _ = b.SendMessage(
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
					T(lang, "addlabels_error").
					String(),
				&gotgbot.SendMessageOpts{
					DisableNotification: isSilenced,
					ParseMode:           "html",
				},
			)
			return err
		}

		cmd, err := commands.FindByKey(
			ctx,
			&model.CommandKey{
				ChatId: msg.Chat.Id,
				UserId: msg.From.Id,
			},
		)
		if err != nil {
			logger.Error("error occurred while adding labels: failed to fetch command data", slogger.Err(err))
			_, _ = b.SendMessage(
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
					T(lang, "addlabels_error").
					String(),
				&gotgbot.SendMessageOpts{
					DisableNotification: isSilenced,
					ParseMode:           "html",
				},
			)
			return err
		}
		rawUsers, ok := cmd.Data["users"].([]interface{})
		if !ok {
			logger.Error("error occurred while adding labels: failed to unmarshal command data, missing users")
			_, _ = b.SendMessage(
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
					T(lang, "addlabels_error").
					String(),
				&gotgbot.SendMessageOpts{
					DisableNotification: isSilenced,
					ParseMode:           "html",
				},
			)
			return err
		}

		users := make([]string, len(rawUsers))
		for i, user := range rawUsers {
			users[i] = user.(string)
		}

		_, err = labelsService.AddByKey(
			ctx,
			msg.Chat.Id,
			users,
			labels,
		)

		if err != nil {
			_, _ = b.SendMessage(
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
					T(lang, "addlabels_failed_apply_labels_error").
					String(),
				&gotgbot.SendMessageOpts{
					DisableNotification: isSilenced,
					ParseMode:           "html",
				},
			)
			return err
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
				T(lang, "addlabels_success").
				String(),
			&gotgbot.SendMessageOpts{
				DisableNotification: isSilenced,
				ParseMode:           "html",
			},
		)
		if err != nil {
			logger.Error("error occurred while adding labels", slogger.Err(err))
		}

		return err
	}
}
