package repositories

import (
	"context"
	"database/sql"
	"errors"

	"qonevo-backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(
	ctx context.Context,
	user *models.User,
) error {

	query := `
		INSERT INTO users (
			first_name,
			last_name,
			email,
			phone,
			password_hash
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		user.PasswordHash,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
}

func (r *UserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {

	query := `
		SELECT
			id,
			first_name,
			last_name,
			email,
			phone,
			password_hash,
			created_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	return user, nil
}