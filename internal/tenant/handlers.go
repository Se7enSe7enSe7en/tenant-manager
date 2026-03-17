package tenant

import (
	"context"
	"net/http"

	web "github.com/Se7enSe7enSe7en/tenant-manager/web/components"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func ListTenantPage(w http.ResponseWriter, r *http.Request) {
	// // call the service -> ListTenant
	// tenantService := NewService()

	// use service get the tenant list

	// put tenant list in context

	// return tenant page with context in an HTTP response
	web.MainPage().Render(context.Background(), w)
}
