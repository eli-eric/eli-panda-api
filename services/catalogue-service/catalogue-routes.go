package catalogueService

import (
	"github.com/labstack/echo/v4"
)

func MapCatalogueRoutes(e *echo.Echo, h ICatalogueHandlers, jwtMiddleware echo.MiddlewareFunc) {
	// categories route
	e.GET("/v1/catalogue/categories/*", h.GetCataloguecategoriesByParentPath(), jwtMiddleware)

}
