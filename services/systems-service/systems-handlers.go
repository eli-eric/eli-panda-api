package systemsService

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"panda/apigateway/helpers"
	codebookModels "panda/apigateway/services/codebook-service/models"
	"panda/apigateway/services/systems-service/models"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"

	"github.com/labstack/echo/v4"
)

type SystemsHandlers struct {
	systemsService ISystemsService
}

type ISystemsHandlers interface {
	GetSubSystemsByParentUID() echo.HandlerFunc
	GetSystemImageByUid() echo.HandlerFunc
	GetSystemDetail() echo.HandlerFunc
	CreateNewSystem() echo.HandlerFunc
	UpdateSystem() echo.HandlerFunc
	DeleteSystemRecursive() echo.HandlerFunc
	GetSystemsWithSearchAndPagination() echo.HandlerFunc
	GetSystemsForRelationship() echo.HandlerFunc
	GetSystemRelationships() echo.HandlerFunc
	DeleteSystemRelationship() echo.HandlerFunc
	CreateNewSystemRelationship() echo.HandlerFunc
	GetSystemCode() echo.HandlerFunc
	GetPhysicalItemProperties() echo.HandlerFunc
	UpdatePhysicalItemProperties() echo.HandlerFunc
	GetSystemHistory() echo.HandlerFunc
	GetSystemTypeGroups() echo.HandlerFunc
	GetSystemTypesBySystemTypeGroup() echo.HandlerFunc
	DeleteSystemTypeGroup() echo.HandlerFunc
	DeleteSystemType() echo.HandlerFunc
	CreateSystemTypeGroup() echo.HandlerFunc
	UpdateSystemTypeGroup() echo.HandlerFunc
	CreateSystemType() echo.HandlerFunc
	UpdateSystemType() echo.HandlerFunc
	GetSystemByEun() echo.HandlerFunc
	GetSystemAsCsv() echo.HandlerFunc
	GetEuns() echo.HandlerFunc
	ImportSubsystems() echo.HandlerFunc
}

// NewCommentsHandlers Comments handlers constructor
func NewsystemsHandlers(systemsSvc ISystemsService) ISystemsHandlers {
	return &SystemsHandlers{systemsService: systemsSvc}
}

func (h *SystemsHandlers) GetSubSystemsByParentUID() echo.HandlerFunc {

	return func(c echo.Context) error {

		parentUID := c.Param("parentUID")
		facilityCode := c.Get("facilityCode").(string)
		subSystems, err := h.systemsService.GetSubSystemsByParentUID(parentUID, facilityCode)

		if err == nil {
			return c.JSON(http.StatusOK, subSystems)
		} else {
			log.Error().Msg(err.Error())
		}

		return echo.ErrInternalServerError
	}
}

func (h *SystemsHandlers) GetSystemImageByUid() echo.HandlerFunc {
	return func(c echo.Context) error {

		//get uid path param
		uid := c.Param("uid")

		imageString, err := h.systemsService.GetSystemImageByUid(uid)

		if err == nil {
			return c.String(200, imageString)

		} else {
			log.Error().Msg(err.Error())
		}

		return echo.ErrInternalServerError
	}
}

func (h *SystemsHandlers) GetSystemDetail() echo.HandlerFunc {

	return func(c echo.Context) error {

		uid := c.Param("uid")
		facilityCode := c.Get("facilityCode").(string)
		systemDetail, err := h.systemsService.GetSystemDetail(uid, facilityCode)

		if err == nil {
			return c.JSON(http.StatusOK, systemDetail)
		} else {
			log.Error().Msg(err.Error())
		}

		return echo.ErrInternalServerError

	}
}

func (h *SystemsHandlers) CreateNewSystem() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		system := new(models.System)

		if err := c.Bind(system); err == nil {

			facilityCode := c.Get("facilityCode").(string)
			userUID := c.Get("userUID").(string)

			uid, err := h.systemsService.CreateNewSystem(system, facilityCode, userUID)

			if err == nil {
				return c.String(http.StatusCreated, uid)
			}

			return echo.ErrInternalServerError

		} else {
			log.Error().Msg(err.Error())
		}
		return echo.ErrBadRequest
	}
}

func (h *SystemsHandlers) UpdateSystem() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		system := new(models.System)

		if err := c.Bind(system); err == nil {

			facilityCode := c.Get("facilityCode").(string)
			userUID := c.Get("userUID").(string)
			system.UID = c.Param("uid")

			err := h.systemsService.UpdateSystem(system, facilityCode, userUID)

			if err == nil {
				return c.NoContent(http.StatusNoContent)
			}

			return echo.ErrInternalServerError

		}
		return echo.ErrBadRequest
	}
}

