package zoneservice

import (
	"encoding/json"
	"errors"
	"net/http"
	"panda/apigateway/helpers"
	"panda/apigateway/services/zone-service/models"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ZoneHandlers struct {
	zoneService IZoneService
}

type IZoneHandlers interface {
	GetAllZones() echo.HandlerFunc
	GetZoneByUID() echo.HandlerFunc
	CreateZone() echo.HandlerFunc
	UpdateZone() echo.HandlerFunc
	DeleteZone() echo.HandlerFunc
	ImportZones() echo.HandlerFunc
}

func NewZoneHandlers(svc IZoneService) IZoneHandlers {
	return &ZoneHandlers{zoneService: svc}
}

// GetAllZones Get all zones godoc
// @Summary Get all zones
// @Description Get all zones for the current facility
// @Tags Zones
// @Security BearerAuth
// @Produce json
// @Success 200 {object} helpers.PaginationResult[models.Zone]
// @Failure 500 "Internal Server Error"
// @Router /v1/zones [get]
func (h *ZoneHandlers) GetAllZones() echo.HandlerFunc {
	return func(c echo.Context) error {
		facilityCode := c.Get("facilityCode").(string)
		search := c.QueryParam("search")

		pagingObject := new(helpers.Pagination)
		pagination := c.QueryParam("pagination")
		json.Unmarshal([]byte(pagination), &pagingObject)

		sortingObject := new([]helpers.Sorting)
		sorting := c.QueryParam("sorting")
		json.Unmarshal([]byte(sorting), &sortingObject)

		zones, totalCount, err := h.zoneService.GetAllZones(facilityCode, search, pagingObject.Page, pagingObject.PageSize, sortingObject)
		if err != nil {
			log.Error().Err(err).Msg("Error getting zones")
			return echo.ErrInternalServerError
		}

		paginationResult := helpers.PaginationResult[models.Zone]{
			TotalCount: totalCount,
			Data:       zones,
		}

		return c.JSON(http.StatusOK, paginationResult)
	}
}

// GetZoneByUID Get zone by uid godoc
// @Summary Get zone by uid
// @Description Get zone by uid
// @Tags Zones
// @Security BearerAuth
// @Produce json
// @Param uid path string true "uid"
// @Success 200 {object} models.Zone
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /v1/zones/{uid} [get]
func (h *ZoneHandlers) GetZoneByUID() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Param("uid")
		facilityCode := c.Get("facilityCode").(string)

		zone, err := h.zoneService.GetZoneByUID(uid, facilityCode)
		if err != nil {
			if err.Error() == "Result contains no more records" {
				return echo.ErrNotFound
			}
			log.Error().Err(err).Msg("Error getting zone")
			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, zone)
	}
}

// CreateZone Create zone godoc
// @Summary Create zone
// @Description Create a new zone
// @Tags Zones
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param zone body models.ZoneCreateRequest true "Zone"
// @Success 201 {object} models.Zone
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /v1/zones [post]
func (h *ZoneHandlers) CreateZone() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(models.ZoneCreateRequest)
		if err := c.Bind(req); err != nil {
			log.Error().Err(err).Msg("Error binding zone create request")
			return helpers.BadRequest(err.Error())
		}

		if req.Name == "" || req.Code == "" {
			return helpers.BadRequest("name and code are required")
		}

		facilityCode := c.Get("facilityCode").(string)
		userUID := c.Get("userUID").(string)

		zone, err := h.zoneService.CreateZone(facilityCode, userUID, req)
		if err != nil {
			if isClientError(err) {
				return helpers.BadRequest(err.Error())
			}
			log.Error().Err(err).Msg("Error creating zone")
			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusCreated, zone)
	}
}

// UpdateZone Update zone godoc
// @Summary Update zone
// @Description Update zone name, code and optionally reassign parent
// @Tags Zones
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param uid path string true "uid"
// @Param zone body models.ZoneUpdateRequest true "Zone"
// @Success 200 {object} models.Zone
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /v1/zones/{uid} [put]
func (h *ZoneHandlers) UpdateZone() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Param("uid")

		req := new(models.ZoneUpdateRequest)
		if err := c.Bind(req); err != nil {
			log.Error().Err(err).Msg("Error binding zone update request")
			return helpers.BadRequest(err.Error())
		}

		if req.Name == "" || req.Code == "" {
			return helpers.BadRequest("name and code are required")
		}

		facilityCode := c.Get("facilityCode").(string)
		userUID := c.Get("userUID").(string)

		zone, err := h.zoneService.UpdateZone(uid, facilityCode, userUID, req)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return echo.ErrNotFound
			}
			if isClientError(err) {
				return helpers.BadRequest(err.Error())
			}
			log.Error().Err(err).Msg("Error updating zone")
			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, zone)
	}
}

// DeleteZone Delete zone godoc
// @Summary Delete zone
// @Description Soft delete zone (rejects if zone has subzones or system references)
// @Tags Zones
// @Security BearerAuth
// @Param uid path string true "uid"
// @Success 200 "OK"
// @Failure 404 "Not Found"
// @Failure 409 "Conflict"
// @Failure 500 "Internal Server Error"
// @Router /v1/zones/{uid} [delete]
func (h *ZoneHandlers) DeleteZone() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Param("uid")
		userUID := c.Get("userUID").(string)
		facilityCode := c.Get("facilityCode").(string)

		err := h.zoneService.DeleteZone(uid, facilityCode, userUID)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return echo.ErrNotFound
			}
			if errors.Is(err, ErrConflictSub) || errors.Is(err, ErrConflictSys) {
				return c.JSON(http.StatusConflict, helpers.ConflictErrorResponse{
					ErrorMessage: err.Error(),
				})
			}
			log.Error().Err(err).Msg("Error deleting zone")
			return echo.ErrInternalServerError
		}

		return c.NoContent(http.StatusOK)
	}
}

// ImportZones Import zones from CSV godoc
// @Summary Import zones from CSV
// @Description Import zones from CSV file. Insert-only, skips existing by code match.
// @Tags Zones
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file"
// @Success 200 {object} models.ZoneImportResult
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /v1/zones/import [post]
func (h *ZoneHandlers) ImportZones() echo.HandlerFunc {
	return func(c echo.Context) error {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			return helpers.BadRequest("file is required")
		}

		file, err := fileHeader.Open()
		if err != nil {
			return helpers.BadRequest("failed to open file")
		}
		defer file.Close()

		facilityCode := c.Get("facilityCode").(string)
		userUID := c.Get("userUID").(string)

		result, err := h.zoneService.ImportZones(facilityCode, userUID, file)
		if err != nil {
			if strings.Contains(err.Error(), "CSV") {
				return helpers.BadRequest(err.Error())
			}
			log.Error().Err(err).Msg("Error importing zones")
			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, result)
	}
}

func isClientError(err error) bool {
	return errors.Is(err, ErrDuplicateCode) ||
		errors.Is(err, ErrSelfParent) ||
		errors.Is(err, ErrParentNotFound) ||
		errors.Is(err, ErrMaxDepth)
}
