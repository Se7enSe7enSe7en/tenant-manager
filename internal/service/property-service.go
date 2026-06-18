package service

import (
	"context"

	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
	"github.com/google/uuid"
)

type PropertyService interface {
	CreateProperty(ctx context.Context, params repo.CreatePropertyParams) (repo.Property, error)
	ListProperties(ctx context.Context, userId uuid.UUID) ([]repo.Property, error)
	ListUnoccupiedProperties(ctx context.Context, userId uuid.UUID) ([]repo.Property, error)
}

type propertyService struct {
	repo repo.Querier
}

func NewPropertyService(repo repo.Querier) *propertyService {
	return new(propertyService{repo})
}

func (s *propertyService) CreateProperty(ctx context.Context, params repo.CreatePropertyParams) (repo.Property, error) {
	return s.repo.CreateProperty(ctx, params)
}

func (s *propertyService) ListProperties(ctx context.Context, userId uuid.UUID) ([]repo.Property, error) {
	return s.repo.ListProperties(ctx, userId)
}

func (s *propertyService) ListUnoccupiedProperties(ctx context.Context, userId uuid.UUID) ([]repo.Property, error) {
	return s.repo.ListUnoccupiedProperties(ctx, userId)
}
