package catalogueService

import (
	"fmt"
	"log"
	"net/http"
	"panda/apigateway/helpers"
	"panda/apigateway/services/catalogue-service/models"

	"github.com/labstack/echo/v4"
)

type CatalogueHandlers struct {
	catalogueService ICatalogueService
}

type ICatalogueHandlers interface {
	GetCataloguecategoriesByParentPath() echo.HandlerFunc
	GetCatalogueCategoryImage() echo.HandlerFunc
	GetCatalogueItems() echo.HandlerFunc
	GetCatalogueItemImage() echo.HandlerFunc
	GetCatalogueItemWithDetailsByUid() echo.HandlerFunc
	GetCatalogueCategoryWithDetailsByUid() echo.HandlerFunc
	UpdateCatalogueCategory() echo.HandlerFunc
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
		categories, err := h.catalogueService.GetCatalogueCategoriesByParentPath(parentPath)

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

		imgData := "assets/no-image.png"

		return c.File(imgData)
	}
}

func (h *CatalogueHandlers) GetCatalogueItemImage() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get query path param
		categoryUid := c.Param("uid")

		fmt.Println(categoryUid)

		imgData := "assets/no-image.png"

		return c.File(imgData)
	}
}

func (h *CatalogueHandlers) GetCatalogueItems() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get query path param
		searchParam := c.QueryParams().Get("search")
		categoryPathParam := c.QueryParams().Get("categoryPath")

		pagination := new(helpers.Pagination)
		if err := c.Bind(pagination); err == nil {

			// get catalogue items
			items, err := h.catalogueService.GetCatalogueItems(searchParam, categoryPathParam, pagination.PageSize, pagination.Page)

			if err == nil {
				return c.JSON(http.StatusOK, items)
			} else {
				log.Println(err)
			}
		}

		return echo.ErrInternalServerError
	}
}

func (h *CatalogueHandlers) GetCatalogueItemWithDetailsByUid() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get uid path param
		uid := c.Param("uid")

		// get catalogue item
		item, err := h.catalogueService.GetCatalogueItemWithDetailsByUid(uid)

		if err == nil {
			return c.JSON(http.StatusOK, item)
		} else {
			log.Println(err)
		}

		return echo.ErrInternalServerError
	}
}

func (h *CatalogueHandlers) GetCatalogueCategoryWithDetailsByUid() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get uid path param
		uid := c.Param("uid")

		// get catalogue item
		item, err := h.catalogueService.GetCatalogueCategoryWithDetailsByUid(uid)

		if err == nil {
			return c.JSON(http.StatusOK, item)
		} else {
			log.Println(err)
		}

		return echo.ErrInternalServerError
	}
}

func (h *CatalogueHandlers) UpdateCatalogueCategory() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		catalogueCatgeory := new(models.CatalogueCategory)

		if err := c.Bind(catalogueCatgeory); err == nil {
			// update catalogue item
			err := h.catalogueService.UpdateCatalogueCategory(catalogueCatgeory)

			if err == nil {
				return c.JSON(http.StatusOK, catalogueCatgeory)
			} else {
				log.Println(err)
			}
		} else {
			log.Println(err)
		}

		return echo.ErrInternalServerError
	}
}
