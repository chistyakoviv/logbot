package commands

import (
	"context"
	"errors"
	"time"

	"github.com/chistyakoviv/logbot/internal/db"
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

	_, err := s.commandsRepository.FindByKey(ctx, key)
	if err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return nil, err
		}

		return s.commandsRepository.Create(ctx, command)
	}

	return s.commandsRepository.Update(ctx, command)
}
