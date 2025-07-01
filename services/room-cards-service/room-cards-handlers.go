package roomcardsservice

import (
	"net/http"

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
// @Tags Room Cards
// @Security BearerAuth
// @Produce json
// @Param code path string true "Location code"
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

		return c.JSON(http.StatusOK, result)
	}
}