// DEPRECATE: not used

package service

// import (
// 	"context"

// 	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
// )

// type TenantService interface {
// 	CreateTenant(ctx context.Context, params repo.CreateTenantParams) (repo.Tenant, error)
// }

// type tenantService struct {
// 	repo repo.Querier
// }

// func NewTenantService(repo repo.Querier) *tenantService {
// 	return &tenantService{repo}
// }

// func (s *tenantService) CreateTenant(ctx context.Context, params repo.CreateTenantParams) (repo.Tenant, error) {
// 	return s.repo.CreateTenant(ctx, params)
// }
