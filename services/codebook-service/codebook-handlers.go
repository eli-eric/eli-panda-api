package codebookService

import (
	"encoding/json"
	"net/http"
	"panda/apigateway/helpers"
	"panda/apigateway/services/codebook-service/models"
	"panda/apigateway/shared"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/rs/zerolog/log"
)

type CodebookHandlers struct {
	codebookService ICodebookService
}

type ICodebookHandlers interface {
	GetCodebook() echo.HandlerFunc
	GetCodebookTree() echo.HandlerFunc
	GetListOfCodebooks() echo.HandlerFunc
	CreateNewCodebook() echo.HandlerFunc
	UpdateCodebook() echo.HandlerFunc
	DeleteCodebook() echo.HandlerFunc
}

// NewCommentsHandlers Comments handlers constructor
func NewCodebookHandlers(codebookService ICodebookService) ICodebookHandlers {
	return &CodebookHandlers{codebookService: codebookService}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// GetCodebook godoc
// @Summary Get codebook
// @Description Returns codebook items and metadata for the given codebook code.
// @Tags Codebooks
// @Produce json
// @Security BearerAuth
// @Param codebookCode path string true "Codebook code"
// @Param parentUID query string false "Parent UID filter"
// @Param searchText query string false "Search text"
// @Param limit query int false "Max number of items"
// @Param filter query string false "Filter JSON"
// @Success 200 {object} models.CodebookResponse
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /v1/codebook/{codebookCode} [get]
func (h *CodebookHandlers) GetCodebook() echo.HandlerFunc {

	return func(c echo.Context) error {
		//get query path param
		codebookCode := c.Param("codebookCode")
		parentUID := c.QueryParams().Get("parentUID")
		searchText := c.QueryParams().Get("searchText")
		limitParam := c.QueryParams().Get("limit")
		facilityCode := c.Get("facilityCode").(string)
		filter := c.QueryParams().Get("filter")
		roles := c.Get("userRoles").([]string)
		userUID := c.Get("userUID").(string)

		//has the user role CODEBOOKS_ADMIN ?
		isCodebooksAdmin := contains(roles, shared.ROLE_CODEBOOKS_ADMIN)

		filterObject := new([]helpers.Filter)
		if filter != "" {
			json.Unmarshal([]byte(filter), filterObject)
		}

		limit := autocompleteDefaultLimit
		limit, err := strconv.Atoi(limitParam)

		if err != nil {
			limit = autocompleteDefaultLimit
		} else if limit > autocompleteMaxLimit && !isCodebooksAdmin {
			limit = autocompleteMaxLimit
		}

		codebookResponse, err := h.codebookService.GetCodebook(codebookCode, searchText, parentUID, limit, facilityCode, filterObject, userUID, roles)

		if err == nil {
			return c.JSON(http.StatusOK, codebookResponse)
		}

		log.Error().Msg(err.Error())
		if err == helpers.ERR_UNAUTHORIZED {
			return echo.ErrUnauthorized
		} else {
			return echo.ErrInternalServerError
		}
	}
}

// GetCodebookTree godoc
// @Summary Get codebook tree
// @Description Returns codebook tree structure for the given codebook code.
// @Tags Codebooks
// @Produce json
// @Security BearerAuth
// @Param codebookCode path string true "Codebook code"
// @Param columnFilter query string false "Column filter JSON"
// @Success 200 {array} models.CodebookTreeItem
// @Failure 500 "Internal server error"
// @Router /v1/codebook/{codebookCode}/tree [get]
func (h *CodebookHandlers) GetCodebookTree() echo.HandlerFunc {

	return func(c echo.Context) error {
		//get query path param
		codebookCode := c.Param("codebookCode")
		facilityCode := c.Get("facilityCode").(string)
		columnFilter := c.QueryParams().Get("columnFilter")

		filterObject := new([]helpers.ColumnFilter)
		if columnFilter != "" {
			json.Unmarshal([]byte(columnFilter), filterObject)
		}

		codebookTree, err := h.codebookService.GetCodebookTree(codebookCode, facilityCode, filterObject)

		if err == nil {
			return c.JSON(http.StatusOK, codebookTree)
		}

		return echo.ErrInternalServerError
	}
}

const autocompleteMaxLimit int = 100
const autocompleteDefaultLimit int = 10

// GetListOfCodebooks godoc
// @Summary Get list of codebooks
// @Description Returns a list of all codebooks (or only editable ones if requested).
// @Tags Codebooks
// @Produce json
// @Security BearerAuth
// @Param editable query bool false "If true, return only editable codebooks"
// @Success 200 {array} models.CodebookType
// @Failure 500 "Internal server error"
// @Router /v1/codebooks [get]
func (h *CodebookHandlers) GetListOfCodebooks() echo.HandlerFunc {

	return func(c echo.Context) error {

		editable := c.QueryParams().Get("editable")
		userRoles := c.Get("userRoles").([]string)
		onlyEditable := false
		if editable == "true" {
			onlyEditable = true
		}

		if onlyEditable {
			codebookList := h.codebookService.GetListOfEditableCodebooks(userRoles)
			return c.JSON(http.StatusOK, codebookList)
		} else {
			codebookList := h.codebookService.GetListOfCodebooks()
			return c.JSON(http.StatusOK, codebookList)
		}
	}
}

// CreateNewCodebook godoc
// @Summary Create new codebook item
// @Description Creates a new item in the specified codebook.
// @Tags Codebooks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param codebookCode path string true "Codebook code"
// @Param body body models.Codebook true "Codebook item"
// @Success 200 {object} models.Codebook
// @Failure 401 "Unauthorized"
// @Failure 409 "Conflict"
// @Failure 500 "Internal server error"
// @Router /v1/codebook/{codebookCode} [post]
func (h *CodebookHandlers) CreateNewCodebook() echo.HandlerFunc {

	return func(c echo.Context) error {

		codebook := new(models.Codebook)
		if err := c.Bind(codebook); err == nil {
			userUID := c.Get("userUID").(string)
			userRoles := c.Get("userRoles").([]string)
			facilityCode := c.Get("facilityCode").(string)
			codebookCode := c.Param("codebookCode")

			// create new codebook
			createdCodebook, err := h.codebookService.CreateNewCodebook(codebookCode, facilityCode, userUID, userRoles, codebook)
			if err != nil {
				if strings.ContainsAny(err.Error(), "already exists") {
					return c.NoContent(http.StatusConflict)
				} else if strings.ContainsAny(err.Error(), "UNAUTHORIZED") {
					return echo.ErrUnauthorized
				} else {
					return err
				}
			}

			return c.JSON(http.StatusOK, createdCodebook)
		} else {
			return echo.ErrInternalServerError
		}
	}
}

// UpdateCodebook godoc
// @Summary Update codebook item
// @Description Updates an existing item in the specified codebook.
// @Tags Codebooks
// @Accept json
// @Security BearerAuth
// @Param codebookCode path string true "Codebook code"
// @Param uid path string true "Codebook item UID"
// @Param body body models.Codebook true "Codebook item"
// @Success 204 "No content"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /v1/codebook/{codebookCode}/{uid} [put]
func (h *CodebookHandlers) UpdateCodebook() echo.HandlerFunc {

	return func(c echo.Context) error {

		codebook := new(models.Codebook)

		err := c.Bind(codebook)
		if err == nil {
			userUID := c.Get("userUID").(string)
			userRoles := c.Get("userRoles").([]string)
			facilityCode := c.Get("facilityCode").(string)
			codebookCode := c.Param("codebookCode")
			codebook.UID = c.Param("uid")

			// create new codebook
			_, err = h.codebookService.UpdateCodebook(codebookCode, facilityCode, userUID, userRoles, codebook)
		}

		if err == nil {
			return c.NoContent(http.StatusNoContent)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

// DeleteCodebook godoc
// @Summary Delete codebook item
// @Description Deletes an item from the specified codebook.
// @Tags Codebooks
// @Security BearerAuth
// @Param codebookCode path string true "Codebook code"
// @Param uid path string true "Codebook item UID"
// @Success 204 "No content"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /v1/codebook/{codebookCode}/{uid} [delete]
func (h *CodebookHandlers) DeleteCodebook() echo.HandlerFunc {

	return func(c echo.Context) error {

		userUID := c.Get("userUID").(string)
		userRoles := c.Get("userRoles").([]string)
		facilityCode := c.Get("facilityCode").(string)
		codebookCode := c.Param("codebookCode")
		codebookUID := c.Param("uid")

		// create new codebook
		err := h.codebookService.DeleteCodebook(codebookCode, facilityCode, userUID, userRoles, codebookUID)

		if err == nil {
			return c.NoContent(http.StatusNoContent)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}
