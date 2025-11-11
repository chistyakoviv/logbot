package handler

import (
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/i18n"
)

func NewCommandStage(logger *slog.Logger, i18n *i18n.I18n) handlers.Response {
	return commandStageHandler(logger, i18n)
}

func commandStageHandler(logger *slog.Logger, i18n *i18n.I18n) handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		_, err := b.SendMessage(ctx.EffectiveMessage.Chat.Id, "No command received", nil)
		return err
	}
}
