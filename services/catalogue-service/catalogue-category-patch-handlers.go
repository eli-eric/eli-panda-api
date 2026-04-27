package catalogueService

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"panda/apigateway/helpers"
	"panda/apigateway/services/catalogue-service/models"
	codebookModels "panda/apigateway/services/codebook-service/models"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// PatchCatalogueCategory godoc
// @Summary Partially update catalogue category
// @Description JSON Merge Patch on category scalar fields (name, code, systemType).
// @Description Image is out of scope — handled via Minio separately.
// @Description Only keys present in the body are modified. systemType null clears the relationship.
// @Tags Catalogue
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uid path string true "Category UID"
// @Param body body object true "Partial category payload"
// @Success 200 {object} models.CatalogueCategory
// @Failure 400 "Bad Request — malformed body or unknown systemType UID"
// @Failure 404 "Not Found — category does not exist"
// @Failure 500 "Internal server error"
// @Router /v1/catalogue/category/{uid} [patch]
func (h *CatalogueHandlers) PatchCatalogueCategory() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Param("uid")

		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return helpers.BadRequest("cannot read request body")
		}

		raw := map[string]json.RawMessage{}
		if err := json.Unmarshal(body, &raw); err != nil {
			return helpers.BadRequest("invalid JSON body")
		}

		fields, err := parsePatchCatalogueCategoryPayload(raw)
		if err != nil {
			return helpers.BadRequest(err.Error())
		}

		userUID := c.Get("userUID").(string)

		updated, err := h.catalogueService.PatchCatalogueCategory(uid, fields, userUID)
		if err == nil {
			return c.JSON(http.StatusOK, updated)
		} else if errors.Is(err, helpers.ERR_NOT_FOUND) {
			return echo.ErrNotFound
		} else if errors.Is(err, ErrPatchValidation) {
			return helpers.BadRequest(err.Error())
		}

		log.Error().Err(err).Msg("Error patching catalogue category")
		return echo.ErrInternalServerError
	}
}

// GetCatalogueCategoryGroup godoc
// @Summary Get a single catalogue category group
// @Description Returns a single group scoped to a category. 404 if the group doesn't belong to the given category.
// @Tags Catalogue
// @Produce json
// @Security BearerAuth
// @Param uid path string true "Category UID"
// @Param gid path string true "Group UID"
// @Success 200 {object} models.CatalogueCategoryPropertyGroup
// @Failure 404 "Not Found"
// @Failure 500 "Internal server error"
// @Router /v1/catalogue/category/{uid}/group/{gid} [get]
func (h *CatalogueHandlers) GetCatalogueCategoryGroup() echo.HandlerFunc {
	return func(c echo.Context) error {
		got, err := h.catalogueService.GetCatalogueCategoryGroup(c.Param("uid"), c.Param("gid"))
		if err == nil {
			return c.JSON(http.StatusOK, got)
		} else if errors.Is(err, helpers.ERR_NOT_FOUND) {
			return echo.ErrNotFound
		}
		log.Error().Err(err).Msg("Error getting catalogue category group")
		return echo.ErrInternalServerError
	}
}

// GetCatalogueCategoryProperty godoc
// @Summary Get a single catalogue category property
// @Description Returns a single property scoped to a category. 404 if the property doesn't belong to this category.
// @Tags Catalogue
// @Produce json
// @Security BearerAuth
// @Param uid path string true "Category UID"
// @Param pid path string true "Property UID"
// @Success 200 {object} models.CatalogueCategoryProperty
// @Failure 404 "Not Found"
// @Failure 500 "Internal server error"
// @Router /v1/catalogue/category/{uid}/property/{pid} [get]
func (h *CatalogueHandlers) GetCatalogueCategoryProperty() echo.HandlerFunc {
	return func(c echo.Context) error {
		got, err := h.catalogueService.GetCatalogueCategoryProperty(c.Param("uid"), c.Param("pid"))
		if err == nil {
			return c.JSON(http.StatusOK, got)
		} else if errors.Is(err, helpers.ERR_NOT_FOUND) {
			return echo.ErrNotFound
		}
		log.Error().Err(err).Msg("Error getting catalogue category property")
		return echo.ErrInternalServerError
	}
}

