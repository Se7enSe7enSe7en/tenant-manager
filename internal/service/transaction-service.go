// Service for user (lessor) and tenant transactions eg. Collect rent, Reduce deposit

package service

import (
	"context"

	"github.com/Se7enSe7enSe7en/tenant-manager/internal/constants"
	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/utils"
	"github.com/google/uuid"
)

type TransactionService interface {
	CollectRent(ctx context.Context, params CollectRentParams) (repo.Trade, error)
	// GetRecentPaidRent(ctx context.Context, params) ()
}

type transactionService struct {
	queries *repo.Queries
}

func NewTransactionService(queries *repo.Queries) *transactionService {
	return new(transactionService{queries})
}

type CollectRentParams struct {
	LeaseId uuid.UUID
}

func (s *transactionService) CollectRent(ctx context.Context, params CollectRentParams) (repo.Trade, error) {
	// get the current state of the lease
	lease, err := s.queries.GetLeaseById(ctx, params.LeaseId)
	if err != nil {
		return repo.Trade{}, err
	}

	// update the lease, move the expiry date to next month
	updatedLease, err := s.queries.UpdateLease(ctx, repo.UpdateLeaseParams{
		ID:         params.LeaseId,
		ExpiryDate: new(utils.ComputeNextExpiryDate(*lease.ExpiryDate, int(lease.ExpectedRentDay))),
	})
	if err != nil {
		return repo.Trade{}, err
	}

	// record the transaction (trade table)
	trade, err := s.queries.CreateTrade(ctx, repo.CreateTradeParams{
		LeaseID:    lease.ID,
		Type:       int16(constants.RENT),
		PaidAmount: lease.PropertyRentAmount.Decimal,
		StartDate:  lease.ExpiryDate,
		EndDate:    updatedLease.ExpiryDate,
	})
	if err != nil {
		return repo.Trade{}, err
	}

	return trade, nil
}
