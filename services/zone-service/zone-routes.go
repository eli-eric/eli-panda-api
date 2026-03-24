package zoneservice

import (
	m "panda/apigateway/middlewares"
	"panda/apigateway/shared"

	"github.com/labstack/echo/v4"
)

func MapZoneRoutes(e *echo.Echo, h IZoneHandlers, jwtMiddleware echo.MiddlewareFunc) {
	e.GET("/v1/zones", m.Authorization(h.GetAllZones(), shared.ROLE_ZONES_VIEW), jwtMiddleware)
	e.GET("/v1/zones/:uid", m.Authorization(h.GetZoneByUID(), shared.ROLE_ZONES_VIEW), jwtMiddleware)
	e.POST("/v1/zones", m.Authorization(h.CreateZone(), shared.ROLE_ZONES_EDIT), jwtMiddleware)
	e.PUT("/v1/zones/:uid", m.Authorization(h.UpdateZone(), shared.ROLE_ZONES_EDIT), jwtMiddleware)
	e.DELETE("/v1/zones/:uid", m.Authorization(h.DeleteZone(), shared.ROLE_ZONES_EDIT), jwtMiddleware)
	e.POST("/v1/zones/import", m.Authorization(h.ImportZones(), shared.ROLE_ZONES_EDIT), jwtMiddleware)
}
