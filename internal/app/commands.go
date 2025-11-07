package app

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/commands"
	"github.com/chistyakoviv/logbot/internal/commands/tgcommand"
	"github.com/chistyakoviv/logbot/internal/di"
)

func BuildTgCommands(ctx context.Context, c di.Container) []*commands.TgCommand {
	logger := resolveLogger(c)
	return []*commands.TgCommand{
		commands.NewTgCommand("start", tgcommand.Start, logger),
	}
}
