package publicationsservice

import (
	m "panda/apigateway/middlewares"
	"panda/apigateway/shared"

	"github.com/labstack/echo/v4"
)

func MapPublicationsRoutes(e *echo.Echo, h IPublicationsHandlers, jwtMiddleware echo.MiddlewareFunc) {
	e.GET("/v1/publication/:uid", m.Authorization(h.GetPublication(), shared.ROLE_PUBLICATIONS_VIEW), jwtMiddleware)

	e.GET("/v1/publications", m.Authorization(h.GetPublications(), shared.ROLE_PUBLICATIONS_VIEW), jwtMiddleware)

	e.GET("/v1/publications/export", m.Authorization(h.GetPublicationsAsCsv(), shared.ROLE_PUBLICATIONS_VIEW), jwtMiddleware)

	e.POST("/v1/publication", m.Authorization(h.CreatePublication(), shared.ROLE_PUBLICATIONS_EDIT), jwtMiddleware)

	e.PUT("/v1/publication/:uid", m.Authorization(h.UpdatePublication(), shared.ROLE_PUBLICATIONS_EDIT), jwtMiddleware)

	e.DELETE("/v1/publication/:uid", m.Authorization(h.DeletePublication(), shared.ROLE_PUBLICATIONS_EDIT), jwtMiddleware)

	e.GET("/v1/publication/wos/:doi", m.Authorization(h.GetWosDataByDoi(), shared.ROLE_PUBLICATIONS_VIEW), jwtMiddleware)

	// Researchers CRUD
	e.GET("/v1/researchers", m.Authorization(h.GetResearchers(), shared.ROLE_PUBLICATIONS_VIEW), jwtMiddleware)

	e.GET("/v1/researcher/:uid", m.Authorization(h.GetResearcher(), shared.ROLE_PUBLICATIONS_VIEW), jwtMiddleware)

	e.POST("/v1/researcher", m.Authorization(h.CreateResearcher(), shared.ROLE_PUBLICATIONS_EDIT), jwtMiddleware)

	e.PUT("/v1/researcher/:uid", m.Authorization(h.UpdateResearcher(), shared.ROLE_PUBLICATIONS_EDIT), jwtMiddleware)

	e.DELETE("/v1/researcher/:uid", m.Authorization(h.DeleteResearcher(), shared.ROLE_PUBLICATIONS_EDIT), jwtMiddleware)
}
