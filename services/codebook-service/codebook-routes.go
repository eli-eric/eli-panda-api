package codebookService

import (
	m "panda/apigateway/middlewares"
	"panda/apigateway/shared"

	"github.com/labstack/echo/v4"
)

func MapCodebookRoutes(e *echo.Echo, h ICodebookHandlers, jwtMiddleware echo.MiddlewareFunc) {
	// get codebook by fixed codebook code and filter optionaly by parentUID
	e.GET("/v1/codebook/:codebookCode", m.Authorization(h.GetCodebook(), shared.ROLE_SYSTEMS_VIEW, shared.ROLE_SYSTEMS_EDIT, shared.ROLE_CATALOGUE_CATEGORY_EDIT, shared.ROLE_CATALOGUE_EDIT, shared.ROLE_CATALOGUE_VIEW, shared.ROLE_ORDERS_VIEW, shared.ROLE_ORDERS_EDIT, shared.ROLE_PUBLICATIONS_VIEW, shared.ROLE_PUBLICATIONS_EDIT), jwtMiddleware)
	// get codebook tree by fixed codebook code
	e.GET("/v1/codebook/:codebookCode/tree", m.Authorization(h.GetCodebookTree(), shared.ROLE_SYSTEMS_VIEW, shared.ROLE_SYSTEMS_EDIT, shared.ROLE_CATALOGUE_CATEGORY_EDIT, shared.ROLE_CATALOGUE_EDIT, shared.ROLE_CATALOGUE_VIEW, shared.ROLE_ORDERS_VIEW, shared.ROLE_ORDERS_EDIT, shared.ROLE_PUBLICATIONS_VIEW, shared.ROLE_PUBLICATIONS_EDIT), jwtMiddleware)
	// get list of all codebooks
	e.GET("/v1/codebooks", m.Authorization(h.GetListOfCodebooks(), shared.ROLE_BASICS_VIEW), jwtMiddleware)
	// create new codebook
	e.POST("/v1/codebook/:codebookCode", m.Authorization(h.CreateNewCodebook(), shared.ROLE_BASICS_VIEW), jwtMiddleware)
	// update codebook
	e.PUT("/v1/codebook/:codebookCode/:uid", m.Authorization(h.UpdateCodebook(), shared.ROLE_BASICS_VIEW), jwtMiddleware)
	// delete codebook
	e.DELETE("/v1/codebook/:codebookCode/:uid", m.Authorization(h.DeleteCodebook(), shared.ROLE_BASICS_VIEW), jwtMiddleware)
}
