package messages

import (
	"context"
	"github.com/google/uuid"
)

type MessageRepositoryInterface interface {
	GetUnreadMessages(ctx context.Context) ([]*Message, error)
	CreateUserMessage(ctx context.Context, userId uuid.UUID, message string) (*Message, error)
}
