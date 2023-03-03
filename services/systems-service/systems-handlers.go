package systemsService

import (
	"log"
	"net/http"
	"panda/apigateway/helpers"

	"github.com/labstack/echo/v4"
)

type SystemsHandlers struct {
	systemsService ISystemsService
}

type ISystemsHandlers interface {
	GetSubSystemsByParentUID() echo.HandlerFunc
	GetSystemImageByUid() echo.HandlerFunc
}

// NewCommentsHandlers Comments handlers constructor
func NewsystemsHandlers(systemsSvc ISystemsService) ISystemsHandlers {
	return &SystemsHandlers{systemsService: systemsSvc}
}

func (h *SystemsHandlers) GetSubSystemsByParentUID() echo.HandlerFunc {

	return func(c echo.Context) error {
		if userInfo := helpers.IsUserInRole(c, helpers.ROLE_SYSTEMS_VIEW); userInfo != nil {
			parentUID := c.Param("parentUID")
			subSystems, err := h.systemsService.GetSubSystemsByParentUID(parentUID, userInfo)

			if err == nil {
				return c.JSON(http.StatusOK, subSystems)
			}

			return echo.ErrInternalServerError
		} else {
			return echo.ErrUnauthorized
		}
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
			log.Println(err)
		}

		return echo.ErrInternalServerError
	}
}
