package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"qonevo-backend/internal/config"
	"qonevo-backend/internal/services"
)

type AuthController struct {
	service *services.AuthService
	config  *config.Config
}

func NewAuthController(
	service *services.AuthService,
	cfg *config.Config,
) *AuthController {

	return &AuthController{
		service: service,
		config:  cfg,
	}
}

// =========================
// Register
// =========================

func (a *AuthController) Register(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	firstName := strings.TrimSpace(r.FormValue("first_name"))
	lastName := strings.TrimSpace(r.FormValue("last_name"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))

	phoneValue := strings.TrimSpace(r.FormValue("phone"))

	var phone *string

	if phoneValue != "" {
		phone = &phoneValue
	}

	// =========================
	// Validation
	// =========================

	if firstName == "" ||
		lastName == "" ||
		email == "" ||
		password == "" {

		http.Error(
			w,
			"all required fields must be provided",
			http.StatusBadRequest,
		)
		return
	}

	if len(password) < 8 {
		http.Error(
			w,
			"password must be at least 8 characters",
			http.StatusBadRequest,
		)
		return
	}

	// =========================
	// Create user
	// =========================

	err := a.service.Register(
		r.Context(),
		firstName,
		lastName,
		email,
		phone,
		password,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "user registered successfully",
	})
}

// =========================
// Login
// =========================

func (a *AuthController) Login(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))

	token, err := a.service.Login(
		r.Context(),
		email,
		password,
	)

	if err != nil {
		http.Error(
			w,
			"invalid credentials",
			http.StatusUnauthorized,
		)
		return
	}

	secure := a.config.AppEnv == "production"

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   3600,
	})

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "login successful",
	})
}

func (a *AuthController) LoginPagePost(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))

	token, err := a.service.Login(
		r.Context(),
		email,
		password,
	)

	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	secure := a.config.AppEnv == "production"

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   3600,
	})

	http.Redirect(
		w,
		r,
		"/dashboard",
		http.StatusSeeOther,
	)
}


func (a *AuthController) RegisterPagePost(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	firstName := strings.TrimSpace(r.FormValue("first_name"))
	lastName := strings.TrimSpace(r.FormValue("last_name"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))

	phoneValue := strings.TrimSpace(r.FormValue("phone"))

	var phone *string

	if phoneValue != "" {
		phone = &phoneValue
	}

	err := a.service.Register(
		r.Context(),
		firstName,
		lastName,
		email,
		phone,
		password,
	)

	if err != nil {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	http.Redirect(
		w,
		r,
		"/login",
		http.StatusSeeOther,
	)
}

func (a *AuthController) Logout(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {

		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	http.Redirect(
		w,
		r,
		"/login",
		http.StatusSeeOther,
	)
}