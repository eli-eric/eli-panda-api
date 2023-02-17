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
	GetCatalogueItems() echo.HandlerFunc
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
	}
}

func (h *CatalogueHandlers) GetCatalogueItems() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get query path param
		searchParam := c.QueryParams().Get("search")

		pagination := new(Pagination)
		if err := c.Bind(pagination); err == nil {

			// get catalogue items
			items, err := h.catalogueService.GetCatalogueItems(searchParam, pagination.PageSize, pagination.Page)

			if err == nil {
				return c.JSON(http.StatusOK, items)
			}
		}

		return echo.ErrInternalServerError
	}
}

type Pagination struct {
	PageSize int `query:"pageSize"`
	Page     int `query:"page"`
}
