package models

import "time"

type Contact struct {
	ID int64 `json:"id"`

	FullName string `json:"full_name"`

	Email string `json:"email"`

	PhoneNumber string `json:"phone_number"`

	CompanyName string `json:"company_name"`

	WebsiteURL string `json:"website_url"`

	HelpMessage string `json:"help_message"`

	CreatedAt time.Time `json:"created_at"`
}
