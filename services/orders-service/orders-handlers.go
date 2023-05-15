package ordersService

import (
	"encoding/json"
	"log"
	"net/http"
	"panda/apigateway/helpers"
	"panda/apigateway/services/orders-service/models"

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
		facilityCode := c.Get("facilityCode").(string)

		pagingObject := new(helpers.Pagination)
		json.Unmarshal([]byte(pagination), &pagingObject)

		sortingObject := new([]helpers.Sorting)
		json.Unmarshal([]byte(sorting), &sortingObject)

		items, err := h.ordersService.GetOrdersWithSearchAndPagination(search, facilityCode, pagingObject, sortingObject)

		if err == nil {
			return c.JSON(http.StatusOK, items)
		} else {
			log.Println(err)
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
			log.Println(err)
			return echo.ErrInternalServerError
		}

	}
}

func (h *OrdersHandlers) InsertNewOrder() echo.HandlerFunc {

	return func(c echo.Context) error {

		order := new(models.OrderDetail)
		err := c.Bind(order)

		if err != nil {
			log.Println(err)
			return echo.ErrInternalServerError
		}

		facilityCode := c.Get("facilityCode").(string)
		userUID := c.Get("userUID").(string)

		newUID, err := h.ordersService.InsertNewOrder(order, facilityCode, userUID)

		if err == nil {
			return c.JSON(http.StatusOK, newUID)
		} else {
			log.Println(err)
			return echo.ErrInternalServerError
		}

	}
}

func (h *OrdersHandlers) UpdateOrder() echo.HandlerFunc {

	return func(c echo.Context) error {

		order := new(models.OrderDetail)
		err := c.Bind(order)

		if err != nil {
			log.Println(err)
			return echo.ErrInternalServerError
		}

		facilityCode := c.Get("facilityCode").(string)
		userUID := c.Get("userUID").(string)

		err = h.ordersService.UpdateOrder(order, facilityCode, userUID)

		if err == nil {
			return c.NoContent(http.StatusNoContent)
		} else {
			log.Println(err)
			return echo.ErrInternalServerError
		}

	}
}

func (h *OrdersHandlers) DeleteOrder() echo.HandlerFunc {

	return func(c echo.Context) error {

		orderUID := c.Param("uid")
		userUID := c.Get("userUID").(string)

		err := h.ordersService.DeleteOrder(orderUID, userUID)

		if err == nil {
			return c.NoContent(http.StatusNoContent)
		} else {
			log.Println(err)
			return echo.ErrInternalServerError
		}

	}
}

func (h *OrdersHandlers) UpdateOrderLineDelivery() echo.HandlerFunc {

	return func(c echo.Context) error {

		itemUID := c.Param("itemUid")
		userUID := c.Get("userUID").(string)
		// get the serial number as is delivered from the request body as as struct of OrderLineDelivery
		orderLineDelivery := new(models.OrderLineDelivery)
		err := c.Bind(orderLineDelivery)

		if err != nil {
			err = h.ordersService.UpdateOrderLineDelivery(itemUID, orderLineDelivery.IsDelivered, orderLineDelivery.SerialNumber, userUID)
		}

		if err == nil {
			return c.NoContent(http.StatusNoContent)
		} else {
			log.Println(err)
			return echo.ErrInternalServerError
		}

	}
}

func (h *OrdersHandlers) GetItemsForEunPrint() echo.HandlerFunc {

	return func(c echo.Context) error {

		euns := c.QueryParam("euns")

		eunsArray := []string{}
		json.Unmarshal([]byte(euns), &eunsArray)

		items, err := h.ordersService.GetItemsForEunPrint(eunsArray)

		if err == nil {
			return c.JSON(http.StatusOK, items)
		} else {
			log.Println(err)
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
			log.Println(err)
			return echo.ErrInternalServerError
		}

	}
}
