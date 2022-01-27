package postgres

import (
	"context"
	"github.com/google/uuid"
	"lets-go-chat-v2/internal/user"
	"lets-go-chat-v2/internal/user/db"
	"lets-go-chat-v2/pkg/client/postgresql"
	"lets-go-chat-v2/pkg/hasher"
	"lets-go-chat-v2/pkg/logging"
	"strings"
)

type UserRepo struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewUserRepo(client postgresql.Client, logger *logging.Logger) repository.Repository {
	return &UserRepo{
		client: client,
		logger: logger,
	}
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (u UserRepo) LoginUser(ctx context.Context) (*user.User, error) {
	panic("implement me")
}

func (u UserRepo) CreateUser(us user.CreateUserRequest) (*user.User, error) {
	hash, err := hasher.HashPassword(us.Password)
	if err != nil {
		return nil, err
	}
	query := `INSERT INTO users ( id , name , hash) VALUES ($1, $2, $3) RETURNING $1, $2 `

	var newResp user.User

	//err = u.client.QueryRow(ctx,query, uuid.New(), user.UserName, hash).Scan(&newResp.ID, &newResp.Name)
	//if err != nil {
	//	log.Printf("this was the error: %v", err.Error())
	//	return nil, err
	//}
	//return &newResp, nil

}
