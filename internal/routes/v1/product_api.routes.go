package v1

import (
	"net/http"

	"qonevo-backend/internal/controllers"
)

func ProductAPIRoutes(
	mux *http.ServeMux,
	product *controllers.ProductController,
) {

	// =====================================
	// GET ALL PRODUCTS API
	// =====================================

	mux.HandleFunc(
		"/api/v1/products",
		func(w http.ResponseWriter, r *http.Request) {

			switch r.Method {

			case http.MethodGet:
				product.GetProductsAPI(w, r)

			default:

				http.Error(
					w,
					"method not allowed",
					http.StatusMethodNotAllowed,
				)
			}
		},
	)

	mux.HandleFunc(
		"/api/v1/products/ifp/sizes",
		func(w http.ResponseWriter, r *http.Request) {

			switch r.Method {

			case http.MethodGet:
				product.GetIFPSizesAPI(w, r)

			default:
				http.Error(
					w,
					"method not allowed",
					http.StatusMethodNotAllowed,
				)
			}
		},
	)

	mux.HandleFunc(
		"/api/v1/ifp/filters",
		func(w http.ResponseWriter, r *http.Request) {

			switch r.Method {

			case http.MethodGet:
				product.GetIFPFiltersAPI(w, r)

			default:
				http.Error(
					w,
					"method not allowed",
					http.StatusMethodNotAllowed,
				)
			}
		},
	)
}
