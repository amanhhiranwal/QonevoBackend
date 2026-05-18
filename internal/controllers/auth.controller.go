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

// func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
// 	var body struct {
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}

// 	json.NewDecoder(r.Body).Decode(&body)

// 	user, err := c.service.Register(r.Context(), body.Email, body.Password)
// 	if err != nil {
// 		http.Error(w, err.Error(), 400)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(user)
// }

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		FirstName string  `json:"first_name"`
		LastName  string  `json:"last_name"`
		Email     string  `json:"email"`
		Password  string  `json:"password"`
		Phone     *string `json:"phone,omitempty"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "invalid request body", 400)
		return
	}

	user, err := c.service.Register(
		r.Context(),
		body.FirstName,
		body.LastName,
		body.Email,
		body.Password,
		body.Phone,
	)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// IMPORTANT: never expose password_hash
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"phone":      user.Phone,
		"created_at": user.CreatedAt,
	})
}