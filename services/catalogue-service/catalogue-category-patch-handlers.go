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

