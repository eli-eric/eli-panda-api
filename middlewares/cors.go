package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CORSMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,
			echo.HeaderAuthorization,
			"hx-current-url",
			"hx-request",
			"hx-target",
			"hx-trigger",
			"hx-history",
			"hx-swap",
			"hx-select",
			"hx-indicator",
			"hx-boost",
			"hx-prompt",
			"hx-push-url",
			"hx-swap-oob",
			"hx-boost-oob",
			"hx-prompt-oob",
			"hx-prompt-override",
			"hx-boost-override",
			"hx-swap-override",
			"hx-trigger-override",
			"hx-history-override",
		},
		AllowCredentials: true,
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	})
}
