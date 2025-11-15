package commands

import (
	"context"
	"time"

	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) UpdateByKey(
	ctx context.Context,
	key *model.CommandKey,
	stage int,
	data map[string]interface{},
) (*model.Command, error) {
	command := &model.Command{
		UserId:    key.UserId,
		ChatId:    key.ChatId,
		Stage:     stage,
		Data:      data,
		UpdatedAt: time.Now(),
	}
	return s.commandsRepository.UpdateOrCreate(ctx, command)
}
