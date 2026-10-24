// DEPRECATE: not used

package service

import (
	"context"

	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
	"github.com/google/uuid"
)

type TenantService interface {
	// CreateTenant(ctx context.Context, params repo.CreateTenantParams) (repo.Tenant, error)
	DeleteTenantAndLease(ctx context.Context, id string) error
}

type tenantService struct {
	repo repo.Querier
}

func NewTenantService(repo repo.Querier) *tenantService {
	return &tenantService{repo}
}

// func (s *tenantService) CreateTenant(ctx context.Context, params repo.CreateTenantParams) (repo.Tenant, error) {
// 	return s.repo.CreateTenant(ctx, params)
// }

func (s *tenantService) DeleteTenantAndLease(ctx context.Context, id string) error {
	// parse uuid
	tenantUuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	tenant, err := s.repo.GetTenantById(ctx, tenantUuid)
	if err != nil {
		return err
	}

	err = s.repo.DeleteLease(ctx, tenant.LeaseID)
	if err != nil {
		return err
	}

	err = s.repo.DeleteTenant(ctx, tenantUuid)
	if err != nil {
		return err
	}

	return nil
}
