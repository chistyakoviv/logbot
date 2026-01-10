package labels

import (
	"context"
	"time"

	"github.com/chistyakoviv/logbot/internal/model"
)

func (s *service) DeleteByKey(ctx context.Context, chatId int64, users []string, labels []string) ([]*model.Label, error) {
	result := make([]*model.Label, 0, len(users))
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		for _, user := range users {
			oldLabel, err := s.labelsRepository.FindByKey(ctx, &model.LabelKey{
				ChatId:   chatId,
				Username: user,
			})
			if err != nil {
				return err
			}

			newLabel := &model.Label{
				ChatId:    chatId,
				Username:  user,
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
			item, err := s.labelsRepository.Update(ctx, newLabel)
			if err != nil {
				return err
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
