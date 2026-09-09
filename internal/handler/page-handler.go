package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/Se7enSe7enSe7en/tenant-manager/internal/ctxkeys"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/errs"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/service"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/utils"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/web/component/propertycard"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/web/component/tenantcard"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/web/page"
	"github.com/google/uuid"
)

type PageHandler struct {
	PropertyService service.PropertyService
	LeaseService    service.LeaseService
}

func NewPageHandler(propertyService service.PropertyService, leaseService service.LeaseService) *PageHandler {
	return new(PageHandler{
		PropertyService: propertyService,
		LeaseService:    leaseService,
	})
}

func (h *PageHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	page.LoginPage().Render(r.Context(), w)
}

func (h *PageHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	page.RegisterPage().Render(r.Context(), w)
}

func (h *PageHandler) DashboardPage(w http.ResponseWriter, r *http.Request) {
	// get user, "ok" is not needed since this handler is already protected from the middleware
	user, ok := ctxkeys.UserFrom(r.Context())
	if !ok {
		http.Error(w, "no user", http.StatusInternalServerError)
		return
	}

	dbLeaseList, err := h.LeaseService.ListLease(r.Context(), user.ID)
	if err != nil {
		errs.Http(w, r, err, http.StatusInternalServerError)
		return
	}

	// convert []repo.Tenant -> []component.TenantCardProps
	tenantList := make([]tenantcard.TenantCardProps, len(dbLeaseList))
	for i, l := range dbLeaseList {
		tenantList[i] = tenantcard.TenantCardProps{
			Id:   l.TenantID.UUID.String(),
			Name: *l.TenantName,
			Unit: *l.PropertyName,
			// Status: , // TODO: add status
			RentAmount:  l.PropertyRentAmount.Decimal.String(), // TODO: should get from property as well
			NextDueDate: l.ExpiryDate.Format("Jan-02-2006"),
			Email:       l.TenantEmail,
			PhoneNumber: l.TenantPhoneNumber,
			LeaseId:     l.ID.String(),
		}
	}

	// property from db
	dbPropertyList, err := h.PropertyService.ListUnoccupiedProperty(r.Context(), user.ID)
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
			RentAmount: dbProperty.RentAmount.String(),
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

	// type conversion
	propertyIdUuid, err := uuid.Parse(propertyId)
	if err != nil {
		errs.Http(w, r, err, http.StatusInternalServerError)
		return
	}

	property, err := h.PropertyService.GetProperty(r.Context(), propertyIdUuid)
	if err != nil {
		errs.Http(w, r, err, http.StatusInternalServerError)
		return
	}

	page.CreateTenantPage(page.CreateTenantPageProps{
		PropertyId:           propertyId,
		DefaultDepositAmount: property.RentAmount.InexactFloat64(),
	}).Render(r.Context(), w)
}

func (h *PageHandler) CreateTradePage(w http.ResponseWriter, r *http.Request) {
	leaseId := r.URL.Query().Get("lease_id")

	leaseIdUuid, err := uuid.Parse(leaseId)
	if err != nil {
		errs.Http(w, r, err, http.StatusInternalServerError)
		return
	}

	// get tenant details using id
	// get property details from the connected property
	leaseDetails, err := h.LeaseService.GetLease(r.Context(), leaseIdUuid)
	if err != nil {
		errs.Http(w, r, err, http.StatusInternalServerError)
		return
	}

	page.CreateTradePage(page.CreateTradePageProps{
		LeaseId:      leaseId,
		TenantName:   *leaseDetails.TenantName,
		PropertyName: *leaseDetails.PropertyName,

		// current month until next month (prev rent day to next rent day)
		RentValidityPeriod: leaseDetails.ExpiryDate.Format("Jan-02-2006") + " to " + utils.ComputeNextExpiryDate(*leaseDetails.ExpiryDate, int(leaseDetails.ExpectedRentDay)).Format("Jan-02-2006"),
		RentAmount:         leaseDetails.PropertyRentAmount.Decimal.String(),
		TransactionDate:    time.Now().Format("Jan-02-2006"),
	}).Render(r.Context(), w)
}
