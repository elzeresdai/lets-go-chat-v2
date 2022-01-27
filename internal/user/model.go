package user

import (
	"github.com/gofrs/uuid"
	"lets-go-chat-v2/internal/messages"
)

type User struct {
	ID           uuid.UUID
	UserName     string           `json:"username"`
	PasswordHash string           `json:"-"`
	Messages     messages.Message `json:"messages"`
}

type CreateUserRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
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
	Message string
}
