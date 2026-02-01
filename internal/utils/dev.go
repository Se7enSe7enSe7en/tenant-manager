//go:build debug

package utils

import "net/http"

func DisableCacheInDevMode(next http.Handler) http.Handler {
	return next
}
