package logs

import "context"

func (s *service) Delete(ctx context.Context, id int) error {
	return s.logsRepository.Delete(ctx, id)
}
