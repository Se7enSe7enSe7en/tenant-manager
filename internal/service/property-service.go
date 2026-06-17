package service

import (
	"context"

	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/utils"
)

type PropertyService interface {
	ListProperties(ctx context.Context) ([]repo.Property, error)
	CreateProperty(ctx context.Context, params repo.CreatePropertyParams) (repo.Property, error)
	ListUnoccupiedProperties(ctx context.Context, userId string) ([]repo.Property, error)
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

func (s *propertyService) CreateProperty(ctx context.Context, params repo.CreatePropertyParams) (repo.Property, error) {
	return s.repo.CreateProperty(ctx, params)
}

func (s *propertyService) ListUnoccupiedProperties(ctx context.Context, userId string) ([]repo.Property, error) {
	userIdPgtypeUuid, err := utils.StringToPgtypeUuid(userId)
	if err != nil {
		return nil, err
	}

	return s.repo.ListUnoccupiedProperties(ctx, userIdPgtypeUuid)
}
