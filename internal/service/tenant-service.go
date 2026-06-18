package service

import (
	"context"

	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
	"github.com/google/uuid"
)

type TenantService interface {
	CreateTenant(ctx context.Context, params repo.CreateTenantParams) (repo.Tenant, error)
	ListTenantsWithProperty(ctx context.Context, userId uuid.UUID) ([]repo.ListTenantsWithPropertyRow, error)
}

type tenantService struct {
	repo repo.Querier
}

func NewTenantService(repo repo.Querier) *tenantService {
	return &tenantService{repo}
}

func (s *tenantService) CreateTenant(ctx context.Context, params repo.CreateTenantParams) (repo.Tenant, error) {
	return s.repo.CreateTenant(ctx, params)
}

func (s *tenantService) ListTenantsWithProperty(ctx context.Context, userId uuid.UUID) ([]repo.ListTenantsWithPropertyRow, error) {
	return s.repo.ListTenantsWithProperty(ctx, userId)
}
