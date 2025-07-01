package roomcardsservice

import (
	m "panda/apigateway/middlewares"
	"panda/apigateway/shared"

	"github.com/labstack/echo/v4"
)

func MapRoomCardsRoutes(e *echo.Echo, h IRoomCardsHandlers, jwtMiddleware echo.MiddlewareFunc) {
	e.GET("/v1/room-card/layout/location/:code", m.Authorization(h.GetLayoutRoomCardInfo(), shared.ROLE_ROOM_CARDS_VIEW), jwtMiddleware)
}
