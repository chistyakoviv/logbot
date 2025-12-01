package labels

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/chistyakoviv/logbot/internal/repository/labels"
)

type ServiceInterface interface {
	FindByLabel(ctx context.Context, label string) ([]*model.Label, error)
	FindByKey(ctx context.Context, key *model.LabelKey) (*model.Label, error)
	FindAllByChat(ctx context.Context, chatId int64) ([]*model.Label, error)
	AddByKey(ctx context.Context, key *model.LabelKey, labels []string) (*model.Label, error)
	DeleteByKey(ctx context.Context, key *model.LabelKey, labels []string) (*model.Label, error)
}

type service struct {
	labelsRepository labels.RepositoryInterface
	txManager        db.TxManager
}

func NewService(labelsRepository labels.RepositoryInterface, txManager db.TxManager) ServiceInterface {
	return &service{
		labelsRepository: labelsRepository,
		txManager:        txManager,
	}
}
