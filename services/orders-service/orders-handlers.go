package ordersService

import (
	"encoding/json"
	"net/http"
	"panda/apigateway/helpers"
	codebookModels "panda/apigateway/services/codebook-service/models"
	"panda/apigateway/services/orders-service/models"

	"github.com/rs/zerolog/log"

	"github.com/labstack/echo/v4"
)

var _ = codebookModels.Codebook{}

type OrdersHandlers struct {
	ordersService IOrdersService
}

type IOrdersHandlers interface {
	GetOrderStatusesCodebook() echo.HandlerFunc
	GetOrdersWithSearchAndPagination() echo.HandlerFunc
	GetOrderWithOrderLinesByUid() echo.HandlerFunc
	InsertNewOrder() echo.HandlerFunc
	UpdateOrder() echo.HandlerFunc
	DeleteOrder() echo.HandlerFunc
	UpdateOrderLineDelivery() echo.HandlerFunc
	UpdateMultipleOrderLineDelivery() echo.HandlerFunc
	UpdateServiceLineDelivery() echo.HandlerFunc
	UpdateMultipleServiceLineDelivery() echo.HandlerFunc
	GetItemsForEunPrint() echo.HandlerFunc
	SetItemPrintEUN() echo.HandlerFunc
	GetOrderUidByOrderNumber() echo.HandlerFunc
	GetOrdersForCatalogueItem() echo.HandlerFunc
	GetMinAndMaxOrderLinePrice() echo.HandlerFunc
}

// NewCommentsHandlers Comments handlers constructor
func NewOrdersHandlers(ordersSvc IOrdersService) IOrdersHandlers {
	return &OrdersHandlers{ordersService: ordersSvc}
}

// GetOrderStatusesCodebook godoc
// @Summary Get order statuses
// @Description Returns order statuses codebook.
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Success 200 {array} codebookModels.Codebook
// @Failure 500 "Internal server error"
// @Router /v1/orders/statuses [get]
func (h *OrdersHandlers) GetOrderStatusesCodebook() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get all categories by parent path
		orderStatuses, err := h.ordersService.GetOrderStatusesCodebook()

		if err == nil {
			return c.JSON(http.StatusOK, orderStatuses)
		}

		return echo.ErrInternalServerError
	}
}

// GetOrdersWithSearchAndPagination godoc
// @Summary Get orders
// @Description Returns a paginated list of orders with optional search/sorting/filtering.
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Param pagination query string true "Pagination JSON (e.g. {\"page\":1,\"pageSize\":100})"
// @Param sorting query string false "Sorting JSON (array of {id, desc})"
// @Param search query string false "Search text"
// @Param columnFilter query string false "Column filter JSON"
// @Success 200 {object} helpers.PaginationResult[models.OrderListItem]
// @Failure 500 "Internal server error"
// @Router /v1/orders [get]
func (h *OrdersHandlers) GetOrdersWithSearchAndPagination() echo.HandlerFunc {

	return func(c echo.Context) error {

		pagination := c.QueryParam("pagination")
		sorting := c.QueryParam("sorting")
		search := c.QueryParam("search")

		facilityCode := c.Get("facilityCode").(string)

		pagingObject := new(helpers.Pagination)
		json.Unmarshal([]byte(pagination), &pagingObject)

		sortingObject := new([]helpers.Sorting)
		json.Unmarshal([]byte(sorting), &sortingObject)

		filterObject := new([]helpers.ColumnFilter)
		filter := c.QueryParam("columnFilter")
		json.Unmarshal([]byte(filter), &filterObject)

		items, err := h.ordersService.GetOrdersWithSearchAndPagination(search, facilityCode, pagingObject, sortingObject, filterObject)

		if err == nil {
			return c.JSON(http.StatusOK, items)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}

	}
}

