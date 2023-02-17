package catalogueService

import (
	"github.com/labstack/echo/v4"
)

func MapCatalogueRoutes(e *echo.Echo, h ICatalogueHandlers, jwtMiddleware echo.MiddlewareFunc) {
	// categories route/s - return categories by parent path
	e.GET("/v1/catalogue/categories/*", h.GetCataloguecategoriesByParentPath(), jwtMiddleware)
	e.GET("/v1/catalogue/categories", h.GetCataloguecategoriesByParentPath(), jwtMiddleware)

	//fake image get - get only eli logo img for now
	e.GET("/v1/catalogue/category/:uid/image", h.GetCatalogueCategoryImage())

	e.GET("/v1/catalogue/items", h.GetCatalogueItems(), jwtMiddleware)
}
