package general

import (
	"panda/apigateway/services/general/models"

	"github.com/labstack/echo/v4"
)

type GeneralHandlers struct {
	generalService IGeneralService
}

type IGeneralHandlers interface {
	GetGraphByUid() echo.HandlerFunc
	GetUUID() echo.HandlerFunc
}

// NewGeneralHandlers General handlers constructor
func NewGeneralHandlers(generalSvc IGeneralService) IGeneralHandlers {
	return &GeneralHandlers{generalService: generalSvc}
}

// GetGraphByUid Get graph by uid godoc
// @Summary Get graph by uid
// @Description Get graph by uid
// @Tags General
// @Security BearerAuth
// @Produce json
// @Param uid path string true "uid"
// @Success 200 {object} models.GraphResponse
// @Failure 500 "Internal Server Error"
// @Router /v1/general/{uid}/graph [get]
func (h *GeneralHandlers) GetGraphByUid() echo.HandlerFunc {

	return func(c echo.Context) error {

		uid := c.Param("uid")

		nodes, err := h.generalService.GetGraphNodesByUid(uid)
		if err != nil {
			return echo.ErrInternalServerError
		}

		links, err := h.generalService.GetGraphLinksByUid(uid)
		if err != nil {
			return echo.ErrInternalServerError
		}

		result := models.GraphResponse{}
		result.Nodes = nodes
		result.Links = links

		return c.JSON(200, result)
	}
}

// GetUUID Get UUID v4 godoc
// @Summary Get UUID V4
// @Description Get UUID v4 string
// @Tags General
// @Security BearerAuth
// @Produce plain
// @Success 200 {string} string
// @Failure 500 "Internal Server Error"
// @Router /v1/uuid/v4 [get]
func (h *GeneralHandlers) GetUUID() echo.HandlerFunc {

	return func(c echo.Context) error {

		uuid, err := h.generalService.GetUUID()
		if err != nil {
			return echo.ErrInternalServerError
		}

		return c.String(200, uuid)
	}
}
