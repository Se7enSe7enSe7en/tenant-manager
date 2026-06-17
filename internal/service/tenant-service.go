package service

import (
	"context"

	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/utils"
)

type TenantService interface {
	ListTenants(ctx context.Context) ([]repo.Tenant, error)
	CreateTenant(ctx context.Context, params repo.CreateTenantParams) (repo.Tenant, error)
	ListTenantsWithProperty(ctx context.Context, userId string) ([]repo.ListTenantsWithPropertyRow, error)
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

func (s *tenantService) CreateTenant(ctx context.Context, params repo.CreateTenantParams) (repo.Tenant, error) {
	return s.repo.CreateTenant(ctx, params)
}

func (s *tenantService) ListTenantsWithProperty(ctx context.Context, userId string) ([]repo.ListTenantsWithPropertyRow, error) {
	userIdPgtypeUuid, err := utils.StringToPgtypeUuid(userId)
	if err != nil {
		return nil, err
	}

	return s.repo.ListTenantsWithProperty(ctx, userIdPgtypeUuid)
}
