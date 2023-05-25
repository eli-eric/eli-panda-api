package codebookService

import (
	"encoding/json"
	"net/http"
	"panda/apigateway/helpers"
	"panda/apigateway/services/codebook-service/models"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type CodebookHandlers struct {
	codebookService ICodebookService
}

type ICodebookHandlers interface {
	GetCodebook() echo.HandlerFunc
	GetListOfCodebooks() echo.HandlerFunc
	CreateNewCodebook() echo.HandlerFunc
}

// NewCommentsHandlers Comments handlers constructor
func NewCodebookHandlers(codebookService ICodebookService) ICodebookHandlers {
	return &CodebookHandlers{codebookService: codebookService}
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

		filterObject := new([]helpers.Filter)
		if filter != "" {
			json.Unmarshal([]byte(filter), filterObject)
		}

		limit := autocompleteDefaultLimit
		limit, err := strconv.Atoi(limitParam)

		if err != nil {
			limit = autocompleteDefaultLimit
		} else if limit > autocompleteMaxLimit {
			limit = autocompleteMaxLimit
		}

		codebookResponse, err := h.codebookService.GetCodebook(codebookCode, searchText, parentUID, limit, facilityCode, filterObject)

		if err == nil {
			return c.JSON(http.StatusOK, codebookResponse)
		}

		return echo.ErrInternalServerError
	}
}

const autocompleteMaxLimit int = 100
const autocompleteDefaultLimit int = 10

func (h *CodebookHandlers) GetListOfCodebooks() echo.HandlerFunc {

	return func(c echo.Context) error {

		codebookList := h.codebookService.GetListOfCodebooks()

		return c.JSON(http.StatusOK, codebookList)
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
