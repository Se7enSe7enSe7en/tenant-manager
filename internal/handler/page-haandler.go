package handler

import (
	"context"
	"net/http"

	"github.com/Se7enSe7enSe7en/tenant-manager/internal/service"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/utils"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/web/component/propertycard"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/web/page"
)

type PageHandler struct {
	PropertyService service.PropertyService
	TenantService   service.TenantService
}

func NewPageHandler(h PageHandler) *PageHandler {
	return &h
}

func (h *PageHandler) DashboardPage(w http.ResponseWriter, r *http.Request) {
	// // call the service -> ListTenant
	// dbTenantList, err := h.tenantService.ListTenants(r.Context())
	// if err != nil {
	// 	logger.Error("ListTenantPage(): ", err)
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// logger.Debug("ListTenantPage() dbTenantList: %v", dbTenantList)

	// // convert []repo.Tenant -> []component.TenantCardProps
	// tenantList := make([]tenantcard.TenantCardProps, len(dbTenantList))
	// for i, t := range dbTenantList {
	// 	tenantList[i] = tenantcard.TenantCardProps{
	// 		Name: t.Name,
	// 		Unit: t.PropertyID.String(), // TODO: should be handled from the query, getting the property connected to the user
	// 		//  Status: , // TODO: add status
	// 		RentAmount: t.PropertyID.String(), // TODO: should get from property as well
	// 		// LastPaymentDate: , TODO: get from last transaction
	// 		Email:       new(t.Email),
	// 		PhoneNumber: new(t.PhoneNumber),
	// 	}
	// }

	// property from db
	dbPropertyList, err := h.PropertyService.ListProperties(r.Context())
	if err != nil {
		http.Error(w, "cannot get property list", http.StatusInternalServerError)
		return
	}

	// convert dbProperty for property in front end
	propertyList := make([]propertycard.PropertyCardProps, len(dbPropertyList))
	for i, dbProperty := range dbPropertyList {

		propertyList[i] = propertycard.PropertyCardProps{
			Id:         dbProperty.ID.String(),
			Name:       dbProperty.Name,
			RentAmount: utils.PgtypeNumericToString(dbProperty.RentAmount),
		}
	}

	// return tenant page with context in an HTTP response
	page.DashboardPage(page.DashboardPageProps{
		PropertyList: propertyList,
		// TenantList: tenantList,

	}).Render(context.Background(), w)
}
