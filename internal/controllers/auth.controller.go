package controllers

import (
	"encoding/json"
	"net/http"

	"qonevo-backend/internal/services"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(s *services.AuthService) *AuthController {
	return &AuthController{service: s}
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	json.NewDecoder(r.Body).Decode(&body)

	token, err := c.service.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"access_token": token,
	})
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	json.NewDecoder(r.Body).Decode(&body)

	user, err := c.service.Register(r.Context(), body.Email, body.Password)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(user)
}