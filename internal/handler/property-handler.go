package handler

import (
	"net/http"

	"github.com/Se7enSe7enSe7en/tenant-manager/internal/ctxkeys"
	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/service"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/utils"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/validation"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/web/page"
	"github.com/starfederation/datastar-go/datastar"
)

type propertyHandler struct {
	service service.PropertyService
}

func NewPropertyHandler(service service.PropertyService) *propertyHandler {
	return new(propertyHandler{service})
}

func (h *propertyHandler) CreatePropertyPage(w http.ResponseWriter, r *http.Request) {
	page.CreatePropertyPage().Render(r.Context(), w)
}

func (h *propertyHandler) CreateProperty(w http.ResponseWriter, r *http.Request) {
	// get user from context added by AttachUser middleware
	user, ok := ctxkeys.UserFrom(r.Context())
	if !ok {
		http.Error(w, "No user", http.StatusUnauthorized) // TODO: check later if this is the correct http code to use here
		return
	}

	// parse the form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// get form values
	name := r.FormValue("name")
	rentAmount := r.FormValue("rent_amount")

	// TODO: make adapters or combine with validation step to make these more straight forward
	// string conversions
	rentAmountNumeric, err := utils.StringToPgtypeNumeric(rentAmount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: validate form values, this function is not yet complete
	if err := validation.CheckCreatePropertyForm(validation.CreatePropertyForm{
		Name:       name,
		RentAmount: float64(rentAmountNumeric.Exp),
	}); err != nil {
		// TODO: handle error, with datastar
		http.Error(w, err.Error(), http.StatusInternalServerError) // TMP
		return
	}

	// call property service
	_, err = h.service.CreateProperty(r.Context(), repo.CreatePropertyParams{
		UserID:     user.ID,
		Name:       name,
		RentAmount: rentAmountNumeric,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: send a success signal
	sse := datastar.NewSSE(w, r)

	type Output struct {
		Success bool `json:"success"`
	}

	sse.MarshalAndPatchSignals(Output{Success: true})

	// TODO: change later to show a modal then make the user choose to add another or go back to the dashboard
	// change auto redirect later
	sse.Redirect("/dashboard") // TMP
}

// TOCHECK
