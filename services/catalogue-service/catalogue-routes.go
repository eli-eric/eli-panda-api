package catalogueService

import (
	"panda/apigateway/helpers"

	"github.com/labstack/echo/v4"
)

func MapCatalogueRoutes(e *echo.Echo, h ICatalogueHandlers, jwtMiddleware echo.MiddlewareFunc) {
	// categories route/s - return categories by parent path
	e.GET("/v1/catalogue/categories/*", h.GetCataloguecategoriesByParentPath(), jwtMiddleware)
	e.GET("/v1/catalogue/categories", h.GetCataloguecategoriesByParentPath(), jwtMiddleware)
	e.GET("/v1/catalogue/category/:uid", h.GetCatalogueCategoryWithDetailsByUid(), jwtMiddleware)
	e.PUT("/v1/catalogue/category/:uid", h.UpdateCatalogueCategory(), jwtMiddleware)
	e.POST("/v1/catalogue/category", h.CreateCatalogueCategory(), jwtMiddleware)
	e.DELETE("/v1/catalogue/category/:uid", h.DeleteCatalogueCategory(), jwtMiddleware)
	e.POST("/v1/catalogue/category/:uid/copy", h.CopyCatalogueCategoryRecursive(), jwtMiddleware)

	//fake image get - get only eli logo img for now
	e.GET("/v1/catalogue/category/:uid/image", h.GetCatalogueCategoryImageByUid())
	e.GET("/v1/catalogue/item/:uid/image", h.GetCatalogueItemImage())

	// get all catalogue items with details
	e.GET("/v1/catalogue/items", RoleBasedHandler(helpers.ROLE_SYSTEMS_VIEW, h.GetCatalogueItems()), jwtMiddleware)

	//get on catalogue item with details by uid
	e.GET("/v1/catalogue/item/:uid", h.GetCatalogueItemWithDetailsByUid(), jwtMiddleware)
}

func RoleBasedHandler(roleToCheck string, handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		user := helpers.GetUserFromJWT(c)

		if user != nil {
			for _, role := range user.Roles {
				if role == roleToCheck {
					return handlerFunc(c)
				}
			}
		}

		return echo.ErrUnauthorized
	}
}
