package handler

import (
	"net/http"

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