// GetCatalogueCategoryPhysicalProperty godoc
// @Summary Get a single physical item property
// @Description Returns a single physical property attached to a category.
// @Tags Catalogue
// @Produce json
// @Security BearerAuth
// @Param uid path string true "Category UID"
// @Param pid path string true "Physical property UID"
// @Success 200 {object} models.CatalogueCategoryProperty
// @Failure 404 "Not Found"
// @Failure 500 "Internal server error"
// @Router /v1/catalogue/category/{uid}/physical-property/{pid} [get]
func (h *CatalogueHandlers) GetCatalogueCategoryPhysicalProperty() echo.HandlerFunc {
	return func(c echo.Context) error {
		got, err := h.catalogueService.GetCatalogueCategoryPhysicalProperty(c.Param("uid"), c.Param("pid"))
		if err == nil {
			return c.JSON(http.StatusOK, got)
		} else if errors.Is(err, helpers.ERR_NOT_FOUND) {
			return echo.ErrNotFound
		}
		log.Error().Err(err).Msg("Error getting physical property")
		return echo.ErrInternalServerError
	}
}

// CreateCatalogueCategoryGroup godoc
// @Summary Create catalogue category group
// @Description Create a new property group under a category. Order is optional — server auto-assigns max(siblings)+10 if omitted.
// @Tags Catalogue
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uid path string true "Category UID"
// @Param body body object true "Group payload (name required)"
// @Success 201 {object} models.CatalogueCategoryPropertyGroup
// @Failure 400 "Bad Request — malformed body or missing name"
// @Failure 404 "Not Found — category does not exist"
// @Failure 500 "Internal server error"
// @Router /v1/catalogue/category/{uid}/group [post]
func (h *CatalogueHandlers) CreateCatalogueCategoryGroup() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Param("uid")

		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return helpers.BadRequest("cannot read request body")
		}
		fields, err := parseCreateCategoryGroupPayload(body)
		if err != nil {
			return helpers.BadRequest(err.Error())
		}

		userUID := c.Get("userUID").(string)
		created, err := h.catalogueService.CreateCatalogueCategoryGroup(uid, fields, userUID)
		if err == nil {
			return c.JSON(http.StatusCreated, created)
		} else if errors.Is(err, helpers.ERR_NOT_FOUND) {
			return echo.ErrNotFound
		}
		log.Error().Err(err).Msg("Error creating catalogue category group")
		return echo.ErrInternalServerError
	}
}

// PatchCatalogueCategoryGroup godoc
// @Summary Update catalogue category group
// @Description Partially update a group's name and/or order. Flat URL — the server validates group belongs to the given category.
// @Tags Catalogue
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uid path string true "Category UID"
// @Param gid path string true "Group UID"
// @Param body body object true "Patch group payload"
// @Success 200 {object} models.CatalogueCategoryPropertyGroup
// @Failure 400 "Bad Request — malformed body"
// @Failure 404 "Not Found — category or group does not exist, or group does not belong to this category"
// @Failure 500 "Internal server error"
// @Router /v1/catalogue/category/{uid}/group/{gid} [patch]
func (h *CatalogueHandlers) PatchCatalogueCategoryGroup() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Param("uid")
		gid := c.Param("gid")

		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return helpers.BadRequest("cannot read request body")
		}
		fields, err := parsePatchCategoryGroupPayload(body)
		if err != nil {
			return helpers.BadRequest(err.Error())
		}

		userUID := c.Get("userUID").(string)
		updated, err := h.catalogueService.PatchCatalogueCategoryGroup(uid, gid, fields, userUID)
		if err == nil {
			return c.JSON(http.StatusOK, updated)
		} else if errors.Is(err, helpers.ERR_NOT_FOUND) {
			return echo.ErrNotFound
		}
		log.Error().Err(err).Msg("Error patching catalogue category group")
		return echo.ErrInternalServerError
	}
}

