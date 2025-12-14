package handler

import (
	"context"
	"errors"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/i18n"
	"github.com/chistyakoviv/logbot/internal/lib/slogger"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/chistyakoviv/logbot/internal/service/commands"
)

func NewCommandStage(
	ctx context.Context,
	logger *slog.Logger,
	i18n i18n.I18nInterface,
	commands commands.ServiceInterface,
	tgCommands command.TgCommands,
) handlers.Response {
	return commandStageHandler(ctx, logger, commands, tgCommands)
}

func commandStageHandler(ctx context.Context,
	logger *slog.Logger,
	commands commands.ServiceInterface,
	tgCommands command.TgCommands,
) handlers.Response {
	return func(b *gotgbot.Bot, ectx *ext.Context) error {
		msg := ectx.EffectiveMessage

		logger.Debug(
			"try to find ongoing command",
			slog.String("chat", msg.Chat.Title),
			slog.String("from", msg.From.Username),
		)

		command, err := commands.FindByKey(ctx, &model.CommandKey{
			ChatId: msg.Chat.Id,
			UserId: msg.From.Id,
		})
		if errors.Is(err, db.ErrNotFound) || (command != nil && command.Stage == model.NoStage) {
			logger.Debug("no ongoing command")
			return nil
		}
		if err != nil {
			logger.Error("error occurred while trying to process a command stage", slogger.Err(err))
			return err
		}

		tgCommand, ok := tgCommands[command.Name]
		if !ok {
			logger.Error("command not found", slog.String("name", command.Name))
			return nil
		}

		stages := tgCommand.GetStages()

		if command.Stage < 0 || command.Stage >= len(stages) {
			logger.Error("command stage is out of range", slog.Int("stage", command.Stage))
			return nil
		}

		return stages[command.Stage](b, ectx)
	}
}
