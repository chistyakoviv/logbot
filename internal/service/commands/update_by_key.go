package commands

import (
	"context"
	"errors"
	"time"

	"github.com/chistyakoviv/logbot/internal/db"
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

	_, err := s.commandsRepository.FindByKey(ctx, key)
	if err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return nil, err
		}

		return s.commandsRepository.Create(ctx, command)
	}

	return s.commandsRepository.Update(ctx, command)
}
