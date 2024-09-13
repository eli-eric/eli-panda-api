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
}

// NewGeneralHandlers General handlers constructor
func NewGeneralHandlers(generalSvc IGeneralService) IGeneralHandlers {
	return &GeneralHandlers{generalService: generalSvc}
}

// GetGraphByUid Get graph by uid godoc
// @Summary Get graph by uid
// @Description Get graph by uid
// @Tags general
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
