package codebookService

import (
	"net/http"
	"panda/apigateway/services/codebook-service/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CodebookHandlers struct {
	codebookService ICodebookService
}

type ICodebookHandlers interface {
	GetCodebook() echo.HandlerFunc
	GetAutocompleteCodebook() echo.HandlerFunc
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
		facilityCode := c.Get("facilityCode").(string)
		// get all categories of the given parentPath
		codebookList, err := h.codebookService.GetCodebook(codebookCode, parentUID, facilityCode)

		if err == nil {
			return c.JSON(http.StatusOK, codebookList)
		}

		return echo.ErrInternalServerError
	}
}

const autocompleteMaxLimit int = 100
const autocompleteDefaultLimit int = 10

func (h *CodebookHandlers) GetAutocompleteCodebook() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get query path param
		codebookCode := c.Param("codebookCode")
		searchText := c.QueryParams().Get("searchText")
		limitParam := c.QueryParams().Get("limit")
		facilityCode := c.Get("facilityCode").(string)

		limit := autocompleteDefaultLimit
		limit, err := strconv.Atoi(limitParam)

		if err != nil {
			limit = autocompleteDefaultLimit
		} else if limit > autocompleteMaxLimit {
			limit = autocompleteMaxLimit
		}

		codebookList, err := h.codebookService.GetAutocompleteCodebook(codebookCode, searchText, limit, facilityCode)

		if err == nil {
			return c.JSONPretty(http.StatusOK, codebookList, "")
		}

		return echo.ErrInternalServerError
	}
}

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
				return err
			}

			return c.JSON(http.StatusOK, createdCodebook)
		} else {
			return echo.ErrInternalServerError
		}
	}
}
