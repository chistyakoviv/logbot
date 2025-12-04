package logs

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/loghasher"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/chistyakoviv/logbot/internal/repository/logs"
	"github.com/google/uuid"
)

type ServiceInterface interface {
	Create(ctx context.Context, log *model.LogInfo) (*model.Log, error)
	FindAllByToken(ctx context.Context, token uuid.UUID) ([]*model.Log, error)
	Delete(ctx context.Context, id int) error
	DeleteByToken(ctx context.Context, token uuid.UUID) error
	DeleteByHash(ctx context.Context, hash string) error
}

type service struct {
	logHasher      loghasher.HasherInterface
	logsRepository logs.RepositoryInterface
	txManager      db.TxManager
}

func NewService(
	logHasher loghasher.HasherInterface,
	logsRepository logs.RepositoryInterface,
	txManager db.TxManager,
) ServiceInterface {
	return &service{
		logHasher:      loghasher.NewHasher(),
		logsRepository: logsRepository,
		txManager:      txManager,
	}
}