// GetOrderWithOrderLinesByUid godoc
// @Summary Get order detail
// @Description Returns order detail with order lines by order UID.
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Param uid path string true "Order UID"
// @Success 200 {object} models.OrderDetail
// @Failure 500 "Internal server error"
// @Router /v1/order/{uid} [get]
func (h *OrdersHandlers) GetOrderWithOrderLinesByUid() echo.HandlerFunc {

	return func(c echo.Context) error {

		orderUid := c.Param("uid")
		facilityCode := c.Get("facilityCode").(string)

		order, err := h.ordersService.GetOrderWithOrderLinesByUid(orderUid, facilityCode)

		if err == nil {
			return c.JSON(http.StatusOK, order)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}

	}
}

// InsertNewOrder godoc
// @Summary Insert a new order
// @Description Insert a new order
// @Tags Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param order body models.OrderDetail true "Order object that needs to be inserted"
// @Success 200 {object} models.OrderDetail
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Router /v1/order [post]
func (h *OrdersHandlers) InsertNewOrder() echo.HandlerFunc {

	return func(c echo.Context) error {

		order := new(models.OrderDetail)
		err := c.Bind(order)

		if err != nil {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}

		facilityCode := c.Get("facilityCode").(string)
		userUID := c.Get("userUID").(string)

		newUID, err := h.ordersService.InsertNewOrder(order, facilityCode, userUID)

		if err == nil {
			// get the updated order from the database
			order, err := h.ordersService.GetOrderWithOrderLinesByUid(newUID, facilityCode)

			if err != nil {
				log.Error().Msg(err.Error())
				return echo.ErrInternalServerError
			}

			return c.JSON(http.StatusOK, order)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}

	}
}

// UpdateOrder godoc
// @Summary Update an order
// @Description Update an order
// @Tags Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param order body models.OrderDetail true "Order object that needs to be updated"
// @Success 200 {object} models.OrderDetail
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Router /v1/order [put]
func (h *OrdersHandlers) UpdateOrder() echo.HandlerFunc {

	return func(c echo.Context) error {

		order := new(models.OrderDetail)
		err := c.Bind(order)

		if err != nil {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}

		facilityCode := c.Get("facilityCode").(string)
		userUID := c.Get("userUID").(string)

		err = h.ordersService.UpdateOrder(order, facilityCode, userUID)

		if err == nil {
			// get the updated order from the database
			order, err := h.ordersService.GetOrderWithOrderLinesByUid(order.UID, facilityCode)

			if err != nil {

				log.Error().Msg(err.Error())
				return echo.ErrInternalServerError
			}

			return c.JSON(http.StatusOK, order)

		} else if err == helpers.ERR_CONFLICT {
			return echo.ErrConflict
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}

	}
}

// DeleteOrder godoc
// @Summary Delete an order
// @Description Delete an order by order UID
// @Tags Orders
// @Security BearerAuth
// @Param uid path string true "Order UID"
// @Success 204 "No Content"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Router /v1/order/{uid} [delete]
func (h *OrdersHandlers) DeleteOrder() echo.HandlerFunc {

	return func(c echo.Context) error {

		orderUID := c.Param("uid")
		userUID := c.Get("userUID").(string)

		err := h.ordersService.DeleteOrder(orderUID, userUID)

		if err == nil {
			return c.NoContent(http.StatusNoContent)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}

	}
}