// DeleteCatalogueCategoryGroup godoc
// @Summary Delete catalogue category group
// @Description Delete a group and its properties. Blocked with 409 if any property under this group is referenced by a catalogue item.
// @Tags Catalogue
// @Security BearerAuth
// @Param uid path string true "Category UID"
// @Param gid path string true "Group UID"
// @Success 204 "No Content"
// @Failure 404 "Not Found"
// @Failure 409 "Conflict — at least one property in this group has item values"
// @Failure 500 "Internal server error"
// @Router /v1/catalogue/category/{uid}/group/{gid} [delete]
func (h *CatalogueHandlers) DeleteCatalogueCategoryGroup() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Param("uid")
		gid := c.Param("gid")
		userUID := c.Get("userUID").(string)

		err := h.catalogueService.DeleteCatalogueCategoryGroup(uid, gid, userUID)
		if err == nil {
			return c.NoContent(http.StatusNoContent)
		} else if errors.Is(err, helpers.ERR_NOT_FOUND) {
			return echo.ErrNotFound
		} else if errors.Is(err, helpers.ERR_DELETE_RELATED_ITEMS) {
			return echo.NewHTTPError(http.StatusConflict, "group contains properties referenced by catalogue items")
		}
		log.Error().Err(err).Msg("Error deleting catalogue category group")
		return echo.ErrInternalServerError
	}
}

// CreateCatalogueCategoryProperty godoc
// @Summary Create catalogue category property
// @Description Create a new property in a group under a category. type.uid required; unit, listOfValues, defaultValue, order are optional.
// @Tags Catalogue
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uid path string true "Category UID"
// @Param gid path string true "Group UID"
// @Param body body object true "Property payload"
// @Success 201 {object} models.CatalogueCategoryProperty
// @Failure 400 "Bad Request — missing/invalid fields or unknown type/unit UID"
// @Failure 404 "Not Found — category or group does not exist"
// @Failure 500 "Internal server error"
// @Router /v1/catalogue/category/{uid}/group/{gid}/property [post]
func (h *CatalogueHandlers) CreateCatalogueCategoryProperty() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Param("uid")
		gid := c.Param("gid")

		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return helpers.BadRequest("cannot read request body")
		}
		fields, err := parseCreateCategoryPropertyPayload(body)
		if err != nil {
			return helpers.BadRequest(err.Error())
		}

		userUID := c.Get("userUID").(string)
		created, err := h.catalogueService.CreateCatalogueCategoryProperty(uid, gid, fields, userUID)
		if err == nil {
			return c.JSON(http.StatusCreated, created)
		} else if errors.Is(err, helpers.ERR_NOT_FOUND) {
			return echo.ErrNotFound
		} else if errors.Is(err, ErrPatchValidation) {
			return helpers.BadRequest(err.Error())
		}
		log.Error().Err(err).Msg("Error creating catalogue category property")
		return echo.ErrInternalServerError
	}
}

// PatchCatalogueCategoryProperty godoc
// @Summary Update catalogue category property
// @Description Partial update. Optional fields: name, defaultValue, listOfValues, order, type, unit, groupUid (move).
// @Tags Catalogue
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uid path string true "Category UID"
// @Param pid path string true "Property UID"
// @Param body body object true "Patch property payload"
// @Success 200 {object} models.CatalogueCategoryProperty
// @Failure 400 "Bad Request — malformed body or unknown ref UID"
// @Failure 404 "Not Found — category or property does not exist, or property does not belong to this category"
// @Failure 500 "Internal server error"
// @Router /v1/catalogue/category/{uid}/property/{pid} [patch]
func (h *CatalogueHandlers) PatchCatalogueCategoryProperty() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Param("uid")
		pid := c.Param("pid")

		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return helpers.BadRequest("cannot read request body")
		}
		fields, err := parsePatchCategoryPropertyPayload(body)
		if err != nil {
			return helpers.BadRequest(err.Error())
		}

		userUID := c.Get("userUID").(string)
		updated, err := h.catalogueService.PatchCatalogueCategoryProperty(uid, pid, fields, userUID)
		if err == nil {
			return c.JSON(http.StatusOK, updated)
		} else if errors.Is(err, helpers.ERR_NOT_FOUND) {
			return echo.ErrNotFound
		} else if errors.Is(err, ErrPatchValidation) {
			return helpers.BadRequest(err.Error())
		}
		log.Error().Err(err).Msg("Error patching catalogue category property")
		return echo.ErrInternalServerError
	}
}

