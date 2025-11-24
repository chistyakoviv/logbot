package cancel

import (
	"context"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/middleware"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/service/commands"
)

const CommandName = "cancel"

func New(
	ctx context.Context,
	mw middleware.TgMiddlewareInterface,
	logger *slog.Logger,
	i18n *I18n.I18n,
	commands commands.IService,
) *command.TgCommand {
	return &command.TgCommand{
		Handler: mw.Pipe(begin(ctx, logger, i18n, commands)).Handler(ctx),
		Stages:  []handlers.Response{},
	}
}
