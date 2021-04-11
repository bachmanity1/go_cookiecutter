package middleware

import (
	"pandita/util"
	"time"

	"github.com/labstack/echo/v4"
)

// ZapLogger ...
func ZapLogger(log *util.MLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}

			fields := []interface{}{
				"id", id,
				"status", res.Status,
				"latency", time.Since(start).String(),
				"method", req.Method,
				"uri", req.RequestURI,
				"host", req.Host,
				"remote_ip", c.RealIP(),
			}

			n := res.Status
			switch {
			case n >= 500:
				log.Errorw("Server error", fields...)
			case n == 404 && req.RequestURI == "/favicon.ico":
				break
			case n >= 400:
				log.Infow("Client error", fields...)
			case n >= 300:
				log.Infow("Redirection", fields...)
			case n == 208:
				break
			default:
				log.Infow("Success", fields...)
			}

			return nil
		}
	}
}
