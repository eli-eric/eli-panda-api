package publicationsservice

import (
	//"panda/apigateway/services/publications-service/models"

	"encoding/csv"
	"encoding/json"
	"panda/apigateway/helpers"
	"panda/apigateway/services/publications-service/models"
	"strconv"
	"time"

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
	GetWosDataByDoi() echo.HandlerFunc
	GetPublicationsAsCsv() echo.HandlerFunc
	// Researcher handlers
	GetResearchers() echo.HandlerFunc
	GetResearcher() echo.HandlerFunc
	CreateResearcher() echo.HandlerFunc
	CreateResearchers() echo.HandlerFunc
	UpdateResearcher() echo.HandlerFunc
	DeleteResearcher() echo.HandlerFunc
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
			// return custom bad request with the message
			return helpers.BadRequest(err.Error())
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
// @Param search query string false "search"
// @Param pagination query string false "pagination"
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

		publications, totalCount, err := h.PublicationsService.GetPublications(search, pagingObject.Page, pagingObject.PageSize)
		if err != nil {
			log.Error().Err(err).Msg("Error getting publications")
			return echo.ErrInternalServerError
		}

		pagiantionResult := helpers.PaginationResult[models.Publication]{
			TotalCount: totalCount,
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
			return helpers.BadRequest(err.Error())
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

// GetWosDataByDoi Get WOS data by DOI godoc
// @Summary Get WOS data by DOI
// @Description Get WOS data by DOI
// @Tags Publications
// @Security BearerAuth
// @Produce json
// @Param doi path string true "doi"
// @Success 200  {object} models.WosAPIResponse
// @Failure 500 "Internal Server Error"
// @Router /v1/publication/wos/{doi} [get]
func (h *PublicationsHandlers) GetWosDataByDoi() echo.HandlerFunc {

	return func(c echo.Context) error {

		doi := c.Param("doi")

		result, err := h.PublicationsService.GetPublicationByDoiFromWOS(doi)
		if err != nil {
			log.Error().Err(err).Msg("Error getting WOS data by DOI")
			return echo.ErrInternalServerError
		}

		return c.JSON(200, result)
	}
}

// GetPublicationsAsCsv Get publications as CSV godoc
// @Summary Get publications as CSV
// @Description CSV header: Media Type,Code,Experimental System,User Call,User Experiment,DOI,Web Link,Open Access Type,Title,Authors,Authors Count,ELI Authors,ELI Authors Count,Journal Title,Volume,Issue,Pages,Pages Count,Cite As,Impact Factor,Quartile Basis,Quartile,Year Of Publication,Date Of Publication,Abstract,Keywords,OECD Ford,Grant,WOS Number,ISSN,E-ISSN,EID Scopus,Publishing Country,Language,Note,UID
// @Tags Publications
// @Security BearerAuth
// @Param search query string false "search"
// @Produce text/csv
// @Success 200 "CSV file"
// @Failure 500 "Internal Server Error"
// @Router /v1/publications/export [get]
func (h *PublicationsHandlers) GetPublicationsAsCsv() echo.HandlerFunc {

	return func(c echo.Context) error {

		search := c.QueryParam("search")
		sorting := c.QueryParam("sorting")

		sortingObject := new([]helpers.Sorting)
		json.Unmarshal([]byte(sorting), &sortingObject)

		filterObject := new([]helpers.ColumnFilter)
		filter := c.QueryParam("columnFilter")
		json.Unmarshal([]byte(filter), &filterObject)

		publications, _, err := h.PublicationsService.GetPublications(search, 1, 1_000_000)

		if err != nil {
			log.Error().Err(err).Msg("Error getting publications")
			return echo.ErrInternalServerError
		}

		c.Response().Header().Set(echo.HeaderContentType, "text/csv")
		// file name - will be publications-yyyy-mm-dd-hh-mm-ss.csv
		fileName := "publications-" + time.Now().Format("2006-01-02-15-04-05") + ".csv"
		c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename="+fileName)

		writer := csv.NewWriter(c.Response())
		writer.UseCRLF = true
		writer.Comma = ','

		defer writer.Flush()
		// write header based on System struct
		// write header
		writer.Write([]string{
			"Media Type",
			"Code",
			"Experimental System",
			"User Call",
			"User Experiment",
			"DOI",
			"Web Link",
			"Open Access Type",
			"Title",
			"Authors",
			"Authors Count",
			"ELI Authors",
			"ELI Authors Count",
			"Journal Title",
			"Volume",
			"Issue",
			"Pages",
			"Pages Count",
			"Cite As",
			"Impact Factor",
			"Quartile Basis",
			"Quartile",
			"Year Of Publication",
			"Date Of Publication",
			"Abstract",
			"Keywords",
			"OECD Ford",
			"Grant",
			"WOS Number",
			"ISSN",
			"E-ISSN",
			"EID Scopus",
			"Publishing Country",
			"Language",
			"Note",
			"UID"})

		for _, item := range publications {

			experimentalSystem := ""
			if item.ExperimentalSystemCb != nil {
				experimentalSystem = item.ExperimentalSystemCb.Name
			} else if item.ExperimentalSystem != nil {
				experimentalSystem = *item.ExperimentalSystem
			}
			userCall := ""
			if item.UserCall != nil {
				userCall = item.UserCall.Name
			}
			userExperiment := ""
			if item.UserExperimentCb != nil {
				userExperiment = item.UserExperimentCb.Name
			} else if item.UserExperiment != nil {
				userExperiment = *item.UserExperiment
			}
			openAccessType := ""
			if item.OpenAccessType != nil {
				openAccessType = item.OpenAccessType.Name
			}
			issue := ""
			if item.Issue != nil {
				issue = strconv.Itoa(*item.Issue)
			}
			impactFactor := ""
			if item.ImpactFactor != nil {
				impactFactor = strconv.FormatFloat(*item.ImpactFactor, 'f', -1, 64)
			}
			quartilBasis := ""
			if item.QuartilBasis != nil {
				quartilBasis = *item.QuartilBasis
			}
			quartil := ""
			if item.Quartil != nil {
				quartil = *item.Quartil
			}
			dateOfPublication := ""
			if item.DateOfPublication != nil {
				dateOfPublication = *item.DateOfPublication
			}
			oecdFord := ""
			if item.OecdFord != nil {
				oecdFord = *item.OecdFord
			}
			grant := ""
			if item.GrantCb != nil {
				grant = item.GrantCb.Name
			} else if item.Grant != nil {
				grant = *item.Grant
			}
			wosNumber := ""
			if item.WosNumber != nil {
				wosNumber = *item.WosNumber
			}
			issn := ""
			if item.Issn != nil {
				issn = *item.Issn
			}
			eIssn := ""
			if item.EIssn != nil {
				eIssn = *item.EIssn
			}
			eidScopus := ""
			if item.EidScopus != nil {
				eidScopus = *item.EidScopus
			}
			note := ""
			if item.Note != nil {
				note = *item.Note
			}
			writer.Write([]string{
				item.MediaType,
				item.Code,
				experimentalSystem,
				userCall,
				userExperiment,
				item.Doi,
				item.WebLink,
				openAccessType,
				item.Title,
				item.AllAuthors,
				strconv.Itoa(item.AllAuthorsCount),
				item.EliAuthors,
				strconv.Itoa(item.EliAuthorsCount),
				item.LongJournalTitle,
				strconv.Itoa(item.Volume),
				issue,
				item.Pages,
				strconv.Itoa(item.PagesCount),
				item.CiteAs,
				impactFactor,
				quartilBasis,
				quartil,
				item.YearOfPublication,
				dateOfPublication,
				item.Abstract,
				item.Keywords,
				oecdFord,
				grant,
				wosNumber,
				issn,
				eIssn,
				eidScopus,
				item.PublishingCountry.Name,
				item.Language,
				note,
				item.Uid})
		}

		return nil
	}
}

// Researcher handlers

// GetResearchers Get researchers godoc
// @Summary Get researchers
// @Description Get researchers
// @Tags Researchers
// @Security BearerAuth
// @Produce json
// @Success 200 {object} helpers.PaginationResult[models.Researcher]
// @Failure 500 "Internal Server Error"
// @Router /v1/researchers [get]
// @Param search query string false "search"
// @Param pagination query string false "pagination"
func (h *PublicationsHandlers) GetResearchers() echo.HandlerFunc {

	return func(c echo.Context) error {

		search := c.QueryParam("search")
		pagination := c.QueryParam("pagination")

		pagingObject := new(helpers.Pagination)
		json.Unmarshal([]byte(pagination), &pagingObject)

		researchers, totalCount, err := h.PublicationsService.GetResearchers(search, pagingObject.Page, pagingObject.PageSize)
		if err != nil {
			log.Error().Err(err).Msg("Error getting researchers")
			return echo.ErrInternalServerError
		}

		paginationResult := helpers.PaginationResult[models.Researcher]{
			TotalCount: totalCount,
			Data:       researchers,
		}

		return c.JSON(200, paginationResult)
	}
}

// GetResearcher Get researcher by uid godoc
// @Summary Get researcher by uid
// @Description Get researcher by uid
// @Tags Researchers
// @Security BearerAuth
// @Produce json
// @Param uid path string true "uid"
// @Success 200 {object} models.Researcher
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /v1/researcher/{uid} [get]
func (h *PublicationsHandlers) GetResearcher() echo.HandlerFunc {

	return func(c echo.Context) error {

		uid := c.Param("uid")

		researcher, err := h.PublicationsService.GetResearcherByUid(uid)
		if err != nil {
			// return 404 if not found - in error message will be result contains no more records
			if err.Error() == "Result contains no more records" {
				return echo.ErrNotFound
			}
			log.Error().Err(err).Msg("Error getting researcher")
			return echo.ErrInternalServerError
		}

		return c.JSON(200, researcher)
	}
}

// CreateResearcher Create researcher godoc
// @Summary Create researcher
// @Description Create researcher
// @Tags Researchers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param researcher body models.Researcher true "Researcher"
// @Success 200 {object} models.Researcher
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /v1/researcher [post]
func (h *PublicationsHandlers) CreateResearcher() echo.HandlerFunc {

	return func(c echo.Context) error {

		researcher := new(models.Researcher)
		if err := c.Bind(researcher); err != nil {
			log.Error().Err(err).Msg("Error binding researcher")
			return helpers.BadRequest(err.Error())
		}

		userUID := c.Get("userUID").(string)

		if researcher.Uid == "" {
			researcher.Uid = uuid.New().String()
		}

		_, err := h.PublicationsService.CreateResearcher(researcher, userUID)
		if err != nil {
			log.Error().Err(err).Msg("Error creating researcher")
			return echo.ErrInternalServerError
		}

		return c.JSON(200, researcher)
	}
}

// CreateResearchers Create multiple researchers godoc
// @Summary Create multiple researchers
// @Description Create multiple researchers at once
// @Tags Researchers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param researchers body []models.Researcher true "Researchers"
// @Success 200 {array} models.Researcher
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /v1/researchers [post]
func (h *PublicationsHandlers) CreateResearchers() echo.HandlerFunc {

	return func(c echo.Context) error {

		researchers := new([]models.Researcher)
		if err := c.Bind(researchers); err != nil {
			log.Error().Err(err).Msg("Error binding researchers")
			return helpers.BadRequest(err.Error())
		}

		userUID := c.Get("userUID").(string)

		// Generate UIDs for researchers without one
		for i := range *researchers {
			if (*researchers)[i].Uid == "" {
				(*researchers)[i].Uid = uuid.New().String()
			}
		}

		result, err := h.PublicationsService.CreateResearchers(*researchers, userUID)
		if err != nil {
			log.Error().Err(err).Msg("Error creating researchers")
			return echo.ErrInternalServerError
		}

		return c.JSON(200, result)
	}
}

// UpdateResearcher Update researcher godoc
// @Summary Update researcher
// @Description Update researcher
// @Tags Researchers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param uid path string true "uid"
// @Param researcher body models.Researcher true "Researcher"
// @Success 200 {object} models.Researcher
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /v1/researcher/{uid} [put]
func (h *PublicationsHandlers) UpdateResearcher() echo.HandlerFunc {

	return func(c echo.Context) error {

		uid := c.Param("uid")

		researcher := new(models.Researcher)
		if err := c.Bind(researcher); err != nil {
			log.Error().Err(err).Msg("Error binding researcher")
			return helpers.BadRequest(err.Error())
		}

		researcher.Uid = uid

		userUID := c.Get("userUID").(string)

		_, err := h.PublicationsService.UpdateResearcher(researcher, userUID)
		if err != nil {
			log.Error().Err(err).Msg("Error updating researcher")
			return echo.ErrInternalServerError
		}

		return c.JSON(200, researcher)
	}
}

// DeleteResearcher Delete researcher by uid godoc
// @Summary Delete researcher by uid
// @Description Delete researcher by uid
// @Tags Researchers
// @Security BearerAuth
// @Produce json
// @Param uid path string true "uid"
// @Success 204 "No Content"
// @Failure 500 "Internal Server Error"
// @Router /v1/researcher/{uid} [delete]
func (h *PublicationsHandlers) DeleteResearcher() echo.HandlerFunc {

	return func(c echo.Context) error {

		uid := c.Param("uid")
		userUID := c.Get("userUID").(string)

		err := h.PublicationsService.DeleteResearcher(uid, userUID)
		if err != nil {
			log.Error().Err(err).Msg("Error deleting researcher")
			return echo.ErrInternalServerError
		}

		return c.NoContent(204)
	}
}
