package labels

import (
	"context"
	"errors"
	"maps"
	"slices"
	"time"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) AddByKey(ctx context.Context, key *model.LabelKey, labels []string) (*model.Label, error) {
	newLabel := &model.Label{
		ChatId:    key.ChatId,
		Username:  key.Username,
		Labels:    nil,
		UpdatedAt: time.Now(),
	}

	oldLabel, err := s.labelsRepository.FindByKey(ctx, key)
	if err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return nil, err
		}

		return s.labelsRepository.Create(ctx, newLabel)
	}

	labelsSet := make(map[string]bool, len(oldLabel.Labels)+len(labels))
	for _, lable := range oldLabel.Labels {
		labelsSet[lable] = true
	}
	for _, label := range labels {
		labelsSet[label] = true
	}

	newLabel.Labels = make([]string, 0, len(labelsSet))
	newLabel.Labels = slices.AppendSeq(newLabel.Labels, maps.Keys(labelsSet))

	return s.labelsRepository.Update(ctx, newLabel)
}
