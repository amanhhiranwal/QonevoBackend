package services

import (
	"context"
	"errors"

	"qonevo-backend/internal/config"
	"qonevo-backend/internal/models"
	"qonevo-backend/internal/repositories"
	"qonevo-backend/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repositories.UserRepository
	cfg  *config.Config
}

func NewAuthService(
	repo *repositories.UserRepository,
	cfg *config.Config,
) *AuthService {

	return &AuthService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	firstName string,
	lastName string,
	email string,
	phone *string,
	password string,
) error {

	existingUser, err := s.repo.FindByEmail(ctx, email)

	if err == nil && existingUser != nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	user := &models.User{
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		Phone:       phone,
		PasswordHash: string(hashedPassword),
	}

	return s.repo.Create(ctx, user)
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {

	user, err := s.repo.FindByEmail(ctx, email)

	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return "", errors.New("invalid credentials")
	}

	return utils.GenerateToken(
		user.ID,
		s.cfg.JWTSecret,
		s.cfg.JWTExpiryHours,
	)
}