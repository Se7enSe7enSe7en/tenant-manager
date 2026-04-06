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
	// TODO: get properties from db
	_, err := h.service.ListProperties(r.Context())
	if err != nil {
		logger.Error("CreatePropertyPage(): ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	page.CreatePropertyPage().Render(r.Context(), w)
}
