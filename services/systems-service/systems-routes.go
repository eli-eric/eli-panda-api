package systemsService

import (
	"github.com/labstack/echo/v4"

	m "panda/apigateway/middlewares"
	"panda/apigateway/shared"
)

func MapSystemsRoutes(e *echo.Echo, h ISystemsHandlers, jwtMiddleware echo.MiddlewareFunc) {

	// get systems with search and pagination
	e.GET("/v1/systems", m.Authorization(h.GetSystemsWithSearchAndPagination(), shared.ROLE_SYSTEMS_VIEW), jwtMiddleware)
	// get all subsystems for given parent system spec. by parentUID
	// if no parentUID is presented then get all root systems
	e.GET("/v1/system/:parentUID/subsystems", m.Authorization(h.GetSubSystemsByParentUID(), shared.ROLE_SYSTEMS_VIEW), jwtMiddleware)

	//get system image - base64string
	e.GET("/v1/system/:uid/image", m.Authorization(h.GetSystemImageByUid(), shared.ROLE_SYSTEMS_VIEW), jwtMiddleware)

	e.POST("/v1/system", m.Authorization(h.CreateNewSystem(), shared.ROLE_SYSTEMS_EDIT), jwtMiddleware)
	// get system detail by uid
	e.GET("/v1/system/:uid", m.Authorization(h.GetSystemDetail(), shared.ROLE_SYSTEMS_VIEW), jwtMiddleware)
	//save new system/sub-system
	e.PUT("/v1/system/:uid", m.Authorization(h.UpdateSystem(), shared.ROLE_SYSTEMS_EDIT), jwtMiddleware)
	e.DELETE("/v1/system/:uid", m.Authorization(h.DeleteSystemRecursive(), shared.ROLE_SYSTEMS_EDIT), jwtMiddleware)
	// this one is only becasue of bad request from ui for now
	e.POST("/v1/system/:xxx", m.Authorization(h.CreateNewSystem(), shared.ROLE_SYSTEMS_EDIT), jwtMiddleware)

	// get systems for relationship
	e.GET("/v1/systems/for-relationship", m.Authorization(h.GetSystemsForRelationship(), shared.ROLE_SYSTEMS_VIEW), jwtMiddleware)

	e.GET("/v1/system/:parentUID/subsystems/for-relationship", m.Authorization(h.GetSubSystemsByParentUID(), shared.ROLE_SYSTEMS_VIEW), jwtMiddleware)

	// get system relationships
	e.GET("/v1/system/:uid/relationships", m.Authorization(h.GetSystemRelationships(), shared.ROLE_SYSTEMS_VIEW), jwtMiddleware)

	// delete system relationship
	e.DELETE("/v1/system/relationship/:uid", m.Authorization(h.DeleteSystemRelationship(), shared.ROLE_SYSTEMS_EDIT), jwtMiddleware)

	// create new system relationship
	e.POST("/v1/system/relationship/:uid", m.Authorization(h.CreateNewSystemRelationship(), shared.ROLE_SYSTEMS_EDIT), jwtMiddleware)

	e.GET("/v1/system/systemCode", m.Authorization(h.GetSystemCode(), shared.ROLE_SYSTEMS_VIEW), jwtMiddleware)

	e.GET("/v1/physical-item/:uid/properties", m.Authorization(h.GetPhysicalItemProperties(), shared.ROLE_SYSTEMS_VIEW, shared.ROLE_CATALOGUE_VIEW, shared.ROLE_ORDERS_VIEW), jwtMiddleware)
}
