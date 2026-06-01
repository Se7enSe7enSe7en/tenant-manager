package middleware

import (
	"net/http"

	"github.com/Se7enSe7enSe7en/tenant-manager/internal/auth"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/ctxkeys"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/service"
)

func AttachUser(svc service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(auth.SessionCookieName)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			user, err := svc.UserFromSession(r.Context(), cookie.Value)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctxkeys.WithUser(r.Context(), user)))
		})
	}
}

func RequireAuth(next http.Handler) http.Handler { // ?: can't I just change the input and output types into this struct instead?
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := ctxkeys.UserFrom(r.Context())
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
