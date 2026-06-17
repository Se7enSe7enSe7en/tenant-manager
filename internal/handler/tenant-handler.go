package handler

import (
	"net/http"

	"github.com/Se7enSe7enSe7en/go-toolkit/pkg/logger"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/ctxkeys"
	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/model"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/service"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/utils"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/validation"
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

	// parse using datastar signals
	var payload struct {
		CreateTenantSignals model.CreateTenantSignals `json:"create_tenant_signals"`
	}
	if err := datastar.ReadSignals(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	form := payload.CreateTenantSignals

	// DEBUG
	logger.Debug("user: %v", user)
	logger.Debug("propertyId: %v", form.PropertyId)
	logger.Debug("name: %v", form.Name)
	logger.Debug("email: %v", form.Email)
	logger.Debug("phoneNumber: %v", form.PhoneNumber)
	logger.Debug("expectedRentDay: %v", form.ExpectedRentDay)

	// string conversions
	propertyIdPgTypeUuid, err := utils.StringToPgtypeUuid(form.PropertyId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// validation
	if err := validation.CheckCreateTenantForm(form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // TMP
		return
	}

	// call tenant service
	_, err = h.service.CreateTenant(r.Context(), repo.CreateTenantParams{
		Email:           form.Email,
		Name:            form.Name,
		PhoneNumber:     form.PhoneNumber,
		ExpectedRentDay: int16(form.ExpectedRentDay),
		PropertyID:      propertyIdPgTypeUuid,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: send a success signal
	sse := datastar.NewSSE(w, r)

	type Response struct {
		Success bool `json:"success"`
	}

	sse.MarshalAndPatchSignals(Response{Success: true})

	// TODO: same as property-handler, show modal
	sse.Redirect("/dashboard") // TMP
}
