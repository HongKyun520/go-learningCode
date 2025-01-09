package memory

import (
	"context"
	"fmt"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Send(ctx context.Context, tql string, args []string, phone ...string) error {
	fmt.Println(args)
	return nil
}
