package systemsService

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"panda/apigateway/helpers"
	"panda/apigateway/services/systems-service/models"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// patchErrorResponse maps PATCH-family service errors to HTTP responses.
func patchErrorResponse(c echo.Context, err error, op string) error {
	var validation *ValidationError
	switch {
	case errors.As(err, &validation):
		return helpers.BadRequest(validation.Msg)
	case errors.Is(err, helpers.ERR_NOT_FOUND), errors.Is(err, ErrNoPhysicalItem), errors.Is(err, ErrRelationNotFound):
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	case errors.Is(err, ErrSystemTypeMissing):
		return helpers.BadRequest(err.Error())
	case errors.Is(err, ErrStaleSystem):
		return c.JSON(http.StatusConflict, models.ConflictResponse{Code: "STALE", Message: err.Error()})
	case errors.Is(err, helpers.ERR_DUPLICATE_SYSTEM_CODE):
		return c.JSON(http.StatusConflict, models.ConflictResponse{Code: "DUPLICATE_SYSTEM_CODE", Message: "duplicate system code"})
	}
	log.Error().Err(err).Msg(op)
	return echo.ErrInternalServerError
}

func readRawJSONBody(c echo.Context) (map[string]json.RawMessage, error) {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return nil, errors.New("cannot read request body")
	}
	raw := map[string]json.RawMessage{}
	if len(body) == 0 {
		return raw, nil
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, errors.New("invalid JSON body")
	}
	return raw, nil
}

// GetSystemDetail godoc
// @Summary Get system detail
// @Description Returns the full system detail (codebooks, people, physical item, statistics, concurrency token).
// @Tags Systems
// @Produce json
// @Security BearerAuth
// @Param uid path string true "System UID"
// @Success 200 {object} models.SystemDetail
// @Failure 404 "System not found"
// @Failure 500 "Internal server error"
// @Router /v1/system/{uid} [get]
func (h *SystemsHandlers) GetSystemDetail() echo.HandlerFunc {
	return func(c echo.Context) error {
		detail, err := h.systemsService.GetSystemDetail(c.Param("uid"), c.Get("facilityCode").(string))
		if err != nil {
			return patchErrorResponse(c, err, "GetSystemDetail")
		}
		return c.JSON(http.StatusOK, detail)
	}
}

// PatchSystem godoc
// @Summary Patch system
// @Description Partial update (absent key = unchanged, null = clear). Optional lastUpdateTime token; codebooks as {uid}; operators/maintainedBy replace the whole list. Writes the WAS_UPDATED_BY audit edge with the change delta.
// @Tags Systems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uid path string true "System UID"
// @Param body body object true "Partial system fields"
// @Success 200 {object} models.SystemDetail
// @Failure 400 "Bad request"
// @Failure 404 "System not found"
// @Failure 409 {object} models.ConflictResponse "STALE or DUPLICATE_SYSTEM_CODE"
// @Failure 500 "Internal server error"
// @Router /v1/system/{uid} [patch]
func (h *SystemsHandlers) PatchSystem() echo.HandlerFunc {
	return func(c echo.Context) error {
		raw, err := readRawJSONBody(c)
		if err != nil {
			return helpers.BadRequest(err.Error())
		}
		fields, err := parsePatchSystemPayload(raw)
		if err != nil {
			return helpers.BadRequest(err.Error())
		}
		detail, err := h.systemsService.PatchSystem(c.Param("uid"), c.Get("facilityCode").(string), c.Get("userUID").(string), fields)
		if err != nil {
			return patchErrorResponse(c, err, "PatchSystem")
		}
		return c.JSON(http.StatusOK, detail)
	}
}

