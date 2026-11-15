package handler

import (
	"net/http"

	"github.com/Se7enSe7enSe7en/tenant-manager/internal/ctxkeys"
	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/errs"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/service"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/store"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/validation"
	"github.com/shopspring/decimal"
	"github.com/starfederation/datastar-go/datastar"
)

type propertyHandler struct {
	service service.PropertyService
}

func NewPropertyHandler(service service.PropertyService) *propertyHandler {
	return new(propertyHandler{service})
}

func (h *propertyHandler) CreateProperty(w http.ResponseWriter, r *http.Request) {
	// get user from context added by AttachUser middleware
	user, ok := ctxkeys.UserFrom(r.Context())
	if !ok {
		http.Error(w, "No user", http.StatusUnauthorized) // TODO: check later if this is the correct http code to use here
		return
	}

	// parse signals
	var payload struct {
		CreatePropertySignals store.CreatePropertySignals `json:"create_property_signals"`
	}
	err := datastar.ReadSignals(r, &payload)
	if err != nil {
		errs.Http(w, r, err, http.StatusInternalServerError)
		return
	}

	// // parse the form data
	// if err := r.ParseForm(); err != nil {
	// 	http.Error(w, "Unable to parse form", http.StatusBadRequest)
	// 	return
	// }
	// // get form values
	// name := r.FormValue("name")
	// rentAmount := r.FormValue("rent_amount")

	// // TODO: make adapters or combine with validation step to make these more straight forward
	// // string conversions
	// rentAmountDecimal, err := decimal.NewFromString(rentAmount)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// TODO: validate form values, this function is not yet complete
	// if err := validation.CheckCreatePropertyForm(validation.CreatePropertyForm{
	// 	Name:       name,
	// 	RentAmount: rentAmountDecimal,
	// }); err != nil {
	// 	// TODO: handle error, with datastar
	// 	http.Error(w, err.Error(), http.StatusInternalServerError) // TMP
	// 	return
	// }
	err = validation.CheckCreatePropertySignals(payload.CreatePropertySignals)
	if err != nil {
		sse := datastar.NewSSE(w, r)
		payload.CreatePropertySignals.FormError = err.Error()
		sse.MarshalAndPatchSignals(payload)
		return
	}

	// call property service
	_, err = h.service.CreateProperty(r.Context(), repo.CreatePropertyParams{
		UserID:     user.ID,
		Name:       payload.CreatePropertySignals.Name,
		RentAmount: decimal.NewFromFloat(payload.CreatePropertySignals.RentAmount),
		Address:    &payload.CreatePropertySignals.Address,
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

	// TODO: change later to show a modal then make the user choose to add another or go back to the dashboard
	// change auto redirect later
	sse.Redirect("/dashboard") // TMP
}

// TOCHECK