// DeleteCatalogueCategoryProperty godoc
// @Summary Delete catalogue category property
// @Description Delete a property. Blocked with 409 if any catalogue item references the property.
// @Tags Catalogue
// @Security BearerAuth
// @Param uid path string true "Category UID"
// @Param pid path string true "Property UID"
// @Success 204 "No Content"
// @Failure 404 "Not Found"
// @Failure 409 "Conflict — property has item values; clear them first via item PATCH"
// @Failure 500 "Internal server error"
// @Router /v1/catalogue/category/{uid}/property/{pid} [delete]
func (h *CatalogueHandlers) DeleteCatalogueCategoryProperty() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Param("uid")
		pid := c.Param("pid")
		userUID := c.Get("userUID").(string)

		err := h.catalogueService.DeleteCatalogueCategoryProperty(uid, pid, userUID)
		if err == nil {
			return c.NoContent(http.StatusNoContent)
		} else if errors.Is(err, helpers.ERR_NOT_FOUND) {
			return echo.ErrNotFound
		} else if errors.Is(err, helpers.ERR_DELETE_RELATED_ITEMS) {
			return echo.NewHTTPError(http.StatusConflict, "property is referenced by catalogue items")
		}
		log.Error().Err(err).Msg("Error deleting catalogue category property")
		return echo.ErrInternalServerError
	}
}

// CreateCatalogueCategoryPhysicalProperty godoc
// @Summary Create physical item property on a category
// @Description Physical properties are default-value templates attached directly to the category (no group). type.uid required.
// @Tags Catalogue
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uid path string true "Category UID"
// @Param body body object true "Physical property payload"
// @Success 201 {object} models.CatalogueCategoryProperty
// @Failure 400 "Bad Request — missing/invalid fields or unknown type/unit UID"
// @Failure 404 "Not Found — category does not exist"
// @Failure 500 "Internal server error"
// @Router /v1/catalogue/category/{uid}/physical-property [post]
func (h *CatalogueHandlers) CreateCatalogueCategoryPhysicalProperty() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Param("uid")
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return helpers.BadRequest("cannot read request body")
		}
		fields, err := parseCreateCategoryPhysicalPropertyPayload(body)
		if err != nil {
			return helpers.BadRequest(err.Error())
		}
		userUID := c.Get("userUID").(string)
		created, err := h.catalogueService.CreateCatalogueCategoryPhysicalProperty(uid, fields, userUID)
		if err == nil {
			return c.JSON(http.StatusCreated, created)
		} else if errors.Is(err, helpers.ERR_NOT_FOUND) {
			return echo.ErrNotFound
		} else if errors.Is(err, ErrPatchValidation) {
			return helpers.BadRequest(err.Error())
		}
		log.Error().Err(err).Msg("Error creating physical property")
		return echo.ErrInternalServerError
	}
}

// PatchCatalogueCategoryPhysicalProperty godoc
// @Summary Update physical item property
// @Description Partial update — same fields as regular property minus groupUid (physicals don't belong to groups).
// @Tags Catalogue
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uid path string true "Category UID"
// @Param pid path string true "Physical property UID"
// @Param body body object true "Patch payload"
// @Success 200 {object} models.CatalogueCategoryProperty
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Failure 500 "Internal server error"
// @Router /v1/catalogue/category/{uid}/physical-property/{pid} [patch]
func (h *CatalogueHandlers) PatchCatalogueCategoryPhysicalProperty() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Param("uid")
		pid := c.Param("pid")
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return helpers.BadRequest("cannot read request body")
		}
		fields, err := parsePatchCategoryPhysicalPropertyPayload(body)
		if err != nil {
			return helpers.BadRequest(err.Error())
		}
		userUID := c.Get("userUID").(string)
		updated, err := h.catalogueService.PatchCatalogueCategoryPhysicalProperty(uid, pid, fields, userUID)
		if err == nil {
			return c.JSON(http.StatusOK, updated)
		} else if errors.Is(err, helpers.ERR_NOT_FOUND) {
			return echo.ErrNotFound
		} else if errors.Is(err, ErrPatchValidation) {
			return helpers.BadRequest(err.Error())
		}
		log.Error().Err(err).Msg("Error patching physical property")
		return echo.ErrInternalServerError
	}
}

