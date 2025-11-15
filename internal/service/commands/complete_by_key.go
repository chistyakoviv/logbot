package commands

import (
	"context"
	"time"

	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) CompleteByKey(ctx context.Context, key *model.CommandKey) (*model.Command, error) {
	command := &model.Command{
		UserId:    key.UserId,
		ChatId:    key.ChatId,
		Stage:     model.NoStage,
		Data:      nil,
		UpdatedAt: time.Now(),
	}
	return s.commandsRepository.UpdateOrCreate(ctx, command)
}
