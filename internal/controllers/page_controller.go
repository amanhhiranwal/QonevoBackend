package controllers

import (
	"html/template"
	"log"
	"net/http"

	"qonevo-backend/internal/services"
)

type PageController struct {
	tpl            *template.Template
	productService *services.ProductService
}

func NewPageController(
	productService *services.ProductService,
) *PageController {

	// =====================================
	// Parse all templates
	// =====================================

	tpl := template.Must(
		template.ParseGlob("templates/**/*.html"),
	)

	return &PageController{
		tpl:            tpl,
		productService: productService,
	}
}

// =====================================
// Shared render helper
// =====================================

func (p *PageController) render(
	w http.ResponseWriter,
	name string,
	data any,
) {

	err := p.tpl.ExecuteTemplate(
		w,
		name,
		data,
	)

	if err != nil {

		log.Printf(
			"template render error: %v",
			err,
		)

		http.Error(
			w,
			"template render error",
			http.StatusInternalServerError,
		)

		return
	}
}

// =====================================
// Login Page
// =====================================

func (p *PageController) LoginPage(
	w http.ResponseWriter,
	r *http.Request,
) {

	p.render(
		w,
		"login.html",
		map[string]any{
			"Title":      "Login",
			"IsLoggedIn": false,
		},
	)
}

// =====================================
// Register Page
// =====================================

func (p *PageController) RegisterPage(
	w http.ResponseWriter,
	r *http.Request,
) {

	p.render(
		w,
		"register.html",
		map[string]any{
			"Title":      "Register",
			"IsLoggedIn": false,
		},
	)
}

// =====================================
// Dashboard
// =====================================

func (p *PageController) Dashboard(
	w http.ResponseWriter,
	r *http.Request,
) {

	totalProducts, err := p.productService.CountProducts()

	if err != nil {

		log.Printf(
			"failed to count products: %v",
			err,
		)

		totalProducts = 0
	}

	userID := r.Context().Value("user_id")

	p.render(
		w,
		"dashboard.html",
		map[string]any{
			"Title":         "Dashboard",
			"IsLoggedIn":    true,
			"UserID":        userID,
			"TotalProducts": totalProducts,
		},
	)
}
