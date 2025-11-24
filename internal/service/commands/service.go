package commands

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/chistyakoviv/logbot/internal/repository/commands"
)

type ServiceInterface interface {
	UpdateByKey(ctx context.Context, key *model.CommandKey, stage int, data map[string]interface{}) (*model.Command, error)
	FindByKey(ctx context.Context, key *model.CommandKey) (*model.Command, error)
	ResetByKey(ctx context.Context, key *model.CommandKey, name string, data map[string]interface{}) (*model.Command, error)
	CompleteByKey(ctx context.Context, key *model.CommandKey) (*model.Command, error)
}

type service struct {
	commandsRepository commands.RepositoryInterface
	txManager          db.TxManager
}

func NewService(commandsRepository commands.RepositoryInterface, txManager db.TxManager) ServiceInterface {
	return &service{
		commandsRepository: commandsRepository,
		txManager:          txManager,
	}
}
