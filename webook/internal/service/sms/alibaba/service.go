package alibaba

import "context"

type Service struct {
}

func (s *Service) Send(ctx context.Context, tql string, args []string, phone ...string) error {
	return nil
}
