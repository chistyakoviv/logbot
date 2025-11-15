package commands

import (
	"context"
	"time"

	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) ResetByKey(
	ctx context.Context,
	key *model.CommandKey,
	name string,
	data map[string]interface{},
) (*model.Command, error) {
	command := &model.Command{
		Name:      name,
		UserId:    key.UserId,
		ChatId:    key.ChatId,
		Stage:     0,
		Data:      data,
		UpdatedAt: time.Now(),
	}
	return s.commandsRepository.UpdateOrCreate(ctx, command)
}
