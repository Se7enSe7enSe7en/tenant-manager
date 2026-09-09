package service

import (
	"context"
	"time"

	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type LeaseService interface {
	// CreateLease(ctx context.Context, params CreateLeaseParams) (repo.Lease, error)
	ListLease(ctx context.Context, userID uuid.UUID) ([]repo.ListLeaseRow, error)
	GetLease(ctx context.Context, leaseID uuid.UUID) (repo.GetLeaseByIdRow, error)
}

type leaseService struct {
	queries *repo.Queries
}

func NewLeaseService(queries *repo.Queries) *leaseService {
	return new(leaseService{queries})
}

type CreateLeaseParams struct {
	PropertyID      uuid.UUID
	TenantID        uuid.UUID
	ExpectedRentDay int
	StartDate       *time.Time
	IsMonthAdvance  bool
	DepositAmount   decimal.Decimal
}

// // DEPRECATED: no longer used
// // 1. "ExpectedRentDay = 31" case, expiry date should be the last day of next month
// // 2. ExpectedRentDay = 1, StartDate = Aug 29 2026, expected output should be ExpiryDate = Oct 1 2026 (not Sep 1 2026)
// // 2. ExpectedRentDay = 28, StartDate = Aug 29 2026, expected output should be ExpiryDate = Sep 28 2026 (not Oct 28 2026)
// func (s *leaseService) CreateLease(ctx context.Context, params CreateLeaseParams) (repo.Lease, error) {
// 	// remove hours, minutes, seconds, nanoseconds from start date
// 	startDateNoHMS := time.Date(
// 		params.StartDate.Year(),
// 		params.StartDate.Month(),
// 		params.StartDate.Day(),
// 		0, 0, 0, 0, // H:M:S:Ns
// 		params.StartDate.Location(),
// 	)

// 	// get expiry date
// 	expiryDate := utils.ComputeNextExpiryDate(startDateNoHMS, params.ExpectedRentDay)

// 	return s.repo.CreateLease(ctx, repo.CreateLeaseParams{
// 		PropertyID:      params.PropertyID,
// 		TenantID:        params.TenantID,
// 		ExpectedRentDay: int16(params.ExpectedRentDay),
// 		StartDate:       startDateNoHMS,
// 		ExpiryDate:      &expiryDate,
// 		IsMonthAdvance:  params.IsMonthAdvance,
// 		DepositAmount:   params.DepositAmount,
// 	})
// }

func (s *leaseService) ListLease(ctx context.Context, userID uuid.UUID) ([]repo.ListLeaseRow, error) {
	return s.queries.ListLease(ctx, userID)
}

func (s *leaseService) GetLease(ctx context.Context, leaseID uuid.UUID) (repo.GetLeaseByIdRow, error) {
	return s.queries.GetLeaseById(ctx, leaseID)
}
