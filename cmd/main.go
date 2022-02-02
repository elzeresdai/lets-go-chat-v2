package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	middleware2 "github.com/labstack/echo/v4/middleware"
	"lets-go-chat-v2/internal/config"
	"lets-go-chat-v2/internal/middleware"
	"lets-go-chat-v2/internal/users"
	"lets-go-chat-v2/internal/users/db"
	websocket2 "lets-go-chat-v2/internal/users/websocket"
	"lets-go-chat-v2/pkg/client/postgresql"
	"lets-go-chat-v2/pkg/logging"
	"os"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("Start app")

	cfg := config.GetConfig()

	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatalf("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

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

	hub := websocket2.NewHub()
	go hub.Run()
	e.GET("/chat/ws.rtm.start/", func(c echo.Context) error {
		websocket2.ServeWs(hub, c.Response(), c.Request(), repository)
		return nil
	})
	logger.Fatal(e.Start(":" + port))
	return
}
