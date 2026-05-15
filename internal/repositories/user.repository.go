package repositories

import (
	"context"
	"database/sql"

	"qonevo-backend/internal/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}

	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, created_at FROM users WHERE email=$1`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)

	return u, err
}

func (r *UserRepo) Create(ctx context.Context, u *models.User) error {
	return r.db.QueryRowContext(
		ctx,
		`INSERT INTO users (email, password_hash)
		 VALUES ($1, $2)
		 RETURNING id, created_at`,
		u.Email,
		u.PasswordHash,
	).Scan(&u.ID, &u.CreatedAt)
}