// DeleteCatalogueCategoryPhysicalProperty godoc
// @Summary Delete physical item property
// @Description Physicals aren't referenced by items, so delete never returns 409 — always 204 on success.
// @Tags Catalogue
// @Security BearerAuth
// @Param uid path string true "Category UID"
// @Param pid path string true "Physical property UID"
// @Success 204 "No Content"
// @Failure 404 "Not Found"
// @Failure 500 "Internal server error"
// @Router /v1/catalogue/category/{uid}/physical-property/{pid} [delete]
func (h *CatalogueHandlers) DeleteCatalogueCategoryPhysicalProperty() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Param("uid")
		pid := c.Param("pid")
		userUID := c.Get("userUID").(string)
		err := h.catalogueService.DeleteCatalogueCategoryPhysicalProperty(uid, pid, userUID)
		if err == nil {
			return c.NoContent(http.StatusNoContent)
		} else if errors.Is(err, helpers.ERR_NOT_FOUND) {
			return echo.ErrNotFound
		}
		log.Error().Err(err).Msg("Error deleting physical property")
		return echo.ErrInternalServerError
	}
}

func parseCreateCategoryPhysicalPropertyPayload(body []byte) (*models.CreateCatalogueCategoryPhysicalPropertyFields, error) {
	regular, err := parseCreateCategoryPropertyPayload(body)
	if err != nil {
		return nil, err
	}
	// Physical props share the same input shape as regular ones — just copy fields.
	return &models.CreateCatalogueCategoryPhysicalPropertyFields{
		Name:         regular.Name,
		DefaultValue: regular.DefaultValue,
		ListOfValues: regular.ListOfValues,
		Order:        regular.Order,
		Type:         regular.Type,
		Unit:         regular.Unit,
	}, nil
}

func parsePatchCategoryPhysicalPropertyPayload(body []byte) (*models.PatchCatalogueCategoryPhysicalPropertyFields, error) {
	regular, err := parsePatchCategoryPropertyPayload(body)
	if err != nil {
		return nil, err
	}
	if regular.GroupUID != nil {
		return nil, fmt.Errorf("groupUid is not valid on physical properties — they are not part of any group")
	}
	return &models.PatchCatalogueCategoryPhysicalPropertyFields{
		Name:         regular.Name,
		DefaultValue: regular.DefaultValue,
		ListOfValues: regular.ListOfValues,
		Order:        regular.Order,
		Type:         regular.Type,
		Unit:         regular.Unit,
	}, nil
}

