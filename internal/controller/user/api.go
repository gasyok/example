package user

import (
	"example/internal/domain"
	"example/pkg/dto"
)

func toUserResponse(u *domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func toUserListResponse(users []*domain.User) []dto.UserResponse {
	result := make([]dto.UserResponse, len(users))
	for i, u := range users {
		result[i] = toUserResponse(u)
	}
	return result
}

func validateListRange(req *dto.ListRange) error {
	if req.Limit < 0 {
		return domain.ErrInvalidInput
	}
	if req.Offset < 0 {
		return domain.ErrInvalidInput
	}
	return nil
}
