package publicationsservice

import (
	//"panda/apigateway/services/publications-service/models"

	"encoding/json"
	"panda/apigateway/helpers"
	"panda/apigateway/services/publications-service/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type PublicationsHandlers struct {
	PublicationsService IPublicationsService
}

type IPublicationsHandlers interface {
	CreatePublication() echo.HandlerFunc
	GetPublication() echo.HandlerFunc
	GetPublications() echo.HandlerFunc
	UpdatePublication() echo.HandlerFunc
	DeletePublication() echo.HandlerFunc
}

// NewPublicationsHandlers General handlers constructor
func NewPublicationsHandlers(svc IPublicationsService) IPublicationsHandlers {
	return &PublicationsHandlers{PublicationsService: svc}
}

// CreatePublication Create publication godoc
// @Summary Create publication
// @Description Create publication
// @Tags Publications
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param publication body models.Publication true "Publication"
// @Success 200 {object} models.Publication
// @Failure 500 "Internal Server Error"
// @Router /v1/publication [post]
func (h *PublicationsHandlers) CreatePublication() echo.HandlerFunc {

	return func(c echo.Context) error {

		publication := new(models.Publication)
		if err := c.Bind(publication); err != nil {
			log.Error().Err(err).Msg("Error binding publication")
			return echo.ErrBadRequest
		}

		userUID := c.Get("userUID").(string)

		if publication.Uid == "" {
			publication.Uid = uuid.New().String()
		}

		_, err := h.PublicationsService.CreatePublication(publication, userUID)
		if err != nil {
			log.Error().Err(err).Msg("Error creating publication")
			return echo.ErrInternalServerError
		}

		return c.JSON(200, publication)
	}
}

// GetPublication Get publication by uid godoc
// @Summary Get publication by uid
// @Description Get publication by uid
// @Tags Publications
// @Security BearerAuth
// @Produce json
// @Param uid path string true "uid"
// @Success 200 {object} models.Publication
// @Failure 500 "Internal Server Error"
// @Router /v1/publication/{uid} [get]
func (h *PublicationsHandlers) GetPublication() echo.HandlerFunc {

	return func(c echo.Context) error {

		uid := c.Param("uid")

		publication, err := h.PublicationsService.GetPublicationByUid(uid)
		if err != nil {
			// return 404 if not found - in error message will be result contains no more records
			if err.Error() == "Result contains no more records" {
				return echo.ErrNotFound
			}
			log.Error().Err(err).Msg("Error getting publication")
			return echo.ErrInternalServerError
		}

		return c.JSON(200, publication)
	}
}

// GetPublications Get publications godoc
// @Summary Get publications
// @Description Get publications
// @Tags Publications
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Publication
// @Failure 500 "Internal Server Error"
// @Router /v1/publications [get]
func (h *PublicationsHandlers) GetPublications() echo.HandlerFunc {

	return func(c echo.Context) error {

		search := c.QueryParam("search")
		pagination := c.QueryParam("pagination")
		sorting := c.QueryParam("sorting")

		pagingObject := new(helpers.Pagination)
		json.Unmarshal([]byte(pagination), &pagingObject)

		sortingObject := new([]helpers.Sorting)
		json.Unmarshal([]byte(sorting), &sortingObject)

		filterObject := new([]helpers.ColumnFilter)
		filter := c.QueryParam("columnFilter")
		json.Unmarshal([]byte(filter), &filterObject)

		publications, err := h.PublicationsService.GetPublications(search, pagingObject.Page, pagingObject.PageSize)
		if err != nil {
			log.Error().Err(err).Msg("Error getting publications")
			return echo.ErrInternalServerError
		}

		pagiantionResult := helpers.PaginationResult[models.Publication]{
			TotalCount: int64(len(publications)),
			Data:       publications,
		}

		return c.JSON(200, pagiantionResult)
	}
}

// UpdatePublication Update publication godoc
// @Summary Update publication
// @Description Update publication
// @Tags Publications
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param uid path string true "uid"
// @Param publication body models.Publication true "Publication"
// @Success 200 {object} models.Publication
// @Failure 500 "Internal Server Error"
// @Router /v1/publication/{uid} [put]
func (h *PublicationsHandlers) UpdatePublication() echo.HandlerFunc {

	return func(c echo.Context) error {

		uid := c.Param("uid")

		publication := new(models.Publication)
		if err := c.Bind(publication); err != nil {
			log.Error().Err(err).Msg("Error binding publication")
			return echo.ErrBadRequest
		}

		publication.Uid = uid

		userUID := c.Get("userUID").(string)

		_, err := h.PublicationsService.UpdatePublication(publication, userUID)
		if err != nil {
			log.Error().Err(err).Msg("Error updating publication")
			return echo.ErrInternalServerError
		}

		return c.JSON(200, publication)
	}
}

// DeletePublication Delete publication by uid godoc
// @Summary Delete publication by uid
// @Description Delete publication by uid
// @Tags Publications
// @Security BearerAuth
// @Produce json
// @Param uid path string true "uid"
// @Success 204 "No Content"
// @Failure 500 "Internal Server Error"
// @Router /v1/publication/{uid} [delete]
func (h *PublicationsHandlers) DeletePublication() echo.HandlerFunc {

	return func(c echo.Context) error {

		uid := c.Param("uid")
		userUID := c.Get("userUID").(string)

		err := h.PublicationsService.DeletePublication(uid, userUID)
		if err != nil {
			log.Error().Err(err).Msg("Error deleting publication")
			return echo.ErrInternalServerError
		}

		return c.NoContent(204)
	}
}