// PatchSystemPhysicalItem godoc
// @Summary Patch the physical item of a system
// @Description Partial update of the contained item (serialNumber, notes, itemUsage, conditionStatus). Uses the system lastUpdateTime token; audit is recorded on the system.
// @Tags Systems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uid path string true "System UID"
// @Param body body object true "Partial physical item fields"
// @Success 200 {object} models.SystemDetail
// @Failure 400 "Bad request"
// @Failure 404 "System or physical item not found"
// @Failure 409 {object} models.ConflictResponse "STALE"
// @Failure 500 "Internal server error"
// @Router /v1/system/{uid}/physical-item [patch]
func (h *SystemsHandlers) PatchSystemPhysicalItem() echo.HandlerFunc {
	return func(c echo.Context) error {
		raw, err := readRawJSONBody(c)
		if err != nil {
			return helpers.BadRequest(err.Error())
		}
		fields, err := parsePatchPhysicalItemPayload(raw)
		if err != nil {
			return helpers.BadRequest(err.Error())
		}
		detail, err := h.systemsService.PatchSystemPhysicalItem(c.Param("uid"), c.Get("facilityCode").(string), c.Get("userUID").(string), fields)
		if err != nil {
			return patchErrorResponse(c, err, "PatchSystemPhysicalItem")
		}
		return c.JSON(http.StatusOK, detail)
	}
}

// CreateSystem godoc
// @Summary Create system
// @Description Creates a root system, or a child of parentUid (absent responsible/location/zone are inherited, parent→child level rules validated). Body uses the PATCH field shapes; name and systemLevel are required. Writes the INSERT audit edge.
// @Tags Systems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body object true "System fields (+ optional parentUid)"
// @Success 201 {object} models.SystemDetail
// @Failure 400 "Bad request"
// @Failure 404 "Parent not found"
// @Failure 409 {object} models.ConflictResponse "DUPLICATE_SYSTEM_CODE"
// @Failure 500 "Internal server error"
// @Router /v1/systems [post]
func (h *SystemsHandlers) CreateSystem() echo.HandlerFunc {
	return func(c echo.Context) error {
		raw, err := readRawJSONBody(c)
		if err != nil {
			return helpers.BadRequest(err.Error())
		}
		fields, err := parseCreateSystemPayload(raw)
		if err != nil {
			return helpers.BadRequest(err.Error())
		}
		detail, err := h.systemsService.CreateSystem(c.Get("facilityCode").(string), c.Get("userUID").(string), fields)
		if err != nil {
			return patchErrorResponse(c, err, "CreateSystem")
		}
		return c.JSON(http.StatusCreated, detail)
	}
}

// CreateSubsystem godoc
// @Summary Create subsystem
// @Description Creates a child system; responsible, location and zone are inherited from the parent and the parent→child level rules are validated.
// @Tags Systems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uid path string true "Parent system UID"
// @Param body body models.CreateSubsystemRequest true "Subsystem"
// @Success 201 {object} models.SystemDetail
// @Failure 400 "Bad request"
// @Failure 404 "Parent not found"
// @Failure 500 "Internal server error"
// @Router /v1/system/{uid}/subsystems [post]
func (h *SystemsHandlers) CreateSubsystem() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(models.CreateSubsystemRequest)
		if err := c.Bind(req); err != nil {
			return helpers.BadRequest(err.Error())
		}
		detail, err := h.systemsService.CreateSubsystem(c.Param("uid"), c.Get("facilityCode").(string), c.Get("userUID").(string), req)
		if err != nil {
			return patchErrorResponse(c, err, "CreateSubsystem")
		}
		return c.JSON(http.StatusCreated, detail)
	}
}

// GenerateSystemCode godoc
// @Summary Generate and save system code
// @Description Generates a code from the system's type/zone/location/parent and saves it atomically (token check, uniqueness, audit).
// @Tags Systems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uid path string true "System UID"
// @Param body body models.SystemCodeGenerateRequest false "Optional token"
// @Success 200 {object} models.SystemDetail
// @Failure 400 "System has no system type"
// @Failure 404 "System not found"
// @Failure 409 {object} models.ConflictResponse "STALE or DUPLICATE_SYSTEM_CODE"
// @Failure 500 "Internal server error"
// @Router /v1/system/{uid}/system-code [post]
func (h *SystemsHandlers) GenerateSystemCode() echo.HandlerFunc {
	return func(c echo.Context) error {
		raw, err := readRawJSONBody(c)
		if err != nil {
			return helpers.BadRequest(err.Error())
		}
		token, err := parsePatchToken(raw)
		if err != nil {
			return helpers.BadRequest(err.Error())
		}
		detail, err := h.systemsService.GenerateSystemCode(c.Param("uid"), c.Get("facilityCode").(string), c.Get("userUID").(string), token)
		if err != nil {
			return patchErrorResponse(c, err, "GenerateSystemCode")
		}
		return c.JSON(http.StatusOK, detail)
	}
}

