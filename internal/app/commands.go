package app

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command/start"
	"github.com/chistyakoviv/logbot/internal/di"
)

func BuildTgCommands(ctx context.Context, c di.Container) command.TgCommands {
	logger := resolveLogger(c)
	i18n := resolveI18n(c)
	return command.TgCommands{
		"start": start.New(logger, i18n),
	}
}