func (h *SystemsHandlers) DeleteSystemRecursive() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get uid path param
		uid := c.Param("uid")

		// get catalogue item
		err := h.systemsService.DeleteSystemRecursive(uid)

		if err == nil {
			return c.NoContent(http.StatusNoContent)
		}

		return echo.ErrInternalServerError
	}
}

func (h *SystemsHandlers) GetSystemsWithSearchAndPagination() echo.HandlerFunc {

	return func(c echo.Context) error {

		pagination := c.QueryParam("pagination")
		sorting := c.QueryParam("sorting")
		search := c.QueryParam("search")
		facilityCode := c.Get("facilityCode").(string)

		pagingObject := new(helpers.Pagination)
		json.Unmarshal([]byte(pagination), &pagingObject)

		sortingObject := new([]helpers.Sorting)
		json.Unmarshal([]byte(sorting), &sortingObject)

		filterObject := new([]helpers.ColumnFilter)
		filter := c.QueryParam("columnFilter")
		json.Unmarshal([]byte(filter), &filterObject)

		items, err := h.systemsService.GetSystemsWithSearchAndPagination(search, facilityCode, pagingObject, sortingObject, filterObject)

		if err == nil {
			return c.JSON(http.StatusOK, items)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

func (h *SystemsHandlers) GetSystemsForRelationship() echo.HandlerFunc {

	return func(c echo.Context) error {

		facilityCode := c.Get("facilityCode").(string)
		search := c.QueryParam("search")

		pagingObject := new(helpers.Pagination)
		pagination := c.QueryParam("pagination")
		json.Unmarshal([]byte(pagination), &pagingObject)

		sortingObject := new([]helpers.Sorting)
		sorting := c.QueryParam("sorting")
		json.Unmarshal([]byte(sorting), &sortingObject)

		systemFromUid := c.QueryParam("systemFromUid")
		relationTypeCode := c.QueryParam("relationTypeCode")

		filterObject := new([]helpers.ColumnFilter)
		filter := c.QueryParam("columnFilter")
		json.Unmarshal([]byte(filter), &filterObject)

		items, err := h.systemsService.GetSystemsForRelationship(search, facilityCode, pagingObject, sortingObject, filterObject, systemFromUid, relationTypeCode)

		if err == nil {
			return c.JSON(http.StatusOK, items)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

func (h *SystemsHandlers) GetSystemRelationships() echo.HandlerFunc {

	return func(c echo.Context) error {

		systemUid := c.Param("uid")

		items, err := h.systemsService.GetSystemRelationships(systemUid)

		if err == nil {
			return c.JSON(http.StatusOK, items)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

func (h *SystemsHandlers) DeleteSystemRelationship() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get userUID
		userUID := c.Get("userUID").(string)
		//get uid path param
		uid := c.Param("uid")
		//convert uid to int64
		uidInt64, err := strconv.ParseInt(uid, 10, 64)

		if err != nil {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}

		err = h.systemsService.DeleteSystemRelationship(uidInt64, userUID)

		if err == nil {
			return c.NoContent(http.StatusNoContent)
		}

		return echo.ErrInternalServerError
	}
}

func (h *SystemsHandlers) CreateNewSystemRelationship() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		systemRelationshipRequest := new(models.SystemRelationshipRequest)

		if err := c.Bind(systemRelationshipRequest); err == nil {

			userUID := c.Get("userUID").(string)
			facilityCode := c.Get("facilityCode").(string)

			newId, err := h.systemsService.CreateNewSystemRelationship(systemRelationshipRequest, facilityCode, userUID)

			if err == nil {
				return c.String(http.StatusCreated, strconv.FormatInt(newId, 10))
			}

			return echo.ErrInternalServerError

		} else {
			log.Error().Msg(err.Error())
		}
		return echo.ErrBadRequest
	}
}

func (h *SystemsHandlers) GetSystemCode() echo.HandlerFunc {

	return func(c echo.Context) error {

		facilityCode := c.Get("facilityCode").(string)
		zoneUID := c.QueryParam("zoneUID")
		locationUID := c.QueryParam("locationUID")
		parentUID := c.QueryParam("parentUID")
		systemTypeUID := c.QueryParam("systemTypeUID")

		code, err := h.systemsService.GetSystemCode(systemTypeUID, zoneUID, locationUID, parentUID, facilityCode)

		if err == nil {
			return c.String(http.StatusOK, code)
		} else if strings.Contains(err.Error(), "missing") {
			return c.String(http.StatusBadRequest, err.Error())
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

func (h *SystemsHandlers) GetPhysicalItemProperties() echo.HandlerFunc {

	return func(c echo.Context) error {

		uid := c.Param("uid")

		properties, err := h.systemsService.GetPhysicalItemProperties(uid)

		if err == nil {
			return c.JSON(http.StatusOK, properties)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

func (h *SystemsHandlers) UpdatePhysicalItemProperties() echo.HandlerFunc {

	return func(c echo.Context) error {

		uid := c.Param("uid")
		userUid := c.Get("userUID").(string)

		properties := new([]models.PhysicalItemDetail)

		if err := c.Bind(properties); err == nil {

			err := h.systemsService.UpdatePhysicalItemProperties(uid, *properties, userUid)

			if err == nil {
				return c.JSON(http.StatusOK, properties)
			}

			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError

		} else {
			log.Error().Msg(err.Error())
		}
		return echo.ErrBadRequest
	}
}

func (h *SystemsHandlers) GetSystemHistory() echo.HandlerFunc {

	return func(c echo.Context) error {

		uid := c.Param("uid")

		history, err := h.systemsService.GetSystemHistory(uid)

		if err == nil {
			return c.JSON(http.StatusOK, history)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

func (h *SystemsHandlers) GetSystemTypeGroups() echo.HandlerFunc {

	return func(c echo.Context) error {

		facilityCode := c.Get("facilityCode").(string)

		items, err := h.systemsService.GetSystemTypeGroups(facilityCode)

		if err == nil {
			return c.JSON(http.StatusOK, items)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

func (h *SystemsHandlers) GetSystemTypesBySystemTypeGroup() echo.HandlerFunc {

	return func(c echo.Context) error {

		facilityCode := c.Get("facilityCode").(string)
		systemTypeGroupUID := c.Param("uid")

		items, err := h.systemsService.GetSystemTypesBySystemTypeGroup(systemTypeGroupUID, facilityCode)

		if err == nil {
			return c.JSON(http.StatusOK, items)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

func (h *SystemsHandlers) DeleteSystemTypeGroup() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get uid path param
		uid := c.Param("uid")

		err, realtedNodes := h.systemsService.DeleteSystemTypeGroup(uid)

		if len(realtedNodes) > 0 {
			relatedItemsRespponse := helpers.ConflictErrorResponse{
				ErrorMessage: "Cannot delete this item because it is related to other items",
				RelatedNodes: realtedNodes,
			}
			return c.JSON(http.StatusConflict, relatedItemsRespponse)
		}

		if err == nil {
			return c.NoContent(http.StatusNoContent)
		}

		return echo.ErrInternalServerError
	}
}

func (h *SystemsHandlers) DeleteSystemType() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get uid path param
		uid := c.Param("uid")

		err, realtedNodes := h.systemsService.DeleteSystemType(uid)

		if len(realtedNodes) > 0 {
			relatedItemsRespponse := helpers.ConflictErrorResponse{
				ErrorMessage: "Cannot delete this item because it is related to other items",
				RelatedNodes: realtedNodes,
			}
			return c.JSON(http.StatusConflict, relatedItemsRespponse)
		}

		if err == nil {
			return c.NoContent(http.StatusNoContent)
		}

		return echo.ErrInternalServerError
	}
}

func (h *SystemsHandlers) CreateSystemTypeGroup() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		systemTypeGroup := new(codebookModels.Codebook)

		if err := c.Bind(systemTypeGroup); err == nil {

			facilityCode := c.Get("facilityCode").(string)
			userUID := c.Get("userUID").(string)

			err := h.systemsService.CreateSystemTypeGroup(systemTypeGroup, facilityCode, userUID)

			if err == nil {
				return c.JSON(http.StatusCreated, systemTypeGroup)
			}

			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError

		} else {
			log.Error().Msg(err.Error())
		}
		return echo.ErrBadRequest
	}
}

func (h *SystemsHandlers) UpdateSystemTypeGroup() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		systemTypeGroup := new(codebookModels.Codebook)

		if err := c.Bind(systemTypeGroup); err == nil {

			userUID := c.Get("userUID").(string)
			systemTypeGroup.UID = c.Param("uid")

			err := h.systemsService.UpdateSystemTypeGroup(systemTypeGroup, userUID)

			if err == nil {
				return c.JSON(http.StatusOK, systemTypeGroup)
			}

			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError

		} else {
			log.Error().Msg(err.Error())
		}

		return echo.ErrBadRequest
	}
}

func (h *SystemsHandlers) CreateSystemType() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		systemType := new(models.SystemType)

		if err := c.Bind(systemType); err == nil {

			facilityCode := c.Get("facilityCode").(string)
			userUID := c.Get("userUID").(string)
			systemTypeGroupUid := c.Param("uid")

			err := h.systemsService.CreateSystemType(systemType, facilityCode, userUID, systemTypeGroupUid)

			if err == nil {
				return c.JSON(http.StatusCreated, systemType)
			}

			return echo.ErrInternalServerError

		} else {
			log.Error().Msg(err.Error())
		}
		return echo.ErrBadRequest
	}
}

func (h *SystemsHandlers) UpdateSystemType() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		systemType := new(models.SystemType)

		if err := c.Bind(systemType); err == nil {

			facilityCode := c.Get("facilityCode").(string)
			userUID := c.Get("userUID").(string)
			systemType.UID = c.Param("uid")

			err := h.systemsService.UpdateSystemType(systemType, facilityCode, userUID)

			if err == nil {
				return c.JSON(http.StatusOK, systemType)
			}

			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError

		} else {
			log.Error().Msg(err.Error())
		}

		return echo.ErrBadRequest
	}
}

func (h *SystemsHandlers) GetSystemByEun() echo.HandlerFunc {

	return func(c echo.Context) error {

		eun := c.Param("eun")

		system, err := h.systemsService.GetSystemByEun(eun)

		if err == nil {
			return c.JSON(http.StatusOK, system)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

func (h *SystemsHandlers) GetSystemAsCsv() echo.HandlerFunc {

	return func(c echo.Context) error {

		sorting := c.QueryParam("sorting")
		search := c.QueryParam("search")
		facilityCode := c.Get("facilityCode").(string)

		pagingObject := new(helpers.Pagination)
		pagingObject.Page = 1
		pagingObject.PageSize = 100000

		sortingObject := new([]helpers.Sorting)
		json.Unmarshal([]byte(sorting), &sortingObject)

		filterObject := new([]helpers.ColumnFilter)
		filter := c.QueryParam("columnFilter")
		json.Unmarshal([]byte(filter), &filterObject)

		items, err := h.systemsService.GetSystemsWithSearchAndPagination(search, facilityCode, pagingObject, sortingObject, filterObject)

		if err == nil {

			c.Response().Header().Set(echo.HeaderContentType, "text/csv")
			c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename=systems.csv")

			writer := csv.NewWriter(c.Response())
			writer.UseCRLF = true
			writer.Comma = ';'

			defer writer.Flush()

			// write header
			writer.Write([]string{"UID", "Name"})

			for _, item := range items.Data {
				writer.Write([]string{item.UID, item.Name})
			}

			return nil

		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

func (h *SystemsHandlers) GetEuns() echo.HandlerFunc {

	return func(c echo.Context) error {

		facilityCode := c.Get("facilityCode").(string)

		euns, err := h.systemsService.GetEuns(facilityCode)

		if err == nil {
			return c.JSON(http.StatusOK, euns)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

// ImportSubsystems import subsystems from csv as text in body
func (h *SystemsHandlers) ImportSubsystems() echo.HandlerFunc {

	return func(c echo.Context) error {

		// get user UID
		//userUID := c.Get("userUID").(string)

		// get parent UID
		//parentUID := c.Param("uid")

		// Get the uploaded file
		file, err := c.FormFile("file")
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to get uploaded file"})
		}

		// Open the uploaded file
		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open uploaded file"})
		}
		defer src.Close()

		// Create a temporary file to save the uploaded content
		dst, err := os.CreateTemp("", "uploaded-*.xlsx")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create temp file"})
		}
		defer dst.Close()

		// Copy the uploaded file content to the temporary file
		if _, err = io.Copy(dst, src); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save uploaded file"})
		}

		// Open the saved file with excelize
		excelFile, err := excelize.OpenFile(dst.Name())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open excel file"})
		}

		// Get the first sheet name
		sheetName := excelFile.GetSheetName(0)

		// Get all the rows in the first sheet
		rows, err := excelFile.GetRows(sheetName)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get rows from excel file"})
		}

		// Assume the first row is the header and map to struct fields
		if len(rows) < 2 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Excel file must have at least two rows (header and data)"})
		}

		header := rows[0]
		var systems []models.SystemImportData
		for _, row := range rows[1:] {
			if len(row) < len(header) {
				continue // skip incomplete rows
			}
			system := models.SystemImportData{}
			for i, cell := range row {
				switch header[i] {
				case "SystemCode":
					system.SystemCode = cell
				case "SystemName":
					system.SystemName = cell
				case "SystemType":
					system.SystemType = cell
				case "ZoneCode":
					system.ZoneCode = cell
				case "LocationCode":
					system.LocationCode = cell
				case "SystemLevel":
					system.SystemLevel = cell
				}
			}
			systems = append(systems, system)
		}

		// Return the data as JSON
		return c.JSON(http.StatusOK, systems)
	}
}
