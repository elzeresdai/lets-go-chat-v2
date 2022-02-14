package messages

import (
	"context"
	"github.com/google/uuid"
)

//go:generate mockgen -source=storage.go -destination=mocks/mock.go

type MessageRepositoryInterface interface {
	GetUnreadMessages(ctx context.Context) ([]*Message, error)
	CreateUserMessage(ctx context.Context, userId uuid.UUID, message string) (*Message, error)
}
