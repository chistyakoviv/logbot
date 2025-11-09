package app

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/commands/tgcommand"
	"github.com/chistyakoviv/logbot/internal/commands/tgcommand/start"
	"github.com/chistyakoviv/logbot/internal/di"
)

func BuildTgCommands(ctx context.Context, c di.Container) tgcommand.TgCommands {
	logger := resolveLogger(c)
	i18n := resolveI18n(c)
	return tgcommand.TgCommands{
		"start": start.New(logger, i18n),
	}
}
