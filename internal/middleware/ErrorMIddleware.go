package middleware

import (
	"errors"
	"github.com/labstack/echo/v4"
	"lets-go-chat-v2/internal/customerrors"
	"net/http"
)

type AppHandler func(e echo.Context) error

func ErrorMiddleware(h AppHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		err := h(c)
		var newErr *customerrors.AppError
		if err != nil {
			if errors.As(err, &newErr) {
				if errors.Is(err, customerrors.ErrNotFound) {
					c.Response().WriteHeader(http.StatusNotFound)
					c.Response().Write(customerrors.ErrNotFound.Marshal())
					return nil
				}

				err = err.(*customerrors.AppError)
				c.Response().WriteHeader(http.StatusBadRequest)
				c.Response().Write(newErr.Marshal())
				return nil
			}

			c.Response().WriteHeader(http.StatusTeapot)
			c.Response().Write(customerrors.SystemError(err).Marshal())
		}
		return nil
	}
}
