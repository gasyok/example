package user

import (
	"context"

	"example/internal/domain"
)

type repo interface {
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context, list, offset int64) ([]*domain.User, error)
	Create(ctx context.Context, u *domain.User) error
	Update(ctx context.Context, u *domain.User) error
	Delete(ctx context.Context, id int64) error
}

type txm interface {
	DoWithOptions(ctx context.Context, opts domain.TxOptions, fn func(ctx context.Context) error) error
	Do(ctx context.Context, fn func(context.Context) error) error
}
