package handler

import (
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/service"
)

type tenantHandler struct {
	service service.TenantService
}

func NewTenantHandler(service service.TenantService) *tenantHandler {
	return &tenantHandler{service}
}
