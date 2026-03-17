package tenant

import "context"

type Service interface {
	ListTenants(ctx context.Context) error
}

type svc struct {
	// repository
}

func (s *svc) ListTenants(ctx context.Context) error {
	return nil
}

func NewService() Service {
	return &svc{}
}
