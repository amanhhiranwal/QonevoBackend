package repositories

import (
	"database/sql"

	"qonevo-backend/internal/models"
)

type ContactRepository struct {
	db *sql.DB
}

func NewContactRepository(
	db *sql.DB,
) *ContactRepository {

	return &ContactRepository{
		db: db,
	}
}

// =====================================
// CREATE CONTACT
// =====================================

func (r *ContactRepository) CreateContact(
	contact *models.Contact,
) error {

	query := `
	INSERT INTO contacts (
		full_name,
		email,
		phone_number,
		company_name,
		website_url,
		help_message
	)
	VALUES ($1,$2,$3,$4,$5,$6)
	`

	_, err := r.db.Exec(
		query,
		contact.FullName,
		contact.Email,
		contact.PhoneNumber,
		contact.CompanyName,
		contact.WebsiteURL,
		contact.HelpMessage,
	)

	return err
}
