package handler

import (
	"net/http"

	"github.com/Se7enSe7enSe7en/tenant-manager/internal/web/page"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	page.LoginPage().Render(r.Context(), w)
}
