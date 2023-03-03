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

	//get system image - base64string
	e.GET("/v1/system/:uid/image", h.GetSystemImageByUid(), jwtMiddleware)

	// get system detail by uid
	e.GET("/v1/system/:uid", h.GetSystemDetail(), jwtMiddleware)
}
