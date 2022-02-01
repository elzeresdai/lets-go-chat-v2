package main

import (
	"context"
	"github.com/labstack/echo/v4"
	middleware2 "github.com/labstack/echo/v4/middleware"
	"lets-go-chat-v2/internal/config"
	"lets-go-chat-v2/internal/middleware"
	"lets-go-chat-v2/internal/users"
	"lets-go-chat-v2/internal/users/db"
	cache2 "lets-go-chat-v2/pkg/cache"
	"lets-go-chat-v2/pkg/client/postgresql"
	"lets-go-chat-v2/pkg/logging"
	"lets-go-chat-v2/pkg/websocket"
	"time"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("Start app")

	cfg := config.GetConfig()

	postgreSQLClient, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	cache2.NewLocalCache(10 * time.Minute)
	repository := db.NewUserRepo(postgreSQLClient, logger)
	logger.Info("register users handler")

	e := echo.New()
	e.Use(middleware2.Recover())
	e.Use(middleware.LoggingMiddleware)

	handler := users.NewHandler(repository, logger)
	handler.Register(e)

	hub := websocket.NewHub()
	go hub.Run()
	e.GET("/chat/ws.rtm.start/{token}", func(c echo.Context) error {
		websocket.ServeWs(hub, c.Response(), c.Request())
		return nil
	})
	logger.Fatal(e.Start(":8088"))
	return
}
