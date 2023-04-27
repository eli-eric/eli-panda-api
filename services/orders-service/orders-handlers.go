package ordersService

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"panda/apigateway/helpers"

	"github.com/labstack/echo/v4"
)

type OrdersHandlers struct {
	ordersService IOrdersService
}

type IOrdersHandlers interface {
	GetOrders() echo.HandlerFunc
	GetOrderStatusesCodebook() echo.HandlerFunc
	GetOrdersWithSearchAndPagination() echo.HandlerFunc
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

func (h *OrdersHandlers) GetOrders() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get all categories by parent path
		//orders, err := h.ordersService.GetOrders()

		pagination := c.QueryParam("pagination")
		sorting := c.QueryParam("sorting")
		search := c.QueryParam("search")

		pagingObject := new(helpers.Pagination)
		json.Unmarshal([]byte(pagination), &pagingObject)

		sortingObject := new([]helpers.Sorting)
		json.Unmarshal([]byte(sorting), &sortingObject)

		fmt.Println("page: ", pagingObject.Page, "page size: ", pagingObject.PageSize)
		fmt.Println(sortingObject)
		fmt.Println(search)

		return c.JSON(http.StatusOK, "ok")

		//return echo.ErrInternalServerError
	}
}

func (h *OrdersHandlers) GetOrdersWithSearchAndPagination() echo.HandlerFunc {

	return func(c echo.Context) error {

		pagination := c.QueryParam("pagination")
		sorting := c.QueryParam("sorting")
		search := c.QueryParam("search")

		pagingObject := new(helpers.Pagination)
		json.Unmarshal([]byte(pagination), &pagingObject)

		sortingObject := new([]helpers.Sorting)
		json.Unmarshal([]byte(sorting), &sortingObject)

		items, err := h.ordersService.GetOrdersWithSearchAndPagination(search, pagingObject, sortingObject)

		if err == nil {
			return c.JSON(http.StatusOK, items)
		} else {
			log.Println(err)
			return echo.ErrInternalServerError
		}

	}
}
