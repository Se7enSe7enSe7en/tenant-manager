// Service for tenant registration including creating the lease

package service

import (
	"context"
	"time"

	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type RegistrationService interface {
	RegisterTenantWithLease(ctx context.Context, params RegisterTenantWithLeaseParams) (repo.Tenant, repo.Lease, error)
}

type registrationService struct {
	db      *pgxpool.Pool
	queries *repo.Queries
}

func NewRegistrationService(db *pgxpool.Pool, queries *repo.Queries) *registrationService {
	return new(registrationService{
		db:      db,
		queries: queries,
	})
}

// TODO: when moving to go 1.27, new feature with embedded structs here
type RegisterTenantWithLeaseParams struct {
	// create tenant params
	Email       string
	Name        string
	PhoneNumber string

	// create lease params
	PropertyID      uuid.UUID
	ExpectedRentDay int
	StartDate       *time.Time
	IsMonthAdvance  bool
	DepositAmount   decimal.Decimal
}

func (s *registrationService) RegisterTenantWithLease(ctx context.Context, params RegisterTenantWithLeaseParams) (repo.Tenant, repo.Lease, error) {
	// init transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Tenant{}, repo.Lease{}, err
	}
	defer tx.Rollback(ctx)
	qtx := s.queries.WithTx(tx)

	// create tenant
	tenant, err := qtx.CreateTenant(ctx, repo.CreateTenantParams{
		Email:       params.Email,
		Name:        params.Name,
		PhoneNumber: params.PhoneNumber,
	})
	if err != nil {
		return repo.Tenant{}, repo.Lease{}, err
	}

	// create lease
	//// remove hours, minutes, seconds, nanoseconds from start date
	startDateNoHMS := time.Date(
		params.StartDate.Year(),
		params.StartDate.Month(),
		params.StartDate.Day(),
		0, 0, 0, 0, // H:M:S:Ns
		params.StartDate.Location(),
	)

	//// get expiry date
	expiryDate := utils.ComputeNextExpiryDate(*params.StartDate, params.ExpectedRentDay)

	lease, err := qtx.CreateLease(ctx, repo.CreateLeaseParams{
		PropertyID:      params.PropertyID,
		TenantID:        tenant.ID,
		ExpectedRentDay: int16(params.ExpectedRentDay),
		StartDate:       startDateNoHMS,
		ExpiryDate:      &expiryDate,
		IsMonthAdvance:  params.IsMonthAdvance,
		DepositAmount:   params.DepositAmount,
	})
	if err != nil {
		return repo.Tenant{}, repo.Lease{}, err
	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return repo.Tenant{}, repo.Lease{}, err
	}

	return tenant, lease, nil
}
