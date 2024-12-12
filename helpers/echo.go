package helpers

import (
	"github.com/labstack/echo/v4"
)

func BadRequest(message string) error {
	return echo.NewHTTPError(400, "Bad request: "+message)
}
