package systemsService

import (
	"encoding/json"
	"net/http"
	"panda/apigateway/helpers"
	"panda/apigateway/services/systems-service/models"
	"strconv"

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
	GetSystemsForRelationship() echo.HandlerFunc
	GetSystemRelationships() echo.HandlerFunc
	DeleteSystemRelationship() echo.HandlerFunc
	CreateNewSystemRelationship() echo.HandlerFunc
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
			system.UID = c.Param("uid")

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

func (h *SystemsHandlers) GetSystemsForRelationship() echo.HandlerFunc {

	return func(c echo.Context) error {

		facilityCode := c.Get("facilityCode").(string)
		search := c.QueryParam("search")

		pagingObject := new(helpers.Pagination)
		pagination := c.QueryParam("pagination")
		json.Unmarshal([]byte(pagination), &pagingObject)

		sortingObject := new([]helpers.Sorting)
		sorting := c.QueryParam("sorting")
		json.Unmarshal([]byte(sorting), &sortingObject)

		systemFromUid := c.QueryParam("systemFromUid")
		relationTypeCode := c.QueryParam("relationTypeCode")

		items, err := h.systemsService.GetSystemsForRelationship(search, facilityCode, pagingObject, sortingObject, systemFromUid, relationTypeCode)

		if err == nil {
			return c.JSON(http.StatusOK, items)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

func (h *SystemsHandlers) GetSystemRelationships() echo.HandlerFunc {

	return func(c echo.Context) error {

		systemUid := c.Param("uid")

		items, err := h.systemsService.GetSystemRelationships(systemUid)

		if err == nil {
			return c.JSON(http.StatusOK, items)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

func (h *SystemsHandlers) DeleteSystemRelationship() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get userUID
		userUID := c.Get("userUID").(string)
		//get uid path param
		uid := c.Param("uid")
		//convert uid to int64
		uidInt64, err := strconv.ParseInt(uid, 10, 64)

		if err != nil {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}

		err = h.systemsService.DeleteSystemRelationship(uidInt64, userUID)

		if err == nil {
			return c.NoContent(http.StatusNoContent)
		}

		return echo.ErrInternalServerError
	}
}

func (h *SystemsHandlers) CreateNewSystemRelationship() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		systemRelationshipRequest := new(models.SystemRelationshipRequest)

		if err := c.Bind(systemRelationshipRequest); err == nil {

			userUID := c.Get("userUID").(string)
			facilityCode := c.Get("facilityCode").(string)

			newId, err := h.systemsService.CreateNewSystemRelationship(systemRelationshipRequest, facilityCode, userUID)

			if err == nil {
				return c.String(http.StatusCreated, strconv.FormatInt(newId, 10))
			}

			return echo.ErrInternalServerError

		} else {
			log.Error().Msg(err.Error())
		}
		return echo.ErrBadRequest
	}
}
