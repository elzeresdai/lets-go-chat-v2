package middleware

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		req := c.Request()
		resp := c.Response()
		stop := time.Now()
		p := req.URL.Path
		bytesIn := req.Header.Get(echo.HeaderContentLength)
		res := next(c)
		log.WithFields(log.Fields{
			"time_rfc3339":  time.Now().Format(time.RFC3339),
			"remote_ip":     c.RealIP(),
			"host":          req.Host,
			"uri":           req.RequestURI,
			"method":        req.Method,
			"path":          p,
			"referer":       req.Referer(),
			"user_agent":    req.UserAgent(),
			"status":        resp.Status,
			"latency":       strconv.FormatInt(stop.Sub(start).Nanoseconds()/1000, 10),
			"latency_human": stop.Sub(start).String(),
			"bytes_in":      bytesIn,
			"bytes_out":     strconv.FormatInt(resp.Size, 10),
		}).Info("Handled request")

		return res
	}
}
