package services

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"qonevo-backend/internal/models"
	"qonevo-backend/internal/repositories"
	"qonevo-backend/internal/utils"
)

type AuthService struct {
	repo   *repositories.UserRepo
	secret string
	expiry int
}

func NewAuthService(r *repositories.UserRepo, secret string, expiry int) *AuthService {
	return &AuthService{repo: r, secret: secret, expiry: expiry}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", errors.New("invalid credentials")
	}

	return utils.GenerateToken(user.ID, s.secret, s.expiry)
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*models.User, error) {
	// check if user exists (optional improvement later)

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: string(hash),
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}