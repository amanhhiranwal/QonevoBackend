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

// 🔍 Find user by email (FIXED)
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}

	err := r.db.QueryRowContext(ctx,
		`SELECT 
			id,
			first_name,
			last_name,
			email,
			phone,
			password_hash,
			created_at
		 FROM users 
		 WHERE email = $1`,
		email,
	).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Phone,
		&u.PasswordHash,
		&u.CreatedAt,
	)

	// ✅ IMPORTANT: user not found is NOT an error
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return u, nil
}

// 🧾 Create user (FIXED & consistent with schema)
func (r *UserRepo) Create(ctx context.Context, u *models.User) error {
	return r.db.QueryRowContext(
		ctx,
		`INSERT INTO users (
			first_name,
			last_name,
			email,
			phone,
			password_hash
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Phone,
		u.PasswordHash,
	).Scan(
		&u.ID,
		&u.CreatedAt,
	)
}