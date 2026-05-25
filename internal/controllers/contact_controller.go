package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"

	"qonevo-backend/internal/models"
	"qonevo-backend/internal/services"
)

type ContactController struct {
	service *services.ContactService
}

func NewContactController(
	service *services.ContactService,
) *ContactController {

	return &ContactController{
		service: service,
	}
}

type CreateContactRequest struct {
	FullName string `json:"fullName"`

	Email string `json:"email"`

	PhoneNumber string `json:"phoneNumber"`

	CompanyName string `json:"companyName"`

	WebsiteURL string `json:"websiteUrl"`

	HelpMessage string `json:"helpMessage"`
}

// =====================================
// CREATE CONTACT API
// =====================================

func (c *ContactController) CreateContact(
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

	var req CreateContactRequest

	err := json.NewDecoder(r.Body).Decode(
		&req,
	)

	if err != nil {

		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================
	// SANITIZE
	// =====================================

	req.FullName = strings.TrimSpace(
		req.FullName,
	)

	req.Email = strings.TrimSpace(
		strings.ToLower(req.Email),
	)

	req.PhoneNumber = strings.TrimSpace(
		req.PhoneNumber,
	)

	req.CompanyName = strings.TrimSpace(
		req.CompanyName,
	)

	req.WebsiteURL = strings.TrimSpace(
		req.WebsiteURL,
	)

	req.HelpMessage = strings.TrimSpace(
		req.HelpMessage,
	)

	// =====================================
	// VALIDATION
	// =====================================

	if req.FullName == "" {

		http.Error(
			w,
			"full name is required",
			http.StatusBadRequest,
		)

		return
	}

	if req.Email == "" {

		http.Error(
			w,
			"email is required",
			http.StatusBadRequest,
		)

		return
	}

	emailRegex := regexp.MustCompile(
		`^[^\s@]+@[^\s@]+\.[^\s@]+$`,
	)

	if !emailRegex.MatchString(
		req.Email,
	) {

		http.Error(
			w,
			"invalid email",
			http.StatusBadRequest,
		)

		return
	}

	phoneRegex := regexp.MustCompile(
		`^\d{10}$`,
	)

	if !phoneRegex.MatchString(
		req.PhoneNumber,
	) {

		http.Error(
			w,
			"phone number must be 10 digits",
			http.StatusBadRequest,
		)

		return
	}

	if req.HelpMessage == "" {

		http.Error(
			w,
			"message is required",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================
	// BUILD MODEL
	// =====================================

	contact := &models.Contact{
		FullName: req.FullName,

		Email: req.Email,

		PhoneNumber: req.PhoneNumber,

		CompanyName: req.CompanyName,

		WebsiteURL: req.WebsiteURL,

		HelpMessage: req.HelpMessage,
	}

	// =====================================
	// SAVE
	// =====================================

	err = c.service.CreateContact(
		contact,
	)

	if err != nil {

		log.Printf(
			"failed to create contact: %v",
			err,
		)

		http.Error(
			w,
			"failed to save contact",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================
	// RESPONSE
	// =====================================

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(
		http.StatusCreated,
	)

	json.NewEncoder(w).Encode(
		map[string]any{
			"success": true,
			"message": "contact created successfully",
		},
	)
}
