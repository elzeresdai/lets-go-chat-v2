package middleware

import (
	"context"
	"github.com/labstack/echo/v4"
	"lets-go-chat-v2/internal/auth"
	"lets-go-chat-v2/internal/customerrors"
	"net/http"
)

type contextKey string

const UserContextKey = contextKey("user")

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, tok := c.Request().URL.Query()["token"]

		if tok && len(token) == 1 {

			user, err := auth.ValidateToken(token[0])
			if err != nil {
				err := customerrors.NewAppError(
					nil,
					"Forbidden",
					"",
					"403",
				)
				c.Response().WriteHeader(http.StatusForbidden)
				c.Response().Write(err.Marshal())
				return nil

			} else {
				c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), UserContextKey, &user)))
			}

		} else {
			err := customerrors.NewAppError(
				nil,
				"Invalid login",
				"",
				"400",
			)
			c.Response().WriteHeader(http.StatusNotFound)
			c.Response().Write(err.Marshal())
			return nil
		}
		return next(c)
	}

}
