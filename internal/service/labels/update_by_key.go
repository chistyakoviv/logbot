package labels

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) UpdateByKey(ctx context.Context, key *model.LabelKey, labels []string) (*model.Label, error) {
	newLabel := &model.Label{
		ChatId:    key.ChatId,
		UserId:    key.UserId,
		Labels:    labels,
		UpdatedAt: time.Now(),
	}

	oldLabel, err := s.labelsRepository.FindByKey(ctx, key)
	if err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return nil, err
		}

		return s.labelsRepository.Create(ctx, newLabel)
	}

	newLabel.Labels = make([]string, 0, len(oldLabel.Labels)+len(labels))
	newLabel.Labels = append(newLabel.Labels, oldLabel.Labels...)
	for _, l := range labels {
		if !slices.Contains(newLabel.Labels, l) {
			newLabel.Labels = append(newLabel.Labels, l)
		}
	}

	return s.labelsRepository.Update(ctx, newLabel)
}
