package users

import (
	"context"
)

type RepositoryInterface interface {
	LoginUser(ctx context.Context, user *LoginUserRequest) (*User, error)
	CreateUser(ctx context.Context, user *CreateUserRequest) (*User, error)
	GetUser(ctx context.Context, name string) ([]*User, bool, error)
}
