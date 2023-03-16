package systemsService

import (
	"github.com/labstack/echo/v4"

	m "panda/apigateway/middlewares"
	"panda/apigateway/shared"
)

func MapSystemsRoutes(e *echo.Echo, h ISystemsHandlers, jwtMiddleware echo.MiddlewareFunc) {
	// get all subsystems for given parent system spec. by parentUID
	// if no parentUID is presented then get all root systems
	e.GET("/v1/system/subsystems/:parentUID", m.Authorization(h.GetSubSystemsByParentUID(), shared.ROLE_SYSTEMS_VIEW), jwtMiddleware)
	e.GET("/v1/system/subsystems/", m.Authorization(h.GetSubSystemsByParentUID(), shared.ROLE_SYSTEMS_VIEW), jwtMiddleware)
	e.GET("/v1/system/subsystems", m.Authorization(h.GetSubSystemsByParentUID(), shared.ROLE_SYSTEMS_VIEW), jwtMiddleware)

	//get system image - base64string
	e.GET("/v1/system/:uid/image", m.Authorization(h.GetSystemImageByUid(), shared.ROLE_SYSTEMS_VIEW), jwtMiddleware)

	// get system detail by uid
	e.GET("/v1/system/:uid", m.Authorization(h.GetSystemDetail(), shared.ROLE_SYSTEMS_VIEW), jwtMiddleware)

	//save new system/sub-system
	e.POST("/v1/system", m.Authorization(h.CreateNewSystem(), shared.ROLE_SYSTEMS_EDIT), jwtMiddleware)
	// this one is only becasue of bad request from ui for now
	e.POST("/v1/system/:xxx", m.Authorization(h.CreateNewSystem(), shared.ROLE_SYSTEMS_EDIT), jwtMiddleware)

	e.PUT("/v1/system/:uid", m.Authorization(h.UpdateSystem(), shared.ROLE_SYSTEMS_EDIT), jwtMiddleware)
}
