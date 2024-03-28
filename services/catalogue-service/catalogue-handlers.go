package catalogueService

import (
	"encoding/json"
	"net/http"
	"panda/apigateway/helpers"
	"panda/apigateway/services/catalogue-service/models"

	"github.com/rs/zerolog/log"

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
	CreateNewCatalogueItem() echo.HandlerFunc
	GetCatalogueCategoryPropertiesByUid() echo.HandlerFunc
	UpdateCatalogueItem() echo.HandlerFunc
	DeleteCatalogueItem() echo.HandlerFunc
	GetCatalogueItemStatistics() echo.HandlerFunc
	CatalogueItemsOverallStatistics() echo.HandlerFunc
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
		categoryUidParam := c.QueryParams().Get("categoryUID")

		pagination := new(helpers.Pagination)
		if err := c.Bind(pagination); err == nil {

			filterObject := new([]helpers.ColumnFilter)
			filter := c.QueryParam("columnFilter")
			json.Unmarshal([]byte(filter), &filterObject)

			sortingObject := new([]helpers.Sorting)
			sorting := c.QueryParam("sorting")
			json.Unmarshal([]byte(sorting), &sortingObject)

			// get catalogue items
			items, err := h.catalogueService.GetCatalogueItems(searchParam, categoryUidParam, pagination.PageSize, pagination.Page, filterObject, sortingObject)

			if err == nil {
				return c.JSON(http.StatusOK, items)
			} else {
				log.Error().Msg(err.Error())
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
			log.Error().Msg(err.Error())
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
			log.Error().Msg(err.Error())
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
			log.Error().Msg(err.Error())
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
			log.Error().Msg(err.Error())
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
				log.Error().Msg(err.Error())
			}
		} else {
			log.Error().Msg(err.Error())
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
				log.Error().Msg(err.Error())
			}
		} else {
			log.Error().Msg(err.Error())
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
			log.Error().Msg(err.Error())
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
			log.Error().Msg(err.Error())
		}

		return echo.ErrInternalServerError
	}
}

func (h *CatalogueHandlers) CreateNewCatalogueItem() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue item data from request body
		catalogueItem := new(models.CatalogueItem)

		if err := c.Bind(catalogueItem); err == nil {

			userUID := c.Get("userUID").(string)
			// create catalogue item
			newUid, err := h.catalogueService.CreateNewCatalogueItem(catalogueItem, userUID)

			if err == nil {
				return c.JSON(http.StatusOK, newUid)
			} else {
				log.Error().Msg(err.Error())
			}
		} else {
			log.Error().Msg(err.Error())
		}

		return echo.ErrInternalServerError
	}
}

func (h *CatalogueHandlers) GetCatalogueCategoryPropertiesByUid() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get uid path param
		uid := c.Param("uid")
		//get item uid from query
		itemUID := c.QueryParam("itemUid")

		// get catalogue item
		properties, err := h.catalogueService.GetCatalogueCategoryPropertiesByUid(uid, &itemUID)

		if err == nil {
			return c.JSON(http.StatusOK, properties)
		} else {
			log.Error().Msg(err.Error())
		}

		return echo.ErrInternalServerError
	}
}

func (h *CatalogueHandlers) UpdateCatalogueItem() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue item data from request body
		catalogueItem := new(models.CatalogueItem)

		if err := c.Bind(catalogueItem); err == nil {

			uid := c.Param("uid")
			catalogueItem.Uid = uid

			userUID := c.Get("userUID").(string)
			// create catalogue item
			err := h.catalogueService.UpdateCatalogueItem(catalogueItem, userUID)

			if err == nil {
				//return c.NoContent(http.StatusNoContent)
				return c.JSON(http.StatusOK, catalogueItem.Uid)
			} else {
				log.Error().Msg(err.Error())
			}
		} else {
			log.Error().Msg(err.Error())
		}

		return echo.ErrInternalServerError
	}
}

func (h *CatalogueHandlers) DeleteCatalogueItem() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get uid path param
		uid := c.Param("uid")
		//get user uid from context
		userUID := c.Get("userUID").(string)

		// delete catalogue item
		err := h.catalogueService.DeleteCatalogueItem(uid, userUID)

		if err == nil {
			return c.NoContent(http.StatusNoContent)
		} else {
			log.Error().Msg(err.Error())
			if err.Error() == helpers.ERR_DELETE_RELATED_ITEMS.Error() {
				return echo.NewHTTPError(http.StatusConflict, err.Error())
			}
		}

		return echo.ErrInternalServerError
	}
}

func (h *CatalogueHandlers) GetCatalogueItemStatistics() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get uid path param
		uid := c.Param("uid")

		statistics, err := h.catalogueService.GetCatalogueItemStatistics(uid)

		if err == nil {
			helpers.ProcessArrayResult[models.CatalogueStatistics](&statistics, err)
			return c.JSON(http.StatusOK, statistics)
		} else {
			log.Error().Msg(err.Error())
		}

		return echo.ErrInternalServerError
	}
}

func (h *CatalogueHandlers) CatalogueItemsOverallStatistics() echo.HandlerFunc {

	return func(c echo.Context) error {

		statistics, err := h.catalogueService.CatalogueItemsOverallStatistics()

		if err == nil {
			helpers.ProcessArrayResult[models.CatalogueStatistics](&statistics, err)
			return c.JSON(http.StatusOK, statistics)
		} else {
			log.Error().Msg(err.Error())
		}

		return echo.ErrInternalServerError
	}
}
