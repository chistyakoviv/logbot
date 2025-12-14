package addlabels

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares/middleware"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/chistyakoviv/logbot/internal/service/commands"
)

func stage0(
	logger *slog.Logger,
	i18n I18n.I18nInterface,
	commands commands.ServiceInterface,
) middlewares.TgMiddlewareHandler {
	return func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) (context.Context, error) {
		msg := ectx.EffectiveMessage

		logger.Debug(
			"add label command: retrieve users to assign labels to",
			slog.Int64("chat_id", msg.Chat.Id),
			slog.String("user", msg.From.Username),
			slog.String("labels", msg.Text),
		)

		lang, ok := ctx.Value(middleware.LangKey).(string)
		if !ok {
			return ctx, middleware.ErrMissingLangMiddleware
		}

		isSilenced, ok := ctx.Value(middleware.SilenceKey).(bool)
		if !ok {
			return ctx, middleware.ErrMissingSilenceMiddleware
		}

		userSet := make(map[string]bool, 0)
		if msg.Entities != nil {
			for _, entity := range msg.Entities {
				logger.Info("entity", slog.Any("entiry", entity))
				if entity.Type == "text_mention" {
					userSet[strconv.FormatInt(entity.User.Id, 10)] = true
				}
				if entity.Type == "mention" {
					mention := msg.Text[entity.Offset+1 : entity.Offset+entity.Length]
					userSet[mention] = true
				}
			}
		}

		if len(userSet) == 0 {
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
					T(lang, "addlabels_no_mentions_error").
					String(),
				&gotgbot.SendMessageOpts{
					DisableNotification: isSilenced,
					ParseMode:           "html",
				},
			)
			return ctx, err
		}

		users := make([]string, 0, len(userSet))
		for k := range userSet {
			users = append(users, k)
		}

		var err error
		_, err = commands.UpdateByKey(
			ctx,
			&model.CommandKey{
				ChatId: msg.Chat.Id,
				UserId: msg.From.Id,
			},
			stageLabels,
			map[string]interface{}{
				"users": users,
			},
		)
		if err != nil {
			logger.Error("error occurred while saving mentioned users", slogger.Err(err))
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
					T(lang, "addlabels_save_mentions_error").
					String(),
				&gotgbot.SendMessageOpts{
					DisableNotification: isSilenced,
					ParseMode:           "html",
				},
			)
			return ctx, err
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
				T(lang, "addlabels_enter_labels").
				String(),
			&gotgbot.SendMessageOpts{
				DisableNotification: isSilenced,
				ParseMode:           "html",
			},
		)
		if err != nil {
			return ctx, err
		}

		return ctx, nil
	}
}
