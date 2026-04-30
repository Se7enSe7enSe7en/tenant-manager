package handler

import (
	"net/http"

	"github.com/Se7enSe7enSe7en/go-toolkit/pkg/logger"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/service"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/web/page"
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
	// parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// get form values
	name := r.FormValue("name")
	rentAmount := r.FormValue("rent_amount")

	logger.Debug("CreateProperty() handler:")
	logger.Debug("name: %v", name)
	logger.Debug("rentAmount: %v", rentAmount)

	// TODO: validate form values
	// TODO: convert form values
	// TODO: call createProperty service
}
