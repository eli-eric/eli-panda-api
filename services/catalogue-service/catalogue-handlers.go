package catalogueService

import (
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
	GetCatalogueItems() echo.HandlerFunc
	GetCatalogueItemImage() echo.HandlerFunc
	GetCatalogueItemWithDetailsByUid() echo.HandlerFunc
	GetCatalogueCategoryWithDetailsByUid() echo.HandlerFunc
	UpdateCatalogueCategory() echo.HandlerFunc
	CreateCatalogueCategory() echo.HandlerFunc
	DeleteCatalogueCategory() echo.HandlerFunc
	GetCatalogueCategoryImageByUid() echo.HandlerFunc
	CopyCatalogueCategoryRecursive() echo.HandlerFunc
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

func (h *CatalogueHandlers) DeleteCatalogueCategory() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get uid path param
		uid := c.Param("uid")

		// get catalogue item
		err := h.catalogueService.DeleteCatalogueCategory(uid)

		if err == nil {
			return c.JSON(http.StatusOK, "ok")
		} else {
			log.Println(err)
			if err.Error() == "category has related items" {
				return echo.ErrForbidden
			}
		}

		return echo.ErrInternalServerError
	}
}

func (h *CatalogueHandlers) CopyCatalogueCategoryRecursive() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get uid path param
		uid := c.Param("uid")

		// get catalogue item
		newUID, err := h.catalogueService.CopyCatalogueCategoryRecursive(uid)

		if err == nil {
			return c.String(http.StatusCreated, newUID)
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
			// update catalogue category
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

func (h *CatalogueHandlers) CreateCatalogueCategory() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		catalogueCatgeory := new(models.CatalogueCategory)

		if err := c.Bind(catalogueCatgeory); err == nil {
			// create catalogue category
			err := h.catalogueService.CreateCatalogueCategory(catalogueCatgeory)

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

func (h *CatalogueHandlers) GetCatalogueItemImage() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get query path param
		uid := c.Param("uid")

		imageString, err := h.catalogueService.GetCatalogueItemImageByUid(uid)
		if err == nil {

			return c.String(200, imageString)

		} else {
			log.Println(err)
			return echo.ErrNotFound
		}
	}
}

func (h *CatalogueHandlers) GetCatalogueCategoryImageByUid() echo.HandlerFunc {
	return func(c echo.Context) error {

		//get uid path param
		uid := c.Param("uid")

		imageString, err := h.catalogueService.GetCatalogueCategoryImageByUid(uid)

		if err == nil {

			// // we have to be sure that we have base64 image string
			// if strings.Index(imageString, "data:image") == 0 {
			// 	baseSplit := strings.Split(imageString, ",")
			// 	mimeType := strings.Split(strings.Split(baseSplit[0], ":")[1], ";")[0]
			// 	data, err := base64.StdEncoding.DecodeString(baseSplit[1])

			// 	if err != nil {
			// 		return c.Blob(500, "image/*", nil)
			// 	}
			// 	return c.Blob(200, mimeType, data)
			// } else {
			// 	return c.Blob(200, "image/*", nil)
			// }

			return c.String(200, imageString)

		} else {
			log.Println(err)
		}

		return echo.ErrInternalServerError
	}
}
