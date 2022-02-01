package main

import (
	"context"
	"github.com/labstack/echo/v4"
	middleware2 "github.com/labstack/echo/v4/middleware"
	"lets-go-chat-v2/internal/config"
	"lets-go-chat-v2/internal/middleware"
	"lets-go-chat-v2/internal/users"
	"lets-go-chat-v2/internal/users/db"
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

	repository := db.NewUserRepo(postgreSQLClient, logger)
	logger.Info("register users handler")

	e := echo.New()
	e.Use(middleware2.Recover())
	e.Use(middleware.LoggingMiddleware)

	handler := users.NewHandler(repository, logger)
	handler.Register(e)
	logger.Fatal(e.Start(":8099"))
	return
}
