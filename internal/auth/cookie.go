package auth

import (
	"net/http"

	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
)

func SetCookie(w http.ResponseWriter, session repo.Session) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    session.ID.String(),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  session.ExpiresAt.Time,
		// Secure: true, // TODO: for prod only
	})
}

func ClearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1, // note: -1 tells the browser to delete the cookie, hence "Expires" is no longer needed to be populated as well
	})
}
