package handler

import (
	"context"
	"net/http"

	"github.com/Se7enSe7enSe7en/go-toolkit/pkg/logger"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/service"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/web/component/tenantcard"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/web/page"
)

type tenantHandler struct {
	service service.TenantService
}

func NewTenantHandler(service service.TenantService) *tenantHandler {
	return &tenantHandler{service}
}

func (h *tenantHandler) ListTenantPage(w http.ResponseWriter, r *http.Request) {
	// call the service -> ListTenant
	dbTenantList, err := h.service.ListTenants(r.Context())
	if err != nil {
		logger.Error("ListTenantPage(): ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Debug("ListTenantPage() dbTenantList: %v", dbTenantList)

	// convert []repo.Tenant -> []component.TenantCardProps
	tenantList := make([]tenantcard.TenantCardProps, len(dbTenantList))
	for i, t := range dbTenantList {
		tenantList[i] = tenantcard.TenantCardProps{
			Name: t.Name,
			Unit: t.PropertyID.String(), // TODO: should be handled from the query, getting the property connected to the user
			//  Status: , // TODO: add status
			RentAmount: t.PropertyID.String(), // TODO: should get from property as well
			// LastPaymentDate: , TODO: get from last transaction
			Email:       new(t.Email),
			PhoneNumber: new(t.PhoneNumber),
		}
	}

	// return tenant page with context in an HTTP response
	page.MainPage(tenantList).Render(context.Background(), w)
}
