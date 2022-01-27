package repository

import (
	"lets-go-chat-v2/internal/user"
)

type Repository interface {
	LoginUser(user user.LoginUserRequest) (*user.User, error)
	CreateUser(user user.CreateUserRequest) (*user.User, error)
}
