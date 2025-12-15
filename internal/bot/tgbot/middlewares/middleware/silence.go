package middleware

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middlewares"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/service/chat_settings"
)

const (
	SilenceKey key = "silence"
)

var (
	ErrMissingSilenceMiddleware error = errors.New("missing silence middleware")
)

func Silence(
	logger *slog.Logger,
	i18n I18n.I18nInterface,
	chatSettings chat_settings.ServiceInterface,
) middlewares.TgMiddlewareHandler {
	return func(ctx context.Context, b *gotgbot.Bot, ectx *ext.Context) (context.Context, error) {
		msg := ectx.EffectiveMessage

		logger.Debug(
			"silence middleware",
			slog.Int64("chat_id", msg.Chat.Id),
			slog.String("from", msg.From.Username),
		)

		settings, err := chatSettings.FindOrDefaults(ctx, msg.Chat.Id)
		if err != nil {
			return ctx, err
		}

		isSilenced, _ := settings.IsSilenced(time.Now().UTC())

		ctx = context.WithValue(ctx, SilenceKey, isSilenced)
		return ctx, nil
	}
}