func parseCreateCategoryPropertyPayload(body []byte) (*models.CreateCatalogueCategoryPropertyFields, error) {
	raw := map[string]json.RawMessage{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("invalid JSON body")
	}
	fields := &models.CreateCatalogueCategoryPropertyFields{}

	rname, ok := raw["name"]
	if !ok || rawMessageIsNull(rname) {
		return nil, fmt.Errorf("name is required")
	}
	if err := json.Unmarshal(rname, &fields.Name); err != nil {
		return nil, fmt.Errorf("invalid name: %w", err)
	}
	if fields.Name == "" {
		return nil, fmt.Errorf("name must not be empty")
	}

	rtype, ok := raw["type"]
	if !ok || rawMessageIsNull(rtype) {
		return nil, fmt.Errorf("type is required")
	}
	if err := json.Unmarshal(rtype, &fields.Type); err != nil {
		return nil, fmt.Errorf("invalid type: %w", err)
	}
	if fields.Type.UID == "" {
		return nil, fmt.Errorf("type.uid is required")
	}

	if r, ok := raw["defaultValue"]; ok && !rawMessageIsNull(r) {
		var v string
		if err := json.Unmarshal(r, &v); err != nil {
			return nil, fmt.Errorf("invalid defaultValue: %w", err)
		}
		fields.DefaultValue = &v
	}
	if r, ok := raw["listOfValues"]; ok && !rawMessageIsNull(r) {
		if err := json.Unmarshal(r, &fields.ListOfValues); err != nil {
			return nil, fmt.Errorf("invalid listOfValues: %w", err)
		}
	}
	if r, ok := raw["order"]; ok && !rawMessageIsNull(r) {
		var v int
		if err := json.Unmarshal(r, &v); err != nil {
			return nil, fmt.Errorf("invalid order: %w", err)
		}
		fields.Order = &v
	}
	if r, ok := raw["unit"]; ok && !rawMessageIsNull(r) {
		var cb codebookModels.Codebook
		if err := json.Unmarshal(r, &cb); err != nil {
			return nil, fmt.Errorf("invalid unit: %w", err)
		}
		if cb.UID == "" {
			return nil, fmt.Errorf("unit.uid is required; omit unit to leave it unset")
		}
		fields.Unit = &cb
	}
	return fields, nil
}

func parsePatchCategoryPropertyPayload(body []byte) (*models.PatchCatalogueCategoryPropertyFields, error) {
	raw := map[string]json.RawMessage{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("invalid JSON body")
	}
	fields := &models.PatchCatalogueCategoryPropertyFields{}

	if r, ok := raw["name"]; ok {
		if rawMessageIsNull(r) {
			return nil, fmt.Errorf("name cannot be null")
		}
		var v string
		if err := json.Unmarshal(r, &v); err != nil {
			return nil, fmt.Errorf("invalid name: %w", err)
		}
		fields.Name = &v
	}
	if r, ok := raw["defaultValue"]; ok {
		if rawMessageIsNull(r) {
			fields.DefaultValue = &models.Optional[string]{Value: nil}
		} else {
			var v string
			if err := json.Unmarshal(r, &v); err != nil {
				return nil, fmt.Errorf("invalid defaultValue: %w", err)
			}
			fields.DefaultValue = &models.Optional[string]{Value: &v}
		}
	}
	if r, ok := raw["listOfValues"]; ok {
		if rawMessageIsNull(r) {
			empty := []string{}
			fields.ListOfValues = &empty
		} else {
			var v []string
			if err := json.Unmarshal(r, &v); err != nil {
				return nil, fmt.Errorf("invalid listOfValues: %w", err)
			}
			fields.ListOfValues = &v
		}
	}
	if r, ok := raw["order"]; ok {
		if rawMessageIsNull(r) {
			return nil, fmt.Errorf("order cannot be null")
		}
		var v int
		if err := json.Unmarshal(r, &v); err != nil {
			return nil, fmt.Errorf("invalid order: %w", err)
		}
		fields.Order = &v
	}
	if r, ok := raw["type"]; ok {
		if rawMessageIsNull(r) {
			return nil, fmt.Errorf("type cannot be null")
		}
		var t models.CatalogueCategoryPropertyType
		if err := json.Unmarshal(r, &t); err != nil {
			return nil, fmt.Errorf("invalid type: %w", err)
		}
		if t.UID == "" {
			return nil, fmt.Errorf("type.uid is required")
		}
		fields.Type = &t
	}
	if r, ok := raw["unit"]; ok {
		if rawMessageIsNull(r) {
			fields.Unit = &models.Optional[codebookModels.Codebook]{Value: nil}
		} else {
			var cb codebookModels.Codebook
			if err := json.Unmarshal(r, &cb); err != nil {
				return nil, fmt.Errorf("invalid unit: %w", err)
			}
			if cb.UID == "" {
				return nil, fmt.Errorf("unit.uid is required; to clear unit send unit: null")
			}
			fields.Unit = &models.Optional[codebookModels.Codebook]{Value: &cb}
		}
	}
	if r, ok := raw["groupUid"]; ok {
		if rawMessageIsNull(r) {
			return nil, fmt.Errorf("groupUid cannot be null")
		}
		var v string
		if err := json.Unmarshal(r, &v); err != nil {
			return nil, fmt.Errorf("invalid groupUid: %w", err)
		}
		if v == "" {
			return nil, fmt.Errorf("groupUid must not be empty")
		}
		fields.GroupUID = &v
	}
	return fields, nil
}

