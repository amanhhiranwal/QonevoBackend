package middleware

import (
	"net/http"

	"qonevo-backend/internal/utils"
)

func RedirectIfAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("token")
		if err == nil && cookie.Value != "" {

			// validate token
			_, err := utils.ParseToken(cookie.Value)
			if err == nil {
				// already logged in → go dashboard
				http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}