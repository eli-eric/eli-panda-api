package publicationsservice

import (
	m "panda/apigateway/middlewares"
	"panda/apigateway/shared"

	"github.com/labstack/echo/v4"
)

func MapPublicationsRoutes(e *echo.Echo, h IPublicationsHandlers, jwtMiddleware echo.MiddlewareFunc) {
	e.GET("/v1/publication/:uid", m.Authorization(h.GetPublication(), shared.ROLE_BASICS_VIEW), jwtMiddleware)

	e.GET("/v1/publications", m.Authorization(h.GetPublications(), shared.ROLE_BASICS_VIEW), jwtMiddleware)

	e.POST("/v1/publication", m.Authorization(h.CreatePublication(), shared.ROLE_BASICS_VIEW), jwtMiddleware)

	e.PUT("/v1/publication/:uid", m.Authorization(h.UpdatePublication(), shared.ROLE_BASICS_VIEW), jwtMiddleware)

	e.DELETE("/v1/publication/:uid", m.Authorization(h.DeletePublication(), shared.ROLE_BASICS_VIEW), jwtMiddleware)
}
