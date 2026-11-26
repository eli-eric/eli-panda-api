package systemsService

import (
	"errors"
	"net/http"
	"panda/apigateway/helpers"
	"panda/apigateway/services/systems-service/models"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// CanEditSystem is the single source of truth for the system edit-guard. It bubbles up the
// HAS_SUBSYSTEM chain and returns whether userUID may edit systemUID plus the responsible users.
// A non-existent/deleted system yields helpers.ERR_NO_ROWS (callers decide 404 vs. pass-through).
func (svc *SystemsService) CanEditSystem(systemUID, userUID string) (models.CanEditSystemResult, error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()
	return helpers.GetNeo4jSingleRecordAndMapToStruct[models.CanEditSystemResult](session, CanEditSystemQuery(systemUID, userUID))
}

// CanEditSystem godoc
// @Summary Can the current user edit this system
// @Description Returns whether the caller may edit the system (responsibility bubbles up HAS_SUBSYSTEM) and the responsible users to contact.
// @Tags Systems
// @Produce json
// @Security BearerAuth
// @Param uid path string true "System UID"
// @Success 200 {object} models.CanEditSystemResult
// @Failure 404 "System not found"
// @Failure 500 "Internal server error"
// @Router /v1/system/{uid}/can-edit [get]
func (h *SystemsHandlers) CanEditSystem() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Param("uid")
		userUID := c.Get("userUID").(string)

		result, err := h.systemsService.CanEditSystem(uid, userUID)
		if err != nil {
			if errors.Is(err, helpers.ERR_NO_ROWS) {
				return echo.ErrNotFound
			}
			log.Error().Err(err).Msg("CanEditSystem")
			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, result)
	}
}
