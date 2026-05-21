package v1

import (
	"net/http"
	"os"
	"path/filepath"

	"qonevo-backend/internal/controllers"
	"qonevo-backend/internal/middleware"
)

func RegisterV1Routes(
	mux *http.ServeMux,
	auth *controllers.AuthController,
	page *controllers.PageController,
	product *controllers.ProductController,
) {

	// ============================================
	// STATIC FILES
	// ============================================

	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	staticPath := filepath.Join(
		wd,
		"static",
	)

	fs := http.FileServer(
		http.Dir(staticPath),
	)

	mux.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			fs,
		),
	)

	// ============================================
	// DEBUG ROUTE
	// ============================================

	mux.HandleFunc(
		"/debug",
		func(w http.ResponseWriter, r *http.Request) {

			w.Write([]byte(
				"server is running",
			))
		},
	)

	// ============================================
	// AUTH ROUTES
	// ============================================

	AuthRoutes(
		mux,
		auth,
	)

	// ============================================
	// PAGE ROUTES
	// ============================================

	PageRoutes(
		mux,
		page,
		auth,
	)

	// ============================================
	// PRODUCT ROUTES
	// ============================================

	ProductRoutes(
		mux,
		product,
	)

	// ============================================
	// PRODUCT API ROUTES
	// ============================================

	ProductAPIRoutes(
		mux,
		product,
	)

	// ============================================
	// FALLBACK 404
	// ============================================

	mux.Handle(
		"/",
		middleware.NotFoundHandler(),
	)
}
