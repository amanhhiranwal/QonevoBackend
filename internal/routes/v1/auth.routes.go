package v1

import (
	"net/http"

	"qonevo-backend/internal/controllers"
)

func AuthRoutes(mux *http.ServeMux, c *controllers.AuthController) {
	mux.HandleFunc("/api/v1/auth/register", c.Register)
	mux.HandleFunc("/api/v1/auth/login", c.Login)
}