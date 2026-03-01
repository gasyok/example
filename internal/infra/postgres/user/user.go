package user

import (
	"context"
	"errors"
	"fmt"

	"example/internal/domain"
	"example/internal/infra/postgres"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	pool postgres.PgxPool
}

func New(pool postgres.PgxPool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	eng := postgres.QueryEngine(ctx, r.pool)

	var u domain.User
	err := eng.QueryRow(ctx,
		`SELECT id, name, email, created_at, updated_at
		 FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query user by id: %w", err)
	}
	return &u, nil
}

func (r *Repository) List(ctx context.Context, limit, offset int64) ([]*domain.User, error) {
	eng := postgres.QueryEngine(ctx, r.pool)

	rows, err := eng.Query(ctx,
		`SELECT id, name, email, created_at, updated_at FROM users ORDER BY id
		limit $1 offset $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, &u)
	}
	return users, rows.Err()
}

func (r *Repository) Create(ctx context.Context, u *domain.User) error {
	eng := postgres.QueryEngine(ctx, r.pool)

	err := eng.QueryRow(ctx,
		`INSERT INTO users (name, email)
		 VALUES ($1, $2)
		 RETURNING id, created_at, updated_at`,
		u.Name, u.Email,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, u *domain.User) error {
	eng := postgres.QueryEngine(ctx, r.pool)

	_, err := eng.Exec(ctx,
		`UPDATE users SET name=$1, email=$2, updated_at=now() WHERE id=$3`,
		u.Name, u.Email, u.ID,
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	eng := postgres.QueryEngine(ctx, r.pool)

	var u domain.User
	err := eng.QueryRow(ctx,
		`SELECT id, name, email, created_at, updated_at
		 FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query user by email: %w", err)
	}
	return &u, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	eng := postgres.QueryEngine(ctx, r.pool)

	_, err := eng.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}
