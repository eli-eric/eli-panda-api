package general

import (
	m "panda/apigateway/middlewares"
	"panda/apigateway/shared"

	"github.com/labstack/echo/v4"
)

func MapGeneralRoutes(e *echo.Echo, h IGeneralHandlers, jwtMiddleware echo.MiddlewareFunc) {
	e.GET("/v1/general/:uid/graph", m.Authorization(h.GetGraphByUid(), shared.ROLE_BASICS_VIEW), jwtMiddleware)
}
