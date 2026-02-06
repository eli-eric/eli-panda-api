package cronservice

import (
	"net/http"
	"panda/apigateway/services/cron-service/models"

	"github.com/rs/zerolog/log"

	"github.com/labstack/echo/v4"
)

var _ = models.CronJobHistory{}

type CronHandlers struct {
	cronService ICronService
}

type ICronHandlers interface {
	GetCronJobHistory() echo.HandlerFunc
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
