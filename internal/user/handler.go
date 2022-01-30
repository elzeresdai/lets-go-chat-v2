package user

import (
	"github.com/labstack/echo/v4"
	"lets-go-chat-v2/internal/handlers"
	"lets-go-chat-v2/internal/middleware"
	repository "lets-go-chat-v2/internal/user/db"
	"lets-go-chat-v2/pkg/logging"
)

type handler struct {
	logger     *logging.Logger
	repository repository.Repository
}

func NewHandler(repository repository.Repository, logger *logging.Logger) handlers.HandlerInterface {
	return &handler{
		repository: repository,
		logger:     logger,
	}
}

func (h *handler) Register(e *echo.Echo) {
	e.POST("/user", middleware.ErrorMiddleware(h.CreateUser))
	e.POST("/user/login", middleware.ErrorMiddleware(h.LoginUser))
}

func (h *handler) CreateUser(ctx echo.Context) error {
	return nil
}

func (h *handler) LoginUser(ctx echo.Context) error {
	return nil
}
