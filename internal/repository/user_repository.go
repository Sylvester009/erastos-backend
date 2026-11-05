package repository

import (
	"context"
	"time"
	"erastos-backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)


func CreateUser(pool *pgxpool.Pool, user *models.User) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, username, email, created_at, updated_at
	`

	err := pool.QueryRow(ctx, query, user.Username, user.Email, user.Password).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByEmail(pool *pgxpool.Pool, email string) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User

	err := pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(pool *pgxpool.Pool, id string) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User

	err := pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}