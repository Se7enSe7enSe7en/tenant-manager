package handler

import (
	"errors"
	"net/http"

	"github.com/Se7enSe7enSe7en/go-toolkit/pkg/logger"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/auth"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/service"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/validation"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/web/page"
	"github.com/starfederation/datastar-go/datastar"
)

type authHandler struct {
	svc service.AuthService
}

func NewAuthHandler(svc service.AuthService) *authHandler {
	return &authHandler{svc: svc}
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	_, session, err := h.svc.Login(r.Context(), email, password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {

			sse := datastar.NewSSE(w, r)
			sse.PatchSignals([]byte(`{"loginError": "Invalid email or password"}`))

			return
		}

		logger.Error("auth-handler Login(): ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	auth.SetCookie(w, session)
	// (!) warning: set cookie writes stuff in the header

	sse := datastar.NewSSE(w, r)
	// (!) warning: datastar.NewSSE() flushes the headers when instantiating,
	// this means as much as possible do this after writing stuff in the header

	sse.Redirect("/dashboard")

}

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	name := r.FormValue("name") // optional (returns "" if empty)

	// validation
	if err := validation.RegisterInput(email, password, name); err != nil {
		// TODO: show inline err message through datastar in the form message
		http.Error(w, "invalid input", http.StatusInternalServerError) // tmp: remove this after setting up datastar + UI
		return
	}

	_, session, err := h.svc.Register(r.Context(), email, password, name)
	if err != nil {
		if errors.Is(err, auth.ErrEmailAlreadyTaken) {
			// TODO: prompt email already taken
			http.Error(w, "invalid input", http.StatusInternalServerError) // tmp: remove this after setting up datastar + UI
			return
		}

		logger.Error("auth-handler Register(): ", err)
		http.Error(w, "Register failed", http.StatusInternalServerError)
		return
	}

	auth.SetCookie(w, session) // (!) warning: same warning above

	sse := datastar.NewSSE(w, r) // (!)
	sse.Redirect("/dashboard")
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(auth.SessionCookieName)
	if err == nil {
		if err := h.svc.Logout(r.Context(), cookie.Value); err != nil {
			logger.Error("auth-handler Logout(): ", err)
			// continue to clear cookie and redirect
		}
	}

	auth.ClearCookie(w) // (!)

	sse := datastar.NewSSE(w, r) // (!)
	sse.Redirect("/login")
}

func (h *authHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	page.LoginPage().Render(r.Context(), w)
}

func (h *authHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	page.RegisterPage().Render(r.Context(), w)
}
