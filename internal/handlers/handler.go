package handlers

import "github.com/labstack/echo/v4"

type HandlerInterface interface {
	Register (e *echo.Echo)
}