func parseCreateCategoryGroupPayload(body []byte) (*models.CreateCatalogueCategoryGroupFields, error) {
	raw := map[string]json.RawMessage{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("invalid JSON body")
	}
	fields := &models.CreateCatalogueCategoryGroupFields{}

	rname, ok := raw["name"]
	if !ok || rawMessageIsNull(rname) {
		return nil, fmt.Errorf("name is required")
	}
	if err := json.Unmarshal(rname, &fields.Name); err != nil {
		return nil, fmt.Errorf("invalid name: %w", err)
	}
	if fields.Name == "" {
		return nil, fmt.Errorf("name must not be empty")
	}

	if r, ok := raw["order"]; ok && !rawMessageIsNull(r) {
		var v int
		if err := json.Unmarshal(r, &v); err != nil {
			return nil, fmt.Errorf("invalid order: %w", err)
		}
		fields.Order = &v
	}
	return fields, nil
}

func parsePatchCategoryGroupPayload(body []byte) (*models.PatchCatalogueCategoryGroupFields, error) {
	raw := map[string]json.RawMessage{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("invalid JSON body")
	}
	fields := &models.PatchCatalogueCategoryGroupFields{}

	if r, ok := raw["name"]; ok {
		if rawMessageIsNull(r) {
			return nil, fmt.Errorf("name cannot be null")
		}
		var v string
		if err := json.Unmarshal(r, &v); err != nil {
			return nil, fmt.Errorf("invalid name: %w", err)
		}
		fields.Name = &v
	}
	if r, ok := raw["order"]; ok {
		if rawMessageIsNull(r) {
			return nil, fmt.Errorf("order cannot be null")
		}
		var v int
		if err := json.Unmarshal(r, &v); err != nil {
			return nil, fmt.Errorf("invalid order: %w", err)
		}
		fields.Order = &v
	}
	return fields, nil
}

// parsePatchCatalogueCategoryPayload maps the raw-message payload into typed fields.
// Absent keys stay nil; explicit null on systemType becomes Optional with nil Value.
func parsePatchCatalogueCategoryPayload(raw map[string]json.RawMessage) (*models.PatchCatalogueCategoryFields, error) {
	fields := &models.PatchCatalogueCategoryFields{}

	if r, ok := raw["name"]; ok {
		if rawMessageIsNull(r) {
			return nil, fmt.Errorf("name cannot be null")
		}
		var v string
		if err := json.Unmarshal(r, &v); err != nil {
			return nil, fmt.Errorf("invalid name: %w", err)
		}
		fields.Name = &v
	}

	if r, ok := raw["code"]; ok {
		if rawMessageIsNull(r) {
			return nil, fmt.Errorf("code cannot be null")
		}
		var v string
		if err := json.Unmarshal(r, &v); err != nil {
			return nil, fmt.Errorf("invalid code: %w", err)
		}
		fields.Code = &v
	}

	if r, ok := raw["systemType"]; ok {
		if rawMessageIsNull(r) {
			fields.SystemType = &models.Optional[codebookModels.Codebook]{Value: nil}
		} else {
			var cb codebookModels.Codebook
			if err := json.Unmarshal(r, &cb); err != nil {
				return nil, fmt.Errorf("invalid systemType: %w", err)
			}
			if cb.UID == "" {
				return nil, fmt.Errorf("systemType.uid is required; to clear the systemType send systemType: null")
			}
			fields.SystemType = &models.Optional[codebookModels.Codebook]{Value: &cb}
		}
	}

	return fields, nil
}

