package labels

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/model"
	"github.com/chistyakoviv/logbot/internal/repository/labels"
)

type ServiceInterface interface {
	FindByLabel(ctx context.Context, label string) (*model.Label, error)
	FindByKey(ctx context.Context, key *model.LabelKey) (*model.Label, error)
	UpdateByKey(ctx context.Context, key *model.LabelKey, labels []string) (*model.Label, error)
}

type service struct {
	labelsRepository labels.RepositoryInterface
}

func NewService(labelsRepository labels.RepositoryInterface) ServiceInterface {
	return &service{
		labelsRepository: labelsRepository,
	}
}
