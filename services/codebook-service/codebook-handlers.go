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

		//has the user role CODEBOOKS_ADMIN ?
		isCodebooksAdmin := contains(roles, shared.ROLE_CODEBOOKS_ADMIN)
		isAdmin := contains(roles, shared.ROLE_ADMIN)

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

		codebookResponse, err := h.codebookService.GetCodebook(codebookCode, searchText, parentUID, limit, facilityCode, filterObject, isAdmin)

		if err == nil {
			return c.JSON(http.StatusOK, codebookResponse)
		}

		return echo.ErrInternalServerError
	}
}

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

func (h *CodebookHandlers) GetListOfCodebooks() echo.HandlerFunc {

	return func(c echo.Context) error {

		editable := c.QueryParams().Get("editable")
		onlyEditable := false
		if editable == "true" {
			onlyEditable = true
		}

		if onlyEditable {
			codebookList := h.codebookService.GetListOfEditableCodebooks()
			return c.JSON(http.StatusOK, codebookList)
		} else {
			codebookList := h.codebookService.GetListOfCodebooks()
			return c.JSON(http.StatusOK, codebookList)
		}
	}
}

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
