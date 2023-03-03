package systemsService

import (
	"github.com/labstack/echo/v4"
)

func MapSystemsRoutes(e *echo.Echo, h ISystemsHandlers, jwtMiddleware echo.MiddlewareFunc) {
	// get all subsystems for given parent system spec. by parentUID
	// if no parentUID is presented then get all root systems
	e.GET("/v1/system/subsystems/:parentUID", h.GetSubSystemsByParentUID(), jwtMiddleware)
	e.GET("/v1/system/subsystems/", h.GetSubSystemsByParentUID(), jwtMiddleware)
	e.GET("/v1/system/subsystems", h.GetSubSystemsByParentUID(), jwtMiddleware)

}
