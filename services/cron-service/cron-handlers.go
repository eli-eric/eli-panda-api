package cronservice

import (
	"net/http"
	"panda/apigateway/services/cron-service/models"

	"github.com/rs/zerolog/log"

	"github.com/labstack/echo/v4"
)

type CronHandlers struct {
	cronService ICronService
}

type ICronHandlers interface {
	GetCronJobHistory() echo.HandlerFunc
	SyncEmployees() echo.HandlerFunc
}

// NewCommentsHandlers Comments handlers constructor
func NewCronHandlers(cronSvc ICronService) ICronHandlers {
	return &CronHandlers{cronService: cronSvc}
}

// GetCronJobHistory godoc
// @Summary Get cron job history
// @Description Returns cron job execution history.
// @Tags Cron
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.CronJobHistory
// @Failure 500 "Internal server error"
// @Router /v1/cron/history [get]
func (h *CronHandlers) GetCronJobHistory() echo.HandlerFunc {
	return func(c echo.Context) error {

		cronJobHistory, err := h.cronService.GetCronJobHistory()

		if err == nil {
			return c.JSON(http.StatusOK, cronJobHistory)
		} else {
			log.Error().Msg(err.Error())
		}

		return echo.ErrInternalServerError
	}
}

// SyncEmployees godoc
// @Summary Run employees sync from HR (OKBase) manually
// @Description Runs the ELI Beamlines employees sync immediately (normally scheduled daily at 02:00).
// @Tags Cron
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.SyncEmployeesResult
// @Failure 500 "Internal server error"
// @Router /v1/cron/sync-employees [post]
func (h *CronHandlers) SyncEmployees() echo.HandlerFunc {
	return func(c echo.Context) error {

		savedCount, errs := h.cronService.SyncEmployees()

		for _, err := range errs {
			log.Error().Msgf("SyncEmployees error: %v", err.Error())
		}

		return c.JSON(http.StatusOK, models.SyncEmployeesResult{SavedCount: savedCount, ErrorCount: len(errs)})
	}
}
