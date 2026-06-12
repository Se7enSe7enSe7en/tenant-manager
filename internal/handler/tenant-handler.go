package handler

import (
	"net/http"

	"github.com/Se7enSe7enSe7en/go-toolkit/pkg/logger"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/ctxkeys"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/service"
	"github.com/starfederation/datastar-go/datastar"
)

type tenantHandler struct {
	service service.TenantService
}

func NewTenantHandler(service service.TenantService) *tenantHandler {
	return &tenantHandler{service}
}

func (h *tenantHandler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	// TODO: move this into auth util
	// get user from context added by AttachUser middleware
	user, ok := ctxkeys.UserFrom(r.Context())
	if !ok {
		http.Error(w, "No user", http.StatusUnauthorized) // TODO: check later if this is the correct http code to use here
		return
	}

	// // parse query param
	// propertyId := r.URL.Query().Get("property_id")

	// // parse form data
	// name := r.FormValue("name")
	// email := r.FormValue("email")
	// phoneNumber := r.FormValue("phone_number")
	// expectedRentDay := r.FormValue("expected_rent_day")

	// parse using datastar
	type Payload struct {
		PropertyId string `json:"property_id"`
		Tenant     struct {
			Name            string `json:"name"`
			Email           string `json:"email"`
			PhoneNumber     string `json:"phone_number"`
			ExpectedRentDay string `json:"expected_rent_day"`
		} `json:"tenant"`
	}

	var payload Payload

	if err := datastar.ReadSignals(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// DEBUG
	logger.Debug("user: %v", user)
	logger.Debug("propertyId: %v", payload.PropertyId)
	logger.Debug("name: %v", payload.Tenant.Name)
	logger.Debug("email: %v", payload.Tenant.Email)
	logger.Debug("phoneNumber: %v", payload.Tenant.PhoneNumber)
	logger.Debug("expectedRentDay: %v", payload.Tenant.ExpectedRentDay)

	// validation

}
