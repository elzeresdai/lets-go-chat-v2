package handlers

import "github.com/labstack/echo/v4"

//go:generate mockgen -source=handler.go -destination=mocks/mock.go

type HandlerInterface interface {
	Register(e *echo.Echo)
}
