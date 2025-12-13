package help

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/addlabels"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/cancel"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/collapse"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/gentoken"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/labels"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/mute"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/rmlabels"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/setlang"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/subscribe"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/subscriptions"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/unsubscribe"
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
			"help command",
			slog.Int64("chat_id", msg.Chat.Id),
			slog.String("from", msg.From.Username),
		)

		lang, ok := ctx.Value(middleware.LangKey).(string)
		if !ok {
			return ctx, middleware.ErrMissingLangMiddleware
		}

		_, err := b.SendMessage(
			msg.Chat.Id,
			i18n.
				Chain().
				T(lang, "help_title").
				Append("\n\n").
				Appendf("/%s - ", CommandName).
				T(lang, "help_description").
				Append("\n\n").
				Appendf("/%s - ", cancel.CommandName).
				T(lang, "cancel_description").
				Append("\n\n").
				Appendf("/%s - ", setlang.CommandName).
				T(lang, "setlang_description").
				Append("\n\n").
				Appendf("/%s - ", gentoken.CommandName).
				T(lang, "gentoken_description").
				Append("\n\n").
				Appendf("/%s - ", collapse.CommandName).
				T(lang, "collapse_description").
				Append("\n\n").
				Appendf("/%s - ", mute.CommandName).
				T(lang, "mute_description").
				Append("\n\n").
				Appendf("/%s - ", addlabels.CommandName).
				T(lang, "addlabels_description").
				Append("\n\n").
				Appendf("/%s - ", rmlabels.CommandName).
				T(lang, "rmlabels_description").
				Append("\n\n").
				Appendf("/%s - ", labels.CommandName).
				T(lang, "labels_description").
				Append("\n\n").
				Appendf("/%s - ", subscribe.CommandName).
				T(lang, "subscribe_description").
				Append("\n\n").
				Appendf("/%s - ", unsubscribe.CommandName).
				T(lang, "unsubscribe_description").
				Append("\n\n").
				Appendf("/%s - ", subscriptions.CommandName).
				T(lang, "subscriptions_description").
				String(),
			&gotgbot.SendMessageOpts{
				ParseMode: "html",
			},
		)
		return ctx, err
	}
}
