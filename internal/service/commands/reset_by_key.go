package commands

import (
	"context"
	"errors"
	"time"

	"github.com/chistyakoviv/logbot/internal/db"
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

	_, err := s.commandsRepository.FindByKey(ctx, key)
	if err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return nil, err
		}

		return s.commandsRepository.Create(ctx, command)
	}

	return s.commandsRepository.Update(ctx, command)
}
