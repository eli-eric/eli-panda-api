package cronservice

import (
	m "panda/apigateway/middlewares"
	"panda/apigateway/shared"

	"github.com/labstack/echo/v4"
)

func MapCronRoutes(e *echo.Echo, h ICronHandlers, jwtMiddleware echo.MiddlewareFunc) {

	e.GET("/v1/cron/history", m.Authorization(h.GetCronJobHistory(), shared.ROLE_ADMIN), jwtMiddleware)
}
