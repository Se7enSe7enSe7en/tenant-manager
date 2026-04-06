package service

import (
	"context"

	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
)

type TenantService interface {
	ListTenants(ctx context.Context) ([]repo.Tenant, error)
}

type tenantService struct {
	repo repo.Querier
}

func NewTenantService(repo repo.Querier) *tenantService {
	return &tenantService{repo}
}

func (s *tenantService) ListTenants(ctx context.Context) ([]repo.Tenant, error) {
	return s.repo.ListTenants(ctx)
}
