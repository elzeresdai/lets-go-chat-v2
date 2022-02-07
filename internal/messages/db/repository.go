package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"lets-go-chat-v2/internal/messages"
	"lets-go-chat-v2/pkg/client/postgresql"
	"lets-go-chat-v2/pkg/logging"
	"strings"
)

type MessageRepo struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewMessageRepo(client postgresql.Client, logger *logging.Logger) messages.MessageRepositoryInterface {
	return &MessageRepo{
		client: client,
		logger: logger,
	}
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (m MessageRepo) GetUnreadMessages(ctx context.Context) ([]*messages.Message, error) {
	query := `SELECT user_id, message FROM user_messages ORDER BY created_at desc`
	m.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(query)))

	row, err := m.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	userMessages := make([]*messages.Message, 0)
	for row.Next() {
		userMessage := messages.Message{}
		err := row.Scan(&userMessage.UserId, &userMessage.Message)
		if err != nil {
			m.logger.Println(err)
			continue
		}
		userMessages = append(userMessages, &userMessage)
	}
	return userMessages, nil
}

func (m MessageRepo) CreateUserMessage(ctx context.Context, userId uuid.UUID, message string) (*messages.Message, error) {
	query := `INSERT INTO user_messages (user_id, message) VALUES ($1, $2) RETURNING $1, $2`
	m.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(query)))
	var newResp messages.Message
	err := m.client.QueryRow(ctx, query, userId, message).Scan(&newResp.UserId, &newResp.Message)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		pgErr = err.(*pgconn.PgError)
		newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
		m.logger.Error(newErr)
		return nil, newErr
	}

	return &newResp, nil
}