// DeleteRelationship godoc
// @Summary Delete system relationship
// @Description Removes (uid)-[type]->(otherUid). Only functional relationship types; IS_SPARE_FOR also recalculates affected coverages. Audited on both systems.
// @Tags Systems
// @Security BearerAuth
// @Param uid path string true "Start system UID"
// @Param type path string true "Relationship type"
// @Param otherUid path string true "End system UID"
// @Success 204 "No content"
// @Failure 400 "Relationship type not removable"
// @Failure 404 "Relationship not found"
// @Failure 500 "Internal server error"
// @Router /v1/system/{uid}/relationships/{type}/{otherUid} [delete]
func (h *SystemsHandlers) DeleteRelationship() echo.HandlerFunc {
	return func(c echo.Context) error {
		err := h.systemsService.DeleteRelationship(c.Param("uid"), c.Param("type"), c.Param("otherUid"), c.Get("facilityCode").(string), c.Get("userUID").(string))
		if err != nil {
			return patchErrorResponse(c, err, "DeleteRelationship")
		}
		return c.NoContent(http.StatusNoContent)
	}
}

// GetSystemSpareParts godoc
// @Summary Get system spare parts
// @Description Spare parts for the system (inbound IS_SPARE_FOR) and systems it is a spare for, with coverage.
// @Tags Systems
// @Produce json
// @Security BearerAuth
// @Param uid path string true "System UID"
// @Success 200 {object} models.SystemSparePartsResponse
// @Failure 404 "System not found"
// @Failure 500 "Internal server error"
// @Router /v1/system/{uid}/spare-parts [get]
func (h *SystemsHandlers) GetSystemSpareParts() echo.HandlerFunc {
	return func(c echo.Context) error {
		result, err := h.systemsService.GetSystemSpareParts(c.Param("uid"), c.Get("facilityCode").(string))
		if err != nil {
			return patchErrorResponse(c, err, "GetSystemSpareParts")
		}
		return c.JSON(http.StatusOK, result)
	}
}

// GetPhysicalItemCatalogueProperties godoc
// @Summary Get resolved catalogue properties of a physical item
// @Description Catalogue properties grouped by property group with the latest service-item overrides applied.
// @Tags Systems
// @Produce json
// @Security BearerAuth
// @Param uid path string true "Physical item UID"
// @Success 200 {object} models.PhysicalItemCatalogueProperties
// @Failure 404 "Item not found"
// @Failure 500 "Internal server error"
// @Router /v1/physical-item/{uid}/catalogue-properties [get]
func (h *SystemsHandlers) GetPhysicalItemCatalogueProperties() echo.HandlerFunc {
	return func(c echo.Context) error {
		result, err := h.systemsService.GetPhysicalItemCatalogueProperties(c.Param("uid"), c.Get("facilityCode").(string))
		if err != nil {
			return patchErrorResponse(c, err, "GetPhysicalItemCatalogueProperties")
		}
		return c.JSON(http.StatusOK, result)
	}
}

// GetPhysicalItemOrders godoc
// @Summary Get orders of a physical item
// @Description The order the item was delivered by and its service items (newest first).
// @Tags Systems
// @Produce json
// @Security BearerAuth
// @Param uid path string true "Physical item UID"
// @Success 200 {object} models.PhysicalItemOrders
// @Failure 404 "Item not found"
// @Failure 500 "Internal server error"
// @Router /v1/physical-item/{uid}/orders [get]
func (h *SystemsHandlers) GetPhysicalItemOrders() echo.HandlerFunc {
	return func(c echo.Context) error {
		result, err := h.systemsService.GetPhysicalItemOrders(c.Param("uid"), c.Get("facilityCode").(string))
		if err != nil {
			return patchErrorResponse(c, err, "GetPhysicalItemOrders")
		}
		return c.JSON(http.StatusOK, result)
	}
}
