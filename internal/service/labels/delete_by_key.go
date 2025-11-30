package labels

import (
	"context"
	"slices"
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

	for _, l := range oldLabel.Labels {
		// Skip if the label in the list of labels to delete
		if !slices.Contains(labels, l) {
			newLabel.Labels = append(newLabel.Labels, l)
		}
	}

	return s.labelsRepository.Update(ctx, newLabel)
}
