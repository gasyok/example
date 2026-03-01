package user

import (
	"context"
	"errors"
	"fmt"

	"example/internal/domain"
)

type CreateInput struct {
	Name  string
	Email string
}

type UpdateInput struct {
	Name  string
	Email string
}

type Service struct {
	repo repo
	txm  txm
}

func New(repo repo, txm txm) *Service {
	return &Service{repo: repo, txm: txm}
}

func (s *Service) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, limit, offset int64) ([]*domain.User, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *Service) Create(ctx context.Context, in CreateInput) (*domain.User, error) {
	u := &domain.User{
		Name:  in.Name,
		Email: in.Email,
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) Update(ctx context.Context, id int64, in UpdateInput) (*domain.User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	u.Name = in.Name
	u.Email = in.Email

	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

// Можно также переносить это в рамках одного метода в repo, иногда делают так смещая
// эту логику в сервисы
func (s *Service) Upsert(ctx context.Context, in CreateInput) (*domain.User, error) {
	var u *domain.User
	err := s.txm.Do(ctx, func(ctx context.Context) error {
		existing, err := s.repo.GetByEmail(ctx, in.Email)
		if err != nil && !errors.Is(err, domain.ErrNotFound) {
			return fmt.Errorf("get by email: %w", err)
		}
		if existing != nil {
			existing.Name = in.Name
			if err := s.repo.Update(ctx, existing); err != nil {
				return fmt.Errorf("update: %w", err)
			}
			u = existing
		} else {
			u = &domain.User{Name: in.Name, Email: in.Email}
			if err := s.repo.Create(ctx, u); err != nil {
				return fmt.Errorf("create: %w", err)
			}
		}
		return nil
	})
	return u, err
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.txm.Do(ctx, func(ctx context.Context) error {
		if _, err := s.repo.GetByID(ctx, id); err != nil {
			return err
		}
		return s.repo.Delete(ctx, id)
	})
}
