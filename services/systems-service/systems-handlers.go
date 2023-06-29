package systemsService

import (
	"encoding/json"
	"net/http"
	"panda/apigateway/helpers"
	"panda/apigateway/services/systems-service/models"

	"github.com/rs/zerolog/log"

	"github.com/labstack/echo/v4"
)

type SystemsHandlers struct {
	systemsService ISystemsService
}

type ISystemsHandlers interface {
	GetSubSystemsByParentUID() echo.HandlerFunc
	GetSystemImageByUid() echo.HandlerFunc
	GetSystemDetail() echo.HandlerFunc
	CreateNewSystem() echo.HandlerFunc
	UpdateSystem() echo.HandlerFunc
	DeleteSystemRecursive() echo.HandlerFunc
	GetSystemsWithSearchAndPagination() echo.HandlerFunc
}

// NewCommentsHandlers Comments handlers constructor
func NewsystemsHandlers(systemsSvc ISystemsService) ISystemsHandlers {
	return &SystemsHandlers{systemsService: systemsSvc}
}

func (h *SystemsHandlers) GetSubSystemsByParentUID() echo.HandlerFunc {

	return func(c echo.Context) error {

		parentUID := c.Param("parentUID")
		facilityCode := c.Get("facilityCode").(string)
		subSystems, err := h.systemsService.GetSubSystemsByParentUID(parentUID, facilityCode)

		if err == nil {
			return c.JSON(http.StatusOK, subSystems)
		} else {
			log.Error().Msg(err.Error())
		}

		return echo.ErrInternalServerError
	}
}

func (h *SystemsHandlers) GetSystemImageByUid() echo.HandlerFunc {
	return func(c echo.Context) error {

		//get uid path param
		uid := c.Param("uid")

		imageString, err := h.systemsService.GetSystemImageByUid(uid)

		if err == nil {
			return c.String(200, imageString)

		} else {
			log.Error().Msg(err.Error())
		}

		return echo.ErrInternalServerError
	}
}

func (h *SystemsHandlers) GetSystemDetail() echo.HandlerFunc {

	return func(c echo.Context) error {

		uid := c.Param("uid")
		facilityCode := c.Get("facilityCode").(string)
		systemDetail, err := h.systemsService.GetSystemDetail(uid, facilityCode)

		if err == nil {
			return c.JSON(http.StatusOK, systemDetail)
		} else {
			log.Error().Msg(err.Error())
		}

		return echo.ErrInternalServerError

	}
}

func (h *SystemsHandlers) CreateNewSystem() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		system := new(models.System)

		if err := c.Bind(system); err == nil {

			facilityCode := c.Get("facilityCode").(string)
			userUID := c.Get("userUID").(string)

			uid, err := h.systemsService.CreateNewSystem(system, facilityCode, userUID)

			if err == nil {
				return c.String(http.StatusCreated, uid)
			}

			return echo.ErrInternalServerError

		} else {
			log.Error().Msg(err.Error())
		}
		return echo.ErrBadRequest
	}
}

func (h *SystemsHandlers) UpdateSystem() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		system := new(models.System)

		if err := c.Bind(system); err == nil {

			facilityCode := c.Get("facilityCode").(string)
			userUID := c.Get("userUID").(string)

			err := h.systemsService.UpdateSystem(system, facilityCode, userUID)

			if err == nil {
				return c.NoContent(http.StatusNoContent)
			}

			return echo.ErrInternalServerError

		}
		return echo.ErrBadRequest
	}
}

func (h *SystemsHandlers) DeleteSystemRecursive() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get uid path param
		uid := c.Param("uid")

		// get catalogue item
		err := h.systemsService.DeleteSystemRecursive(uid)

		if err == nil {
			return c.NoContent(http.StatusNoContent)
		}

		return echo.ErrInternalServerError
	}
}

func (h *SystemsHandlers) GetSystemsWithSearchAndPagination() echo.HandlerFunc {

	return func(c echo.Context) error {

		pagination := c.QueryParam("pagination")
		sorting := c.QueryParam("sorting")
		search := c.QueryParam("search")
		facilityCode := c.Get("facilityCode").(string)

		pagingObject := new(helpers.Pagination)
		json.Unmarshal([]byte(pagination), &pagingObject)

		sortingObject := new([]helpers.Sorting)
		json.Unmarshal([]byte(sorting), &sortingObject)

		items, err := h.systemsService.GetSystemsWithSearchAndPagination(search, facilityCode, pagingObject, sortingObject)

		if err == nil {
			return c.JSON(http.StatusOK, items)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}
