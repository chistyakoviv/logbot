package labels

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) FindByLabel(ctx context.Context, label string) ([]*model.Label, error) {
	return s.labelsRepository.FindByLabel(ctx, label)
}
