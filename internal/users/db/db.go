package db

import (
	"lets-go-chat-v2/internal/users"
	"lets-go-chat-v2/pkg/client/postgresql"
	"lets-go-chat-v2/pkg/logging"
)

type DB struct {
	client postgresql.Client
	logger *logging.Logger
	repo   users.RepositoryInterface
}
