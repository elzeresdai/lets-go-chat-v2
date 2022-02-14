package users

import (
	"context"
)

//go:generate mockgen -source=storage.go -destination=../tests/mocks/mock.go

type RepositoryInterface interface {
	CreateUser(ctx context.Context, user *CreateUserRequest) (*User, error)
	GetUser(ctx context.Context, name string) ([]*User, bool, error)
}
