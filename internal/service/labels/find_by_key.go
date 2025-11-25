package labels

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) FindByKey(ctx context.Context, key *model.LabelKey) (*model.Label, error) {
	return s.labelsRepository.FindByKey(ctx, key)
}
