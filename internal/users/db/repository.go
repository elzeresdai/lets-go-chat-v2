package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"lets-go-chat-v2/internal/users"
	"lets-go-chat-v2/pkg/client/postgresql"
	"lets-go-chat-v2/pkg/hasher"
	"lets-go-chat-v2/pkg/logging"
	"strings"
)

type UserRepo struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewUserRepo(client postgresql.Client, logger *logging.Logger) users.RepositoryInterface {
	return &UserRepo{
		client: client,
		logger: logger,
	}
}

func (u UserRepo) LoginUser(ctx context.Context, user *users.LoginUserRequest) (*users.User, error) {

	panic("implement me")
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (u UserRepo) CreateUser(ctx context.Context, user *users.CreateUserRequest) (*users.User, error) {
	hash, err := hasher.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	query := `INSERT INTO users ( id , name , hash) VALUES ($1, $2, $3) RETURNING $1, $2`
	u.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(query)))
	var newResp users.User
	err = u.client.QueryRow(ctx, query, uuid.New(), user.UserName, hash).Scan(&newResp.UserName, &newResp.UserName)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		pgErr = err.(*pgconn.PgError)
		newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
		u.logger.Error(newErr)
		return nil, newErr
	}

	return &newResp, nil
}

func (u UserRepo) GetUser(ctx context.Context, name string) ([]*users.User, bool, error) {
	query := `SELECT id, name FROM public.users where name = $1`
	u.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(query)))

	row, err := u.client.Query(ctx, query, name)

	var exist bool
	if err != nil {
		u.logger.Error(err)
		return nil, exist, err
	}
	defer row.Close()
	var user users.User
	existUsers := make([]*users.User, 0)
	for row.Next() {
		err := row.Scan(&user.ID, &user.UserName)
		if err != nil {
			u.logger.Error(err)
			continue
		}
		existUsers = append(existUsers, &user)
	}
	if len(existUsers) > 0 {
		exist = true
	}
	return existUsers, exist, nil
}
