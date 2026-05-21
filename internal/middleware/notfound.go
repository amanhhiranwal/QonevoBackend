package middleware

import (
	"net/http"
)

func NotFoundHandler() http.Handler {

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		http.NotFound(w, r)
	})
}