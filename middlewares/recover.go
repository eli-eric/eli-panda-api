package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RecoverMiddleware() echo.MiddlewareFunc {

	return middleware.Recover()
}
