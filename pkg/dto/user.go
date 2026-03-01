package dto

import "time"

type UserResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpsertUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ListData struct {
	Items interface{} `json:"items"`
}

type ListRange struct {
	Count  int64 `json:"count"`
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}
