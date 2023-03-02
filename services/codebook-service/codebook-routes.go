package codebookService

import (
	"github.com/labstack/echo/v4"
)

func MapCodebookRoutes(e *echo.Echo, h ICodebookHandlers, jwtMiddleware echo.MiddlewareFunc) {
	// get codebook by fixed codebook code and filter optionaly by parentUID
	e.GET("/v1/codebook/:codebookCode", h.GetCodebook(), jwtMiddleware)
	// get autocomplete variant of the codebook - get results by searchText - limit param is optional - default 10
	e.GET("/v1/codebook/autocomplete/:codebookCode", h.GetAutocompleteCodebook(), jwtMiddleware)
}
