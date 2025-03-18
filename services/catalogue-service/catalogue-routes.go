package catalogueService

import (
	m "panda/apigateway/middlewares"
	"panda/apigateway/shared"

	"github.com/labstack/echo/v4"
)

func MapCatalogueRoutes(e *echo.Echo, h ICatalogueHandlers, jwtMiddleware echo.MiddlewareFunc) {
	// categories route/s - return categories by parent path
	e.GET("/v1/catalogue/categories/*", m.Authorization(h.GetCataloguecategoriesByParentPath(), shared.ROLE_CATALOGUE_VIEW, shared.ROLE_CATALOGUE_EDIT, shared.ROLE_CATALOGUE_CATEGORY_EDIT), jwtMiddleware)
	e.GET("/v1/catalogue/categories", m.Authorization(h.GetCataloguecategoriesByParentPath(), shared.ROLE_CATALOGUE_VIEW, shared.ROLE_CATALOGUE_EDIT, shared.ROLE_CATALOGUE_CATEGORY_EDIT), jwtMiddleware)

	e.POST("/v1/catalogue/category", m.Authorization(h.CreateCatalogueCategory(), shared.ROLE_CATALOGUE_CATEGORY_EDIT), jwtMiddleware)
	e.GET("/v1/catalogue/category/:uid", m.Authorization(h.GetCatalogueCategoryWithDetailsByUid(), shared.ROLE_CATALOGUE_VIEW, shared.ROLE_CATALOGUE_EDIT, shared.ROLE_CATALOGUE_CATEGORY_EDIT), jwtMiddleware)
	e.PUT("/v1/catalogue/category/:uid", m.Authorization(h.UpdateCatalogueCategory(), shared.ROLE_CATALOGUE_CATEGORY_EDIT), jwtMiddleware)
	e.DELETE("/v1/catalogue/category/:uid", m.Authorization(h.DeleteCatalogueCategory(), shared.ROLE_CATALOGUE_CATEGORY_EDIT), jwtMiddleware)

	e.GET("/v1/catalogue/category/:uid/properties", m.Authorization(h.GetCatalogueCategoryPropertiesByUid(), shared.ROLE_CATALOGUE_VIEW, shared.ROLE_CATALOGUE_EDIT, shared.ROLE_CATALOGUE_CATEGORY_EDIT), jwtMiddleware)

	e.GET("/v1/catalogue/category/:uid/physical-item-properties", m.Authorization(h.GetCatalogueCategoryPhysicalItemPropertiesByUid(), shared.ROLE_CATALOGUE_VIEW, shared.ROLE_CATALOGUE_EDIT, shared.ROLE_CATALOGUE_CATEGORY_EDIT), jwtMiddleware)

	e.POST("/v1/catalogue/category/:uid/copy", m.Authorization(h.CopyCatalogueCategoryRecursive(), shared.ROLE_CATALOGUE_CATEGORY_EDIT), jwtMiddleware)

	//fake image get - get only eli logo img for now
	e.GET("/v1/catalogue/category/:uid/image", m.Authorization(h.GetCatalogueCategoryImageByUid(), shared.ROLE_CATALOGUE_VIEW, shared.ROLE_CATALOGUE_EDIT, shared.ROLE_CATALOGUE_CATEGORY_EDIT), jwtMiddleware)
	e.GET("/v1/catalogue/item/:uid/image", m.Authorization(h.GetCatalogueItemImage(), shared.ROLE_CATALOGUE_VIEW, shared.ROLE_CATALOGUE_EDIT, shared.ROLE_CATALOGUE_CATEGORY_EDIT), jwtMiddleware)

	// get all catalogue items with details
	e.GET("/v1/catalogue/items", m.Authorization(h.GetCatalogueItems(), shared.ROLE_CATALOGUE_VIEW, shared.ROLE_CATALOGUE_EDIT, shared.ROLE_CATALOGUE_CATEGORY_EDIT), jwtMiddleware)

	//get on catalogue item with details by uid
	e.GET("/v1/catalogue/item/:uid", m.Authorization(h.GetCatalogueItemWithDetailsByUid(), shared.ROLE_CATALOGUE_VIEW, shared.ROLE_CATALOGUE_EDIT, shared.ROLE_CATALOGUE_CATEGORY_EDIT), jwtMiddleware)
	//create new catalogue item
	e.POST("/v1/catalogue/item", m.Authorization(h.CreateNewCatalogueItem(), shared.ROLE_CATALOGUE_EDIT), jwtMiddleware)
	//update catalogue item
	e.PUT("/v1/catalogue/item/:uid", m.Authorization(h.UpdateCatalogueItem(), shared.ROLE_CATALOGUE_EDIT), jwtMiddleware)
	//delete catalogue item
	e.DELETE("/v1/catalogue/item/:uid", m.Authorization(h.DeleteCatalogueItem(), shared.ROLE_CATALOGUE_EDIT), jwtMiddleware)
	//get catalogue item statistics
	e.GET("/v1/catalogue/item/:uid/statistics", m.Authorization(h.GetCatalogueItemStatistics(), shared.ROLE_CATALOGUE_VIEW), jwtMiddleware)
	//get catalogue items overall statistics
	e.GET("/v1/catalogue/items/statistics", m.Authorization(h.CatalogueItemsOverallStatistics(), shared.ROLE_CATALOGUE_VIEW), jwtMiddleware)
	//get catalogue service type by uid
	//TODO: add roles
	e.GET("/v1/catalogue/service/type/:uid", m.Authorization(h.GetCatalogueServiceTypeByUid(), shared.ROLE_CATALOGUE_SERVICE_VIEW, shared.ROLE_CATALOGUE_SERVICE_EDIT), jwtMiddleware)
	//get catalogue service types
	e.GET("/v1/catalogue/service/types", m.Authorization(h.GetCatalogueServiceTypes(), shared.ROLE_CATALOGUE_SERVICE_VIEW, shared.ROLE_CATALOGUE_SERVICE_EDIT), jwtMiddleware)
	//create catalogue service type
	e.POST("/v1/catalogue/service/type", m.Authorization(h.CreateCatalogueServiceType(), shared.ROLE_CATALOGUE_SERVICE_EDIT), jwtMiddleware)
	//update catalogue service type
	e.PUT("/v1/catalogue/service/type/:uid", m.Authorization(h.UpdateCatalogueServiceType(), shared.ROLE_CATALOGUE_SERVICE_EDIT), jwtMiddleware)
	//delete catalogue service type
	e.DELETE("/v1/catalogue/service/type/:uid", m.Authorization(h.DeleteCatalogueServiceType(), shared.ROLE_CATALOGUE_SERVICE_EDIT), jwtMiddleware)
}
