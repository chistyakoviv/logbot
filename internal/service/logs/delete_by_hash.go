package logs

import "context"

func (s *service) DeleteByHash(ctx context.Context, hash string) error {
	return s.logsRepository.DeleteByHash(ctx, hash)
}
