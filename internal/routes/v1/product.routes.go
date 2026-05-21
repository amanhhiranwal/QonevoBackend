package v1

import (
	"net/http"

	"qonevo-backend/internal/controllers"
	"qonevo-backend/internal/middleware"
)

func ProductRoutes(
	mux *http.ServeMux,
	product *controllers.ProductController,
) {

	// =====================================
	// Products Listing
	// =====================================

	mux.Handle(
		"/products",
		middleware.RequireAuth(
			http.HandlerFunc(func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				switch r.Method {

				case http.MethodGet:
					product.Index(w, r)

				case http.MethodPost:
					product.Store(w, r)

				default:

					http.Error(
						w,
						"method not allowed",
						http.StatusMethodNotAllowed,
					)
				}
			}),
		),
	)

	// =====================================
	// Create Product Page
	// =====================================

	mux.Handle(
		"/products/create",
		middleware.RequireAuth(
			http.HandlerFunc(
				product.CreatePage,
			),
		),
	)

	// =====================================
	// Edit Product Page
	// =====================================

	mux.Handle(
		"/products/edit/",
		middleware.RequireAuth(
			http.HandlerFunc(
				product.EditPage,
			),
		),
	)

	// =====================================
	// Update Product
	// =====================================

	mux.Handle(
		"/products/update/",
		middleware.RequireAuth(
			http.HandlerFunc(
				product.Update,
			),
		),
	)

	// =====================================
	// Delete Product
	// =====================================

	mux.Handle(
		"/products/delete/",
		middleware.RequireAuth(
			http.HandlerFunc(
				product.Delete,
			),
		),
	)
	mux.Handle(
		"/products/image/delete/",
		middleware.RequireAuth(
			http.HandlerFunc(
				product.DeleteImage,
			),
		),
	)
}
