package logs

import (
	"context"

	"github.com/google/uuid"
)

func (s *service) DeleteByToken(ctx context.Context, token uuid.UUID) error {
	return s.logsRepository.DeleteByToken(ctx, token)
}
