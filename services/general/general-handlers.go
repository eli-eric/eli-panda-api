package general

import (
	"encoding/json"
	"net/http"
	"panda/apigateway/helpers"
	"panda/apigateway/services/general/models"

	"github.com/labstack/echo/v4"
)

type GeneralHandlers struct {
	generalService IGeneralService
}

type IGeneralHandlers interface {
	GetGraphByUid() echo.HandlerFunc
	GetUUID() echo.HandlerFunc
	GlobalSearch() echo.HandlerFunc
}

// NewGeneralHandlers General handlers constructor
func NewGeneralHandlers(generalSvc IGeneralService) IGeneralHandlers {
	return &GeneralHandlers{generalService: generalSvc}
}

// GetGraphByUid Get graph by uid godoc
// @Summary Get graph by uid
// @Description Get graph by uid
// @Tags General
// @Security BearerAuth
// @Produce json
// @Param uid path string true "uid"
// @Success 200 {object} models.GraphResponse
// @Failure 500 "Internal Server Error"
// @Router /v1/general/{uid}/graph [get]
func (h *GeneralHandlers) GetGraphByUid() echo.HandlerFunc {

	return func(c echo.Context) error {

		uid := c.Param("uid")

		nodes, err := h.generalService.GetGraphNodesByUid(uid)
		if err != nil {
			return echo.ErrInternalServerError
		}

		links, err := h.generalService.GetGraphLinksByUid(uid)
		if err != nil {
			return echo.ErrInternalServerError
		}

		result := models.GraphResponse{}
		result.Nodes = nodes
		result.Links = links

		return c.JSON(200, result)
	}
}

// GetUUID Get UUID v4 godoc
// @Summary Get UUID V4
// @Description Get UUID v4 string
// @Tags General
// @Security BearerAuth
// @Produce plain
// @Success 200 {string} string
// @Failure 500 "Internal Server Error"
// @Router /v1/uuid/v4 [get]
func (h *GeneralHandlers) GetUUID() echo.HandlerFunc {

	return func(c echo.Context) error {

		uuid, err := h.generalService.GetUUID()
		if err != nil {
			return echo.ErrInternalServerError
		}

		return c.String(200, uuid)
	}
}

// GlobalSearch Global search across systems, orders, and catalogue items godoc
// @Summary Global search
// @Description Search across systems, orders, and catalogue items by free text
// @Tags General
// @Security BearerAuth
// @Produce json
// @Param searchText query string true "Search text"
// @Param pagination query string false "Pagination object as JSON string"
// @Success 200 {object} helpers.PaginationResult[models.GlobalSearchResult]
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /v1/global-search [get]
func (h *GeneralHandlers) GlobalSearch() echo.HandlerFunc {

	return func(c echo.Context) error {

		pagination := c.QueryParam("pagination")
		searchText := c.QueryParam("searchText")

		if searchText == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "searchText parameter is required")
		}

		// Parse pagination
		pagingObject := new(helpers.Pagination)
		if pagination != "" {
			if err := json.Unmarshal([]byte(pagination), &pagingObject); err != nil {
				log.Warn().Err(err).Str("pagination", pagination).Msg("Failed to parse pagination parameter")
			}
		}

		// Set default pagination if not provided or invalid
		if pagingObject.Page == 0 {
			pagingObject.Page = 1
		}
		if pagingObject.PageSize == 0 {
			pagingObject.PageSize = 10 // Default of 10 items when no pagination provided
		}

		// Enforce max limit
		const maxPageSize = 100
		if pagingObject.PageSize > maxPageSize {
			log.Warn().Int("requested", pagingObject.PageSize).Int("max", maxPageSize).Msg("PageSize exceeds maximum, capping to max")
			pagingObject.PageSize = maxPageSize
		}

		log.Debug().Int("page", pagingObject.Page).Int("pageSize", pagingObject.PageSize).Str("searchText", searchText).Msg("GlobalSearch pagination")

		// Call the service method
		result, err := h.generalService.GlobalSearch(searchText, pagingObject.Page, pagingObject.PageSize)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, helpers.PaginationResult[models.GlobalSearchResult]{
			TotalCount: result.TotalCount,
			Data:       result.Data,
		})
	}
}
