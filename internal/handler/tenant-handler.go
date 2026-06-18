package handler

import (
	"net/http"

	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/model"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/service"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/validation"
	"github.com/google/uuid"
	"github.com/starfederation/datastar-go/datastar"
)

type tenantHandler struct {
	service service.TenantService
}

func NewTenantHandler(service service.TenantService) *tenantHandler {
	return &tenantHandler{service}
}

func (h *tenantHandler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	// parse using datastar signals
	var payload struct {
		CreateTenantSignals model.CreateTenantSignals `json:"create_tenant_signals"`
	}
	if err := datastar.ReadSignals(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	form := payload.CreateTenantSignals

	// string conversions
	propertyIdUuid, err := uuid.Parse(form.PropertyId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// validation
	if err := validation.CheckCreateTenantForm(form); err != nil {
		// TODO: use datastar to show errors in the front end
		http.Error(w, err.Error(), http.StatusInternalServerError) // TMP
		return
	}

	// call tenant service
	_, err = h.service.CreateTenant(r.Context(), repo.CreateTenantParams{
		Email:           form.Email,
		Name:            form.Name,
		PhoneNumber:     form.PhoneNumber,
		ExpectedRentDay: int16(form.ExpectedRentDay),
		PropertyID:      propertyIdUuid,
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
