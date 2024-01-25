package cronservice

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/labstack/echo/v4"
)

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
