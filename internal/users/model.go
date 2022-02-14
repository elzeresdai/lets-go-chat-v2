package users

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	UserName     string `json:"username"`
	PasswordHash string `json:"-"`
}

type CreateUserRequest struct {
	UserName string `json:"username" validate:"required,min=4"`
	Password string `json:"password" validate:"required,min=8"`
}

// swagger:model LoginUserRequest
type LoginUserRequest struct {
	UserName string `json:"userName" validate:"required,min=4"`
	Password string `json:"password" validate:"required,min=8"`
}

// swagger:model CreateUserResponse
type CreateUserResponse struct {
	Id       uuid.UUID
	UserName string
}

// swagger:model LoginUserResponse
type LoginUserResponse struct {
	Url string
}
type ValidationResponse struct {
	Message validator.FieldError
}

type ActiveUsersResponse struct {
	count int
}

type ActiveUsers struct {
	UserName, Token string
}
