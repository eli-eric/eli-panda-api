package catalogueService

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CatalogueHandlers struct {
	catalogueService ICatalogueService
}

type ICatalogueHandlers interface {
	GetCataloguecategoriesByParentPath() echo.HandlerFunc
	GetCatalogueCategoryImage() echo.HandlerFunc
}

// NewCommentsHandlers Comments handlers constructor
func NewCatalogueHandlers(catalogueSvc ICatalogueService) ICatalogueHandlers {
	return &CatalogueHandlers{catalogueService: catalogueSvc}
}

func (h *CatalogueHandlers) GetCataloguecategoriesByParentPath() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get query path param
		parentPath := c.Param("*")
		// get all categories of the given parentPath
		categories, err := h.catalogueService.GetCataloguecategoriesByParentPath(parentPath)

		if err == nil {
			return c.JSON(http.StatusOK, categories)
		}

		return echo.ErrInternalServerError
	}
}

func (h *CatalogueHandlers) GetCatalogueCategoryImage() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get query path param
		categoryUid := c.Param("uid")

		fmt.Println(categoryUid)

		imgData := "open-api-specification/eli-logo-small.png"

		return c.File(imgData)

		// categories, err := h.catalogueService.GetCataloguecategoriesByParentPath(parentPath)

		// if err == nil {
		// 	return c.JSON(http.StatusOK, categories)
		// }

		// return echo.ErrInternalServerError
	}
}
