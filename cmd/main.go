package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"lets-go-chat-v2/internal/config"
	"lets-go-chat-v2/internal/middleware"
	"lets-go-chat-v2/internal/user"
	"lets-go-chat-v2/internal/user/db/postgres"
	"lets-go-chat-v2/pkg/client/postgresql"
	"lets-go-chat-v2/pkg/logging"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("Start app")

	cfg := config.GetConfig()
	//
	postgreSQLClient, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("%v", err)
	}

	//repository := postgres.NewUserRepo(postgreSQLClient, logger)

	e := echo.New()
	handler := user.NewHandler()
	e.Use(middleware.ErrorMiddleware)
	e.Use()
	handler.Register(e)
	logger.Fatal(e.Start(":8070"))
	return
}
