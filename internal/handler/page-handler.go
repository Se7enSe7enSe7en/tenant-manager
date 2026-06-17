package handler

import (
	"context"
	"net/http"

	"github.com/Se7enSe7enSe7en/tenant-manager/internal/ctxkeys"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/service"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/utils"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/web/component/propertycard"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/web/component/tenantcard"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/web/page"
)

type PageHandler struct {
	PropertyService service.PropertyService
	TenantService   service.TenantService
}

func NewPageHandler(h PageHandler) *PageHandler {
	return &h
}

func (h *PageHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	page.LoginPage().Render(r.Context(), w)
}

func (h *PageHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	page.RegisterPage().Render(r.Context(), w)
}

func (h *PageHandler) DashboardPage(w http.ResponseWriter, r *http.Request) {
	// get user
	user, ok := ctxkeys.UserFrom(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// call the service -> ListTenant
	dbTenantList, err := h.TenantService.ListTenantsWithProperty(r.Context(), user.ID.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// logger.Debug("ListTenantPage() dbTenantList: %v", dbTenantList)

	// convert []repo.Tenant -> []component.TenantCardProps
	tenantList := make([]tenantcard.TenantCardProps, len(dbTenantList))
	for i, t := range dbTenantList {
		tenantList[i] = tenantcard.TenantCardProps{
			Name: t.TenantName,
			Unit: t.PropertyName.String,
			// Status: , // TODO: add status
			RentAmount: utils.PgtypeNumericToString(t.PropertyRentAmount), // TODO: should get from property as well
			// LastPaymentDate: , TODO: get from last transaction
			Email:       &t.TenantEmail,
			PhoneNumber: &t.TenantPhoneNumber,
		}
	}

	// property from db
	dbPropertyList, err := h.PropertyService.ListUnoccupiedProperties(r.Context(), user.ID.String())
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
		TenantList:   tenantList,
	}).Render(context.Background(), w)
}

func (h *PageHandler) CreatePropertyPage(w http.ResponseWriter, r *http.Request) {
	page.CreatePropertyPage().Render(r.Context(), w)
}

func (h *PageHandler) CreateTenantPage(w http.ResponseWriter, r *http.Request) {
	propertyId := r.URL.Query().Get("property_id")

	page.CreateTenantPage(page.CreateTenantPageProps{
		PropertyId: propertyId,
	}).Render(r.Context(), w)
}
