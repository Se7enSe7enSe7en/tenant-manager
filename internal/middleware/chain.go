package middleware

import "net/http"

func Chain(h http.Handler, ml ...func(http.Handler) http.Handler) http.Handler {
	for i := len(ml) - 1; i >= 0; i-- {
		h = ml[i](h)
	}

	return h
}
