package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"github.com/SuperIntelligence-Labs/go-backend-template/internal/logger"
)

func Zerolog() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			latency := time.Since(start)

			req := c.Request()
			res := c.Response()
			status := res.Status

			var event *zerolog.Event

			switch {
			case status >= 500:
				event = logger.Error()
			case status >= 400:
				event = logger.Warn()
			default:
				event = logger.Info()
			}

			event.
				Str("remote_ip", c.RealIP()).
				Str("method", req.Method).
				Str("path", req.URL.Path).
				Str("host", req.Host).
				Int("status", status).
				Int64("bytes_in", req.ContentLength).
				Int64("bytes_out", res.Size).
				Dur("latency", latency).
				Str("latency_human", latency.String()).
				Str("user_agent", req.UserAgent())

			if err != nil {
				event = event.Err(err)
			}

			if id := res.Header().Get(echo.HeaderXRequestID); id != "" {
				event = event.Str("request_id", id)
			}

			event.Msg("HTTP request")
			return nil
		}
	}
}
