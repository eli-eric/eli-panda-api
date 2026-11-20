package cronservice

import (
	"errors"
	"net/http"

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
// @Description Runs the ELI Beamlines employees sync immediately (normally scheduled daily at 02:00). The run is recorded in the cron job history.
// @Tags Cron
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.SyncEmployeesResult
// @Failure 409 "Sync is already running"
// @Failure 502 "Employees could not be fetched from HR (OKBase)"
// @Router /v1/cron/sync-employees [post]
func (h *CronHandlers) SyncEmployees() echo.HandlerFunc {
	return func(c echo.Context) error {

		result, err := h.cronService.SyncEmployees()

		if errors.Is(err, ErrEmployeesSyncRunning) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}

		if err != nil {
			log.Error().Msgf("SyncEmployees error: %v", err.Error())
			return echo.NewHTTPError(http.StatusBadGateway, "Employees could not be fetched from HR")
		}

		return c.JSON(http.StatusOK, result)
	}
}
