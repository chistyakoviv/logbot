package commands

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) FindByKey(ctx context.Context, key *model.CommandKey) (*model.Command, error) {
	return s.commandsRepository.FindByKey(ctx, key)
}
