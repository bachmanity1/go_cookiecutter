package middleware

import (
	"context"
	"pandita/util"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// TransID Transaction ID(yyyymmddhhmi + 5 numbers)
func TransID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, util.TransIDKey, util.NewID())
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

// JWTInit ...
func JWTInit(handler middleware.JWTErrorHandlerWithContext, key string) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:              []byte(key),
		ErrorHandlerWithContext: handler,
	})
}
