package codebookService

import (
	"github.com/labstack/echo/v4"
)

func MapCodebookRoutes(e *echo.Echo, h ICodebookHandlers, jwtMiddleware echo.MiddlewareFunc) {
	// get codebook by fixed codebook code and filter optionaly by parentUID
	e.GET("/v1/codebook/:codebookCode", h.GetCodebook(), jwtMiddleware)
}
