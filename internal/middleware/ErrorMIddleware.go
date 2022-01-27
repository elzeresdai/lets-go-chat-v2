package middleware

import (
	"errors"
	"github.com/labstack/echo/v4"
	"lets-go-chat-v2/internal/customerrors"
)

type AppHandler func(e *echo.HandlerFunc) error

func ErrorMiddleware(h AppHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		var newErr *customerrors.AppError
		err :=
		//if err != nil {
		//	if errors.As(err, &newErr) {
		//		if errors.Is(err, customerrors.ErrNotFound) {
		//			w.WriteHeader(http.StatusNotFound)
		//			w.Write(customerrors.ErrNotFound.Marshal())
		//			return
		//		}
		//
		//		err = err.(*customerrors.AppError)
		//		w.WriteHeader(http.StatusBadRequest)
		//		w.Write(newErr.Marshal())
		//		return
		//	}
		//
		//	w.WriteHeader(http.StatusTeapot)
		//	w.Write(customerrors.SystemError(err).Marshal())
		//}
	}
}
// func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		start := time.Now()
//		res := next(c)
//		log.WithFields(log.Fields{
//			"method":     c.Request().Method,
//			"path":       c.Path(),
//			"status":     c.Response().Status,
//			"latency_ns": time.Since(start).Nanoseconds(),
//		}).Info("request details")
//
//		return res
//	}
//}