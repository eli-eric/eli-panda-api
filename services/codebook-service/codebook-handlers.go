package codebookService

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type CodebookHandlers struct {
	codebookService ICodebookService
}

type ICodebookHandlers interface {
	GetCodebook() echo.HandlerFunc
}

// NewCommentsHandlers Comments handlers constructor
func NewCodebookHandlers(codebookService ICodebookService) ICodebookHandlers {
	return &CodebookHandlers{codebookService: codebookService}
}

func (h *CodebookHandlers) GetCodebook() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get query path param
		codebookCode := c.Param("codebookCode")
		parentUID := c.QueryParams().Get("parentUID")
		// get all categories of the given parentPath
		codebookList, err := h.codebookService.GetCodebook(codebookCode, parentUID)

		if err == nil {
			return c.JSON(http.StatusOK, codebookList)
		}

		return echo.ErrInternalServerError
	}
}
