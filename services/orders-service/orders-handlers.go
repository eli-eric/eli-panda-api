package ordersService

import (
	"encoding/json"
	"net/http"
	"panda/apigateway/helpers"
	"panda/apigateway/services/orders-service/models"

	"github.com/rs/zerolog/log"

	"github.com/labstack/echo/v4"
)

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

func (h *OrdersHandlers) GetOrdersWithSearchAndPagination() echo.HandlerFunc {

	return func(c echo.Context) error {

		pagination := c.QueryParam("pagination")
		sorting := c.QueryParam("sorting")
		search := c.QueryParam("search")
		supplierUID := c.QueryParam("supplierUID")
		orderStatusUID := c.QueryParam("orderStatusUID")
		procurementResponsibleUID := c.QueryParam("procurementResponsibleUID")
		requestorUID := c.QueryParam("requestorUID")
		facilityCode := c.Get("facilityCode").(string)

		pagingObject := new(helpers.Pagination)
		json.Unmarshal([]byte(pagination), &pagingObject)

		sortingObject := new([]helpers.Sorting)
		json.Unmarshal([]byte(sorting), &sortingObject)

		items, err := h.ordersService.GetOrdersWithSearchAndPagination(search, supplierUID, orderStatusUID, procurementResponsibleUID, requestorUID, facilityCode, pagingObject, sortingObject)

		if err == nil {
			return c.JSON(http.StatusOK, items)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}

	}
}

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
// @Router /v1/order/items/eun/print [get]
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
