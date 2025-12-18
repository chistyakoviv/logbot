package labels

import (
	"context"
	"time"

	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) DeleteByKey(ctx context.Context, key *model.LabelKey, labels []string) (*model.Label, error) {
	oldLabel, err := s.labelsRepository.FindByKey(ctx, key)
	if err != nil {
		return nil, err
	}

	newLabel := &model.Label{
		ChatId:    key.ChatId,
		Username:  key.Username,
		Labels:    make([]string, 0, len(oldLabel.Labels)),
		UpdatedAt: time.Now(),
	}

	labelsToDelete := make(map[string]bool, len(labels))
	for _, label := range labels {
		labelsToDelete[label] = true
	}

	for _, label := range oldLabel.Labels {
		if !labelsToDelete[label] {
			newLabel.Labels = append(newLabel.Labels, label)
		}
	}

	return s.labelsRepository.Update(ctx, newLabel)
}
