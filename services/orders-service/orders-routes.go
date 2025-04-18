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

	e.GET("/v1/order/:uid", m.Authorization(h.GetOrderWithOrderLinesByUid(), shared.ROLE_ORDERS_VIEW, shared.ROLE_ORDERS_EDIT), jwtMiddleware)

	e.POST("/v1/order", m.Authorization(h.InsertNewOrder(), shared.ROLE_ORDERS_EDIT), jwtMiddleware)

	e.DELETE("/v1/order/:uid", m.Authorization(h.DeleteOrder(), shared.ROLE_ORDERS_EDIT), jwtMiddleware)

	e.PUT("/v1/order/:uid", m.Authorization(h.UpdateOrder(), shared.ROLE_ORDERS_EDIT), jwtMiddleware)

	e.PUT("/v1/order/:uid/orderline/:itemUid/delivery", m.Authorization(h.UpdateOrderLineDelivery(), shared.ROLE_ORDERS_DELIVERY_EDIT, shared.ROLE_ORDERS_EDIT), jwtMiddleware)
	e.PUT("/v1/order/:uid/serviceline/:serviceItemUid/delivery", m.Authorization(h.UpdateServiceLineDelivery(), shared.ROLE_ORDERS_DELIVERY_EDIT, shared.ROLE_ORDERS_EDIT), jwtMiddleware)

	e.PUT("/v1/order/:uid/orderlines/delivery", m.Authorization(h.UpdateMultipleOrderLineDelivery(), shared.ROLE_ORDERS_DELIVERY_EDIT, shared.ROLE_ORDERS_EDIT), jwtMiddleware)
	e.PUT("/v1/order/:uid/servicelines/delivery", m.Authorization(h.UpdateMultipleServiceLineDelivery(), shared.ROLE_ORDERS_DELIVERY_EDIT, shared.ROLE_ORDERS_EDIT), jwtMiddleware)

	e.GET("/v1/orders/eun-for-print", m.Authorization(h.GetItemsForEunPrint(), shared.ROLE_ORDERS_VIEW, shared.ROLE_ORDERS_DELIVERY_EDIT), jwtMiddleware)

	e.PUT("/v1/orders/eun-for-print/:eun", m.Authorization(h.SetItemPrintEUN(), shared.ROLE_ORDERS_VIEW, shared.ROLE_ORDERS_DELIVERY_EDIT), jwtMiddleware)

	e.GET("/v1/order-uid-by-order-number/:orderNumber", h.GetOrderUidByOrderNumber())

	e.GET("/v1/catalogue/:catalogueItemUid/orders", m.Authorization(h.GetOrdersForCatalogueItem(), shared.ROLE_ORDERS_VIEW, shared.ROLE_ORDERS_EDIT), jwtMiddleware)

	e.GET("/v1/orders/order-lines/min-max-prices", m.Authorization(h.GetMinAndMaxOrderLinePrice(), shared.ROLE_ORDERS_VIEW, shared.ROLE_ORDERS_EDIT), jwtMiddleware)
}
