package service

import (
	"context"

	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
)

type TenantService interface {
	ListTenants(ctx context.Context) ([]repo.Tenant, error)
}

type service struct {
	repo repo.Querier
}

func NewTenantService(repo repo.Querier) TenantService {
	return &service{repo: repo}
}

func (s *service) ListTenants(ctx context.Context) ([]repo.Tenant, error) {
	return s.repo.ListTenants(ctx)
}
