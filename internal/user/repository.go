package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, email, password_hash, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(ctx, query, user.ID, user.Email, user.PasswordHash, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE email = $1
	`
	var user User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE id = $1
	`
	var user User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil
}
