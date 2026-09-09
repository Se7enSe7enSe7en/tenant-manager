// TODO: refactor this

package service

import (
	"context"

	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
	"github.com/google/uuid"
)

type PropertyService interface {
	CreateProperty(ctx context.Context, params repo.CreatePropertyParams) (repo.Property, error)
	ListProperty(ctx context.Context, userId uuid.UUID) ([]repo.Property, error)
	ListUnoccupiedProperty(ctx context.Context, userId uuid.UUID) ([]repo.Property, error)
	GetProperty(ctx context.Context, propertyId uuid.UUID) (repo.GetPropertyByIdRow, error)
}

type propertyService struct {
	queries *repo.Queries
}

func NewPropertyService(repo *repo.Queries) *propertyService {
	return new(propertyService{repo})
}

func (s *propertyService) CreateProperty(ctx context.Context, params repo.CreatePropertyParams) (repo.Property, error) {
	return s.queries.CreateProperty(ctx, params)
}

func (s *propertyService) ListProperty(ctx context.Context, userId uuid.UUID) ([]repo.Property, error) {
	return s.queries.ListProperty(ctx, userId)
}

func (s *propertyService) ListUnoccupiedProperty(ctx context.Context, userId uuid.UUID) ([]repo.Property, error) {
	return s.queries.ListUnoccupiedProperty(ctx, userId)
}

func (s *propertyService) GetProperty(ctx context.Context, propertyId uuid.UUID) (repo.GetPropertyByIdRow, error) {
	return s.queries.GetPropertyById(ctx, propertyId)
}
