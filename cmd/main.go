package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	middleware2 "github.com/labstack/echo/v4/middleware"
	"lets-go-chat-v2/internal/config"
	db2 "lets-go-chat-v2/internal/messages/db"
	"lets-go-chat-v2/internal/middleware"
	"lets-go-chat-v2/internal/users"
	"lets-go-chat-v2/internal/users/db"
	"lets-go-chat-v2/pkg/client/postgresql"
	"lets-go-chat-v2/pkg/logging"
	"lets-go-chat-v2/pkg/websocket"
	"log"
	"os"
)

type ContextValue struct {
	echo.Context
}

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
	messageRepo := db2.NewMessageRepo(postgreSQLClient, logger)
	logger.Info("register users handler")

	e := echo.New()
	e.Use(middleware2.Recover())
	e.Use(middleware.LoggingMiddleware)

	handler := users.NewHandler(repository, logger)
	handler.Register(e)

	hub := websocket.NewHub(repository, messageRepo)
	go hub.Run()

	e.GET("/chat/ws.rtm.start", middleware.AuthMiddleware(func(c echo.Context) error {
		log.Println(c.Request().Context())
		websocket.ServeWs(hub, c.Response(), c.Request())
		return nil
	}))
	logger.Fatal(e.Start(":" + port))
	return
}
