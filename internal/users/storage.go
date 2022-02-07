package users

import (
	"context"
)

type RepositoryInterface interface {
	CreateUser(ctx context.Context, user *CreateUserRequest) (*User, error)
	GetUser(ctx context.Context, name string) ([]*User, bool, error)
}
