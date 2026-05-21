package middleware

import (
	"net/http"

	"qonevo-backend/internal/utils"
)

func RequireAuth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		cookie, err := r.Cookie("token")

		if err != nil || cookie.Value == "" {

			http.Redirect(
				w,
				r,
				"/login",
				http.StatusSeeOther,
			)

			return
		}

		_, err = utils.ValidateToken(cookie.Value)

		if err != nil {

			http.SetCookie(w, &http.Cookie{
				Name:   "token",
				Value:  "",
				Path:   "/",
				MaxAge: -1,
			})

			http.Redirect(
				w,
				r,
				"/login",
				http.StatusSeeOther,
			)

			return
		}

		next.ServeHTTP(w, r)
	})
}


func RedirectIfAuthenticated(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		cookie, err := r.Cookie("token")

		if err == nil && cookie.Value != "" {

			_, err := utils.ValidateToken(cookie.Value)

			if err == nil {

				http.Redirect(
					w,
					r,
					"/dashboard",
					http.StatusSeeOther,
				)

				return
			}
		}

		next.ServeHTTP(w, r)
	})
}