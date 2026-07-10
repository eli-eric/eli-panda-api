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

// GetSystemUIDByItemUID resolves the System containing the given physical item, so the edit-guard
// can be applied to physical-item endpoints. Returns helpers.ERR_NO_ROWS if the item has no system.
func (svc *SystemsService) GetSystemUIDByItemUID(itemUID string) (string, error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()
	return helpers.GetNeo4jSingleRecordSingleValue[string](session, GetSystemUIDByItemUIDQuery(itemUID))
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

// guardSystemEdit returns echo.ErrForbidden (403) if the caller may not edit ANY of the given
// systems; nil if allowed. A non-existent system (helpers.ERR_NO_ROWS) passes the guard so the
// handler's own existence check produces the appropriate 404. Empty uids are skipped.
func (h *SystemsHandlers) guardSystemEdit(c echo.Context, systemUIDs ...string) error {
	userUID := c.Get("userUID").(string)
	for _, uid := range systemUIDs {
		if uid == "" {
			continue
		}
		res, err := h.systemsService.CanEditSystem(uid, userUID)
		if err != nil {
			if errors.Is(err, helpers.ERR_NO_ROWS) {
				continue
			}
			log.Error().Err(err).Msg("guardSystemEdit")
			return echo.ErrInternalServerError
		}
		if !res.Result {
			return echo.ErrForbidden
		}
	}
	return nil
}

// guardSystemEditByItem resolves the physical item's system and applies guardSystemEdit to it.
// If the item belongs to no system (ERR_NO_ROWS), the guard passes (handler handles the miss).
func (h *SystemsHandlers) guardSystemEditByItem(c echo.Context, itemUID string) error {
	systemUID, err := h.systemsService.GetSystemUIDByItemUID(itemUID)
	if err != nil {
		if errors.Is(err, helpers.ERR_NO_ROWS) {
			return nil
		}
		log.Error().Err(err).Msg("guardSystemEditByItem")
		return echo.ErrInternalServerError
	}
	return h.guardSystemEdit(c, systemUID)
}