// UpdateOrderLineDelivery godoc
// @Summary Update order line delivery
// @Description Updates delivery status and optional serial/EUN for a single order line item.
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uid path string true "Order UID"
// @Param itemUid path string true "Order line item UID"
// @Param body body models.OrderLineDelivery true "Order line delivery payload"
// @Success 200 {object} models.OrderLine
// @Failure 500 "Internal server error"
// @Router /v1/order/{uid}/orderline/{itemUid}/delivery [put]
func (h *OrdersHandlers) UpdateOrderLineDelivery() echo.HandlerFunc {

	return func(c echo.Context) error {

		itemUID := c.Param("itemUid")
		userUID := c.Get("userUID").(string)
		facilityCode := c.Get("facilityCode").(string)
		// get the serial number as is delivered from the request body as as struct of OrderLineDelivery
		orderLineDelivery := new(models.OrderLineDelivery)
		err := c.Bind(orderLineDelivery)

		if err != nil {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}

		orderLine, err := h.ordersService.UpdateOrderLineDelivery(itemUID, orderLineDelivery.IsDelivered, orderLineDelivery.SerialNumber, orderLineDelivery.EUN, userUID, facilityCode)

		if err == nil {
			return c.JSON(http.StatusOK, orderLine)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

// UpdateMultipleOrderLineDelivery godoc
// @Summary Update multiple order line delivery
// @Description Update multiple order line delivery
// @Tags Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param uid path string true "Order UID"
// @Param itemsUids body []string true "Items UIDs"
// @Success 200 {object} []models.OrderLine
// @Failure 400 "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Router /v1/order/{uid}/orderlines/delivery [put]
func (h *OrdersHandlers) UpdateMultipleOrderLineDelivery() echo.HandlerFunc {
	return func(c echo.Context) error {

		userUID := c.Get("userUID").(string)
		facilityCode := c.Get("facilityCode").(string)
		// get the serial number as is delivered from the request body as as struct of OrderLineDelivery
		itemUids := new([]string)
		err := c.Bind(itemUids)

		if err != nil {
			log.Error().Msg(err.Error())
			return echo.ErrBadRequest
		}

		orderLines, err := h.ordersService.UpdateMultipleOrderLineDelivery(*itemUids, userUID, facilityCode)

		if err == nil {
			return c.JSON(http.StatusOK, orderLines)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

// UpdateServiceLineDelivery godoc
// @Summary Update service line delivery
// @Description Update service line delivery status
// @Tags Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param uid path string true "Order UID"
// @Param serviceItemUid path string true "Service Item UID"
// @Param serviceLineDelivery body models.ServiceLineDelivery true "Service Line Delivery object"
// @Success 200 {object} models.ServiceLine
// @Failure 400 "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Router /v1/order/{uid}/serviceline/{serviceItemUid}/delivery [put]
func (h *OrdersHandlers) UpdateServiceLineDelivery() echo.HandlerFunc {
	return func(c echo.Context) error {
		serviceItemUID := c.Param("serviceItemUid")
		userUID := c.Get("userUID").(string)
		facilityCode := c.Get("facilityCode").(string)

		// get the delivery info from the request body
		serviceLineDelivery := new(models.ServiceLineDelivery)
		err := c.Bind(serviceLineDelivery)

		if err != nil {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}

		serviceLine, err := h.ordersService.UpdateServiceLineDelivery(serviceItemUID, serviceLineDelivery.IsDelivered, userUID, facilityCode)

		if err == nil {
			return c.JSON(http.StatusOK, serviceLine)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

// UpdateMultipleServiceLineDelivery godoc
// @Summary Update multiple service line delivery
// @Description Update multiple service line delivery status to delivered
// @Tags Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param uid path string true "Order UID"
// @Param serviceItemUids body []string true "Service Item UIDs"
// @Success 200 {object} []models.ServiceLine
// @Failure 400 "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Router /v1/order/{uid}/servicelines/delivery [put]
func (h *OrdersHandlers) UpdateMultipleServiceLineDelivery() echo.HandlerFunc {
	return func(c echo.Context) error {
		userUID := c.Get("userUID").(string)
		facilityCode := c.Get("facilityCode").(string)

		// get the service item UIDs from the request body
		serviceItemUids := new([]string)
		err := c.Bind(serviceItemUids)

		if err != nil {
			log.Error().Msg(err.Error())
			return echo.ErrBadRequest
		}

		serviceLines, err := h.ordersService.UpdateMultipleServiceLineDelivery(*serviceItemUids, userUID, facilityCode)

		if err == nil {
			return c.JSON(http.StatusOK, serviceLines)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

// GetItemsForEunPrint godoc
// @Summary Get items for EUN print
// @Description Get items for EUN print
// @Tags Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param euns query string true "EUNs"
// @Success 200 {object} []models.ItemForEunPrint
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Router /v1/orders/eun-for-print [get]
func (h *OrdersHandlers) GetItemsForEunPrint() echo.HandlerFunc {

	return func(c echo.Context) error {

		euns := c.QueryParam("euns")

		eunsArray := []string{}
		json.Unmarshal([]byte(euns), &eunsArray)

		items, err := h.ordersService.GetItemsForEunPrint(eunsArray)

		if err == nil {
			return c.JSON(http.StatusOK, items)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}

	}
}

// SetItemPrintEUN godoc
// @Summary Set EUN print flag
// @Description Sets whether the given EUN should be printed.
// @Tags Orders
// @Produce plain
// @Security BearerAuth
// @Param eun path string true "EUN"
// @Param printEUN query bool true "Print flag"
// @Success 204 "No content"
// @Failure 500 "Internal server error"
// @Router /v1/orders/eun-for-print/{eun} [put]
func (h *OrdersHandlers) SetItemPrintEUN() echo.HandlerFunc {

	return func(c echo.Context) error {

		eun := c.Param("eun")
		printEUN := c.QueryParam("printEUN")

		printEUNBool := false
		json.Unmarshal([]byte(printEUN), &printEUNBool)

		err := h.ordersService.SetItemPrintEUN(eun, printEUNBool)

		if err == nil {
			return c.NoContent(http.StatusNoContent)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}

	}
}

// GetOrderUidByOrderNumber godoc
// @Summary Get order UID by order number
// @Description Returns order UID for the provided order number.
// @Tags Orders
// @Produce json
// @Param orderNumber path string true "Order number"
// @Success 200 {string} string
// @Failure 500 "Internal server error"
// @Router /v1/order-uid-by-order-number/{orderNumber} [get]
func (h *OrdersHandlers) GetOrderUidByOrderNumber() echo.HandlerFunc {

	return func(c echo.Context) error {

		orderNumber := c.Param("orderNumber")

		orderUid, err := h.ordersService.GetOrderUidByOrderNumber(orderNumber)

		if err == nil {
			return c.JSON(http.StatusOK, orderUid)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}

	}
}

// GetOrdersForCatalogueItem godoc
// @Summary Get orders for catalogue item
// @Description Returns orders that reference the specified catalogue item.
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Param catalogueItemUid path string true "Catalogue item UID"
// @Success 200 {array} models.OrderListItem
// @Failure 500 "Internal server error"
// @Router /v1/catalogue/{catalogueItemUid}/orders [get]
func (h *OrdersHandlers) GetOrdersForCatalogueItem() echo.HandlerFunc {

	return func(c echo.Context) error {

		catalogueItemUid := c.Param("catalogueItemUid")
		facilityCode := c.Get("facilityCode").(string)

		orders, err := h.ordersService.GetOrdersForCatalogueItem(catalogueItemUid, facilityCode)

		if err == nil {
			return c.JSON(http.StatusOK, orders)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}

	}
}

// GetMinAndMaxOrderLinePrice godoc
// @Summary Get min/max order line prices
// @Description Returns minimum and maximum order line price for the current facility.
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.OrderLineMinMax
// @Failure 500 "Internal server error"
// @Router /v1/orders/order-lines/min-max-prices [get]
func (h *OrdersHandlers) GetMinAndMaxOrderLinePrice() echo.HandlerFunc {

	return func(c echo.Context) error {

		facilityCode := c.Get("facilityCode").(string)

		minAndMax, err := h.ordersService.GetMinAndMaxOrderLinePrice(facilityCode)

		if err == nil {
			return c.JSON(http.StatusOK, minAndMax)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}

	}
}
