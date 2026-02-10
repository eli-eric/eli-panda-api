package roomcardsservice

import (
	"net/http"
	"panda/apigateway/helpers"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type RoomCardsHandlers struct {
	RoomCardsService IRoomCardsService
}

type IRoomCardsHandlers interface {
	GetLayoutRoomCardInfo() echo.HandlerFunc
}

func NewRoomCardsHandlers(service IRoomCardsService) IRoomCardsHandlers {
	return &RoomCardsHandlers{RoomCardsService: service}
}

// GetLayoutRoomCardInfo Get room card layout info by location code godoc
// @Summary Get room card layout info by location code
// @Description Get room card layout information for external systems to display scientific complex layout
// @Description
// @Description Status values:
// @Description - DIRTY_MODE: Room is dirty and needs cleaning
// @Description - CLEAN_MODE: Room is clean and ready for use
// @Description - IN_PREPARATION_MODE: Room is being prepared for use
// @Description
// @Description PurityClass values (ISO standards):
// @Description - ISO_5: ISO Class 5 cleanroom
// @Description - ISO_6: ISO Class 6 cleanroom
// @Description - ISO_7: ISO Class 7 cleanroom
// @Description - ISO_8: ISO Class 8 cleanroom
// @Description
// @Description StatusColor values (automatically assigned based on status):
// @Description - DIRTY_MODE: #fecaca (light red)
// @Description - CLEAN_MODE: #d9f99d (light green)
// @Description - IN_PREPARATION_MODE: #fdba74 (light orange)
// @Description - Unknown status: #808080 (gray)
// @Description
// @Description OperationalState: Codebook containing uid, name, and code of the operational state
// @Description
// @Description OperationalStateColor values (automatically assigned based on operationalState.code):
// @Description - OS1: #22c55e (green)
// @Description - OS2: #86efac (light green)
// @Description - OS3: #facc15 (yellow)
// @Description - OS4: #fb923c (orange)
// @Description - OS5: #ef4444 (red)
// @Description - OS6: #b91c1c (dark red)
// @Description - Unknown/null: #9ca3af (gray)
// @Tags Room Cards
// @Security BearerAuth
// @Produce json
// @Param code path string true "Location code"
// @Param fields query string false "Comma-separated list of top-level response fields to return (e.g. uid,name). If provided, the response is a partial object."
// @Success 200 {object} models.LayoutRoomCard
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /v1/room-card/layout/location/{code} [get]
func (h *RoomCardsHandlers) GetLayoutRoomCardInfo() echo.HandlerFunc {
	return func(c echo.Context) error {
		locationCode := c.Param("code")

		if locationCode == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Location code is required"})
		}

		result, err := h.RoomCardsService.GetLayoutRoomCardInfo(locationCode)
		if err != nil {
			// return 404 if not found - in error message will be result contains no more records
			if err.Error() == "Result contains no more records" {
				return echo.ErrNotFound
			}
			log.Error().Err(err).Msg("Error getting room cards for location")
			return echo.ErrInternalServerError
		}

		fieldsParam := c.QueryParam("fields")
		if fieldsParam == "" {
			return c.JSON(http.StatusOK, result)
		}

		fields := helpers.ParseFieldsParam(fieldsParam)
		if len(fields) == 0 {
			return c.JSON(http.StatusBadRequest, map[string]any{"error": "fields parameter is empty"})
		}
		if len(fields) > 50 {
			return c.JSON(http.StatusBadRequest, map[string]any{"error": "too many fields requested", "max": 50})
		}

		projected, err := helpers.ProjectJSONFields(result, fields)
		if err != nil {
			if projErr, ok := err.(*helpers.FieldProjectionError); ok {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"error":         "invalid fields",
					"invalidFields": projErr.InvalidFields,
					"allowedFields": projErr.AllowedFields,
				})
			}
			log.Error().Err(err).Msg("Error projecting room card response fields")
			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, projected)
	}
}
