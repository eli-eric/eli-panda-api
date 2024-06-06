package filesservice

import (
	m "panda/apigateway/middlewares"
	"panda/apigateway/shared"

	"github.com/labstack/echo/v4"
)

// MapFilesRoutes maps files routes
func MapFilesRoutes(e *echo.Echo, h IFilesHandlers, jwtMiddleware echo.MiddlewareFunc) {
	e.GET("/v1/files/links/:parentUid", m.Authorization(h.GetFileLinksByParentUid(), shared.ROLE_BASICS_VIEW), jwtMiddleware)
	e.POST("/v1/files/link/:parentUid", m.Authorization(h.CreateFileLink(), shared.ROLE_SYSTEMS_EDIT, shared.ROLE_CATALOGUE_EDIT, shared.ROLE_ORDERS_EDIT, shared.ROLE_ROOM_CARDS_EDIT), jwtMiddleware)
	e.PUT("/v1/files/link/:uid", m.Authorization(h.UpdateFileLink(), shared.ROLE_SYSTEMS_EDIT, shared.ROLE_CATALOGUE_EDIT, shared.ROLE_ORDERS_EDIT, shared.ROLE_ROOM_CARDS_EDIT), jwtMiddleware)
	e.DELETE("/v1/files/link/:uid", m.Authorization(h.DeleteFileLink(), shared.ROLE_SYSTEMS_EDIT, shared.ROLE_CATALOGUE_EDIT, shared.ROLE_ORDERS_EDIT, shared.ROLE_ROOM_CARDS_EDIT), jwtMiddleware)
	e.POST("/v1/files/node/:uid/mini-image-url", m.Authorization(h.SetMiniImageUrlToNode(), shared.ROLE_SYSTEMS_EDIT, shared.ROLE_CATALOGUE_EDIT, shared.ROLE_CATALOGUE_CATEGORY_EDIT), jwtMiddleware)
}
