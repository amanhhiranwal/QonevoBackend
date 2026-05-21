package v1

import (
	"net/http"

	"qonevo-backend/internal/controllers"
	"qonevo-backend/internal/middleware"
)

func PageRoutes(
	mux *http.ServeMux,
	page *controllers.PageController,
	auth *controllers.AuthController,
) {

	// =========================
	// Public routes
	// =========================

	mux.Handle(
		"/login",
		middleware.RedirectIfAuthenticated(
			http.HandlerFunc(page.LoginPage),
		),
	)

	mux.Handle(
		"/register",
		middleware.RedirectIfAuthenticated(
			http.HandlerFunc(page.RegisterPage),
		),
	)

	// =========================
	// Protected routes
	// =========================

	mux.Handle(
		"/dashboard",
		middleware.RequireAuth(
			http.HandlerFunc(page.Dashboard),
		),
	)

	// =========================
	// Form handlers
	// =========================

	mux.HandleFunc(
		"/login-post",
		auth.LoginPagePost,
	)

	mux.HandleFunc(
		"/register-post",
		auth.RegisterPagePost,
	)
	mux.HandleFunc(
		"/logout",
		auth.Logout,
	)
}