package handler

import (
	"net/http"
	"time"

	"github.com/Se7enSe7enSe7en/go-toolkit/pkg/logger"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/errs"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/service"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/store"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/utils"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/validation"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/starfederation/datastar-go/datastar"
)

type TenantHandler struct {
	registrationService service.RegistrationService
}

func NewTenantHandler(registrationService service.RegistrationService) *TenantHandler {
	return &TenantHandler{
		registrationService: registrationService,
	}
}

func (h *TenantHandler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	logger.Debug("CreateTenant() triggered")

	// parse using datastar signals
	var payload struct {
		CreateTenantSignals store.CreateTenantSignals `json:"create_tenant_signals"`
	}
	if err := datastar.ReadSignals(r, &payload); err != nil {
		errs.Http(w, r, err, http.StatusInternalServerError)
		return
	}

	form := payload.CreateTenantSignals

	// validation
	if err := validation.CheckCreateTenantForm(form); err != nil {
		// TODO: use datastar to show errors in the front end
		errs.Http(w, r, err, http.StatusInternalServerError) // TMP
		return
	}

	// string conversions
	propertyIdUuid, err := uuid.Parse(form.PropertyId)
	if err != nil {
		errs.Http(w, r, err, http.StatusInternalServerError)
		return
	}
	startDateTime, err := time.Parse("2006-01-02", form.StartDate)
	if err != nil {
		errs.Http(w, r, err, http.StatusInternalServerError)
		return
	}
	depositAmountDecimal := decimal.NewFromFloat(form.DepositAmount)

	_, _, err = h.registrationService.RegisterTenantWithLease(r.Context(), service.RegisterTenantWithLeaseParams{
		Email:           form.Email,
		Name:            form.Name,
		PhoneNumber:     form.PhoneNumber,
		PropertyID:      propertyIdUuid,
		ExpectedRentDay: form.ExpectedRentDay,
		StartDate:       &startDateTime,
		IsMonthAdvance:  form.IsMonthAdvance,
		DepositAmount:   depositAmountDecimal,
	})
	if err != nil {
		errs.Http(w, r, err, http.StatusInternalServerError)
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

func (h *TenantHandler) ComputeNextDueDate(w http.ResponseWriter, r *http.Request) {
	// parse using datastar signals
	var payload struct {
		CreateTenantSignals store.CreateTenantSignals `json:"create_tenant_signals"`
	}
	if err := datastar.ReadSignals(r, &payload); err != nil {
		errs.Http(w, r, err, http.StatusInternalServerError)
		return
	}

	// type conversions
	startDateTime, err := time.Parse("2006-01-02", payload.CreateTenantSignals.StartDate)
	if err != nil {
		errs.Http(w, r, err, http.StatusInternalServerError)
		return
	}

	// recompute next_due_date
	payload.CreateTenantSignals.NextDueDate = utils.ComputeNextExpiryDate(startDateTime, payload.CreateTenantSignals.ExpectedRentDay).Format("Jan 02 2006")

	sse := datastar.NewSSE(w, r)

	sse.MarshalAndPatchSignals(payload)
}
