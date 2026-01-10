package labels

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) AddByKey(ctx context.Context, chatId int64, users []string, labels []string) ([]*model.Label, error) {
	result := make([]*model.Label, 0, len(users))
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		for _, user := range users {
			newLabel := &model.Label{
				ChatId:    chatId,
				Username:  user,
				Labels:    nil,
				UpdatedAt: time.Now(),
			}

			oldLabel, err := s.labelsRepository.FindByKey(ctx, &model.LabelKey{
				ChatId:   chatId,
				Username: user,
			})
			if err != nil {
				if !errors.Is(err, db.ErrNotFound) {
					return fmt.Errorf("failed to add labels: %w", err)
				}

				item, err := s.labelsRepository.Create(ctx, newLabel)
				if err != nil {
					return fmt.Errorf("failed to add labels: %w", err)
				}
				result = append(result, item)
				continue
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

			item, err := s.labelsRepository.Update(ctx, newLabel)
			if err != nil {
				return fmt.Errorf("failed to add labels: %w", err)
			}
			result = append(result, item)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}
