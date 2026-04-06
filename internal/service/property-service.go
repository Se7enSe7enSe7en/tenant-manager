package service

import (
	"context"

	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
)

type PropertyService interface {
	ListProperties(ctx context.Context) ([]repo.Property, error)
}

type propertyService struct {
	repo repo.Querier
}

func NewPropertyService(repo repo.Querier) *propertyService {
	return new(propertyService{repo})
}

func (s *propertyService) ListProperties(ctx context.Context) ([]repo.Property, error) {
	return s.repo.ListProperties(ctx)
}
