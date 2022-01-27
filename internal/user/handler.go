package user

import (
	"github.com/labstack/echo/v4"
	"lets-go-chat-v2/internal/handlers"
)

type handler struct {
}

func NewHandler() handlers.HandlerInterface {
	return &handler{}
}

func (h *handler) Register(e *echo.Echo) {
	e.POST("/user", h.CreateUser)
	e.POST("/user/login", h.LoginUser)
}

func (h *handler) CreateUser(ctx echo.Context) error {
	return nil
}

func (h *handler) LoginUser(ctx echo.Context) error {
	return nil
}
