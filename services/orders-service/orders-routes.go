package ordersService

import (
	m "panda/apigateway/middlewares"
	"panda/apigateway/shared"

	"github.com/labstack/echo/v4"
)

func MapOrdersRoutes(e *echo.Echo, h IOrdersHandlers, jwtMiddleware echo.MiddlewareFunc) {
	//get order statuses codebook
	e.GET("/v1/orders/statuses", m.Authorization(h.GetOrderStatusesCodebook(), shared.ROLE_ORDERS_VIEW, shared.ROLE_ORDERS_EDIT), jwtMiddleware)

	e.GET("/v1/orders", m.Authorization(h.GetOrdersWithSearchAndPagination(), shared.ROLE_ORDERS_VIEW, shared.ROLE_ORDERS_EDIT), jwtMiddleware)
}
