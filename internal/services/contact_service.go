package services

import (
	"qonevo-backend/internal/models"
	"qonevo-backend/internal/repositories"
)

type ContactService struct {
	repo *repositories.ContactRepository
}

func NewContactService(
	repo *repositories.ContactRepository,
) *ContactService {

	return &ContactService{
		repo: repo,
	}
}

// =====================================
// CREATE CONTACT
// =====================================

func (s *ContactService) CreateContact(
	contact *models.Contact,
) error {

	return s.repo.CreateContact(
		contact,
	)
}
