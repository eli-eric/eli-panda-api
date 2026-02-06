package systemsService

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"panda/apigateway/helpers"
	codebookModels "panda/apigateway/services/codebook-service/models"
	"panda/apigateway/services/systems-service/models"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

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
	CreateNewSystemFromJira() echo.HandlerFunc
	UpdateSystem() echo.HandlerFunc
	DeleteSystemRecursive() echo.HandlerFunc
	GetSystemsWithSearchAndPagination() echo.HandlerFunc
	GetSystemsHierarchy() echo.HandlerFunc
	GetSystemLeavesByParentUID() echo.HandlerFunc
	GetSystemsForControlsSystems() echo.HandlerFunc
	GetNewSystemCodesPreview() echo.HandlerFunc
	SaveNewSystemCodes() echo.HandlerFunc
	GetSystemsForRelationship() echo.HandlerFunc
	GetSystemRelationships() echo.HandlerFunc
	DeleteSystemRelationship() echo.HandlerFunc
	CreateNewSystemRelationship() echo.HandlerFunc
	GetSystemCode() echo.HandlerFunc
	GetPhysicalItemProperties() echo.HandlerFunc
	UpdatePhysicalItemProperties() echo.HandlerFunc
	GetSystemHistory() echo.HandlerFunc
	GetSystemTypeGroups() echo.HandlerFunc
	GetSystemTypeGroupsTree() echo.HandlerFunc
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
	SyncSystemLocationByEUNs() echo.HandlerFunc
	GetAllLocationsFlat() echo.HandlerFunc
	GetAllSystemTypes() echo.HandlerFunc
	GetAllZones() echo.HandlerFunc
	CreateNewSystemCode() echo.HandlerFunc
	RecalculateSpareParts() echo.HandlerFunc
	GetSystemsTreeByUids() echo.HandlerFunc
	MovePhysicalItem() echo.HandlerFunc
	ReplacePhysicalItems() echo.HandlerFunc
	MoveSystems() echo.HandlerFunc
	CopySystem() echo.HandlerFunc
	AssignSpareItem() echo.HandlerFunc
	GetSystemSparePartsDetail() echo.HandlerFunc
}

// NewCommentsHandlers Comments handlers constructor
func NewsystemsHandlers(systemsSvc ISystemsService) ISystemsHandlers {
	return &SystemsHandlers{systemsService: systemsSvc}
}

// Swagger documentation for AssignSpareItem
// @Summary Assign spare item to system
// @Description Assigns a spare item to a system, updating the item's condition and system association
// @Tags Systems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.AssignSpareRequest true "Assign spare item request"
// @Success 200 {object} models.AssignSpareResponse "Returns the updated spare item information"
// @Failure 400 {string} string "Bad request - missing required fields or invalid data"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/system/{systemUid}/assign-spare [post]
func (h *SystemsHandlers) AssignSpareItem() echo.HandlerFunc {

	return func(c echo.Context) error {

		// bind assign spare item request data from request body
		assignSpareRequest := new(models.AssignSpareRequest)
		err := c.Bind(assignSpareRequest)
		if err == nil {

			userUID := c.Get("userUID").(string)

			response, err := h.systemsService.AssignSpareItem(*assignSpareRequest, userUID)

			if err == nil {
				return c.JSON(http.StatusOK, response)
			}

			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError

		} else {
			log.Error().Msg(err.Error())
		}
		return helpers.BadRequest(err.Error())
	}
}

// GetSubSystemsByParentUID godoc
// @Summary Get subsystems by parent UID
// @Description Returns subsystems for the given parent system UID.
// @Tags Systems
// @Produce json
// @Security BearerAuth
// @Param parentUID path string true "Parent system UID"
// @Success 200 {array} models.System
// @Failure 500 "Internal server error"
// @Router /v1/system/{parentUID}/subsystems [get]
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

// GetSystemImageByUid godoc
// @Summary Get system image
// @Description Returns base64-encoded image string for the system.
// @Tags Systems
// @Produce plain
// @Security BearerAuth
// @Param uid path string true "System UID"
// @Success 200 {string} string
// @Failure 500 "Internal server error"
// @Router /v1/system/{uid}/image [get]
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

// GetSystemDetail godoc
// @Summary Get system detail
// @Description Returns system detail by UID.
// @Tags Systems
// @Produce json
// @Security BearerAuth
// @Param uid path string true "System UID"
// @Success 200 {object} models.System
// @Failure 500 "Internal server error"
// @Router /v1/system/{uid} [get]
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

// Swagger documentation for CreateNewSystem
// @Summary Create new system
// @Description Creates a new system with the given details
// @Tags Systems
// @Accept json
// @Produce plain
// @Security BearerAuth
// @Param body body models.System true "System details"
// @Success 201 {string} string "Returns the UID of the created system"
// @Failure 400 {string} string "Bad request - missing required fields"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/system [post]
func (h *SystemsHandlers) CreateNewSystem() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		system := new(models.System)
		err := c.Bind(system)
		if err == nil {

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
		return helpers.BadRequest(err.Error())
	}
}

// Swagger documentation for CreateNewSystemFromJira
// @Summary Create new system from Jira import request
// @Description Creates a new system using data from a Jira import request
// @Tags Systems
// @Accept json
// @Produce plain
// @Security BearerAuth
// @Param body body models.JiraSystemImportRequest true "Jira system import data"
// @Success 201 {string} string "Returns the UID of the created system"
// @Failure 400 {string} string "Bad request - invalid input or system code already exists"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/system/jira-import [post]
func (h *SystemsHandlers) CreateNewSystemFromJira() echo.HandlerFunc {
	return func(c echo.Context) error {
		facilityCode := c.Get("facilityCode").(string)
		userUID := c.Get("userUID").(string)

		// Safely get userRoles with a default empty slice
		var userRoles []string
		if roles, ok := c.Get("userRoles").([]string); ok {
			userRoles = roles
		}

		request := new(models.JiraSystemImportRequest)
		err := c.Bind(request)
		if err != nil {
			log.Error().Msg(err.Error())
			return helpers.BadRequest(err.Error())
		}

		result, err := h.systemsService.CreateNewSystemFromJira(facilityCode, userUID, userRoles, request)

		if err == nil {
			return c.String(http.StatusCreated, result)
		}

		if err == helpers.ERR_DUPLICATE_SYSTEM_CODE {
			return c.String(http.StatusBadRequest, err.Error())
		}

		log.Error().Msg(err.Error())
		return echo.ErrInternalServerError
	}
}

// UpdateSystem godoc
// @Summary Update system
// @Description Updates an existing system.
// @Tags Systems
// @Accept json
// @Security BearerAuth
// @Param uid path string true "System UID"
// @Param body body models.System true "System"
// @Success 204 "No content"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /v1/system/{uid} [put]
func (h *SystemsHandlers) UpdateSystem() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		system := new(models.System)
		err := c.Bind(system)
		if err == nil {

			facilityCode := c.Get("facilityCode").(string)
			userUID := c.Get("userUID").(string)
			system.UID = c.Param("uid")

			err := h.systemsService.UpdateSystem(system, facilityCode, userUID)

			if err == nil {
				return c.NoContent(http.StatusNoContent)
			}

			return echo.ErrInternalServerError

		}
		return helpers.BadRequest(err.Error())
	}
}

// DeleteSystemRecursive godoc
// @Summary Delete system recursively
// @Description Deletes a system and its subsystems. If physical items are attached, returns 409 with details.
// @Tags Systems
// @Security BearerAuth
// @Param uid path string true "System UID"
// @Success 204 "No content"
// @Failure 409 {array} models.SystemPhysicalItemInfo
// @Failure 500 "Internal server error"
// @Router /v1/system/{uid} [delete]
func (h *SystemsHandlers) DeleteSystemRecursive() echo.HandlerFunc {

	return func(c echo.Context) error {

		//get uid path param
		uid := c.Param("uid")
		userUid := c.Get("userUID").(string)

		// first check if there are systems with physical items
		itemsInfo, err := h.systemsService.GetPhysicalItemsBySystemUidRecursive(uid)

		if err != nil {
			log.Err(err).Msg("GetPhysicalItemsBySystemUidRecursive")
			return echo.ErrInternalServerError
		}

		// if there are some connected physical items return confiltc http error and related items and systems
		if len(itemsInfo) > 0 {
			return c.JSON(409, itemsInfo)
		}

		// get catalogue item
		err = h.systemsService.DeleteSystemRecursive(uid, userUid)

		if err == nil {
			return c.NoContent(http.StatusNoContent)
		}

		log.Err(err).Msg("DeleteSystemRecursive")
		return echo.ErrInternalServerError
	}
}

// GetSystemsWithSearchAndPagination godoc
// @Summary Get systems
// @Description Returns a paginated list of systems with optional search/sorting/filtering.
// @Tags Systems
// @Produce json
// @Security BearerAuth
// @Param pagination query string true "Pagination JSON (e.g. {\"page\":1,\"pageSize\":100})"
// @Param sorting query string false "Sorting JSON (array of {id, desc})"
// @Param search query string false "Search text"
// @Param columnFilter query string false "Column filter JSON"
// @Success 200 {object} helpers.PaginationResult[models.System]
// @Failure 500 "Internal server error"
// @Router /v1/systems [get]
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

// GetSystemsHierarchy godoc
// @Summary Get systems hierarchy (parents only)
// @Description Returns a hierarchy tree containing only systems that have subsystems (no leaves).
// @Tags Systems
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.SystemHierarchyNode
// @Failure 500 "Internal server error"
// @Router /v1/systems/hierarchy [get]
func (h *SystemsHandlers) GetSystemsHierarchy() echo.HandlerFunc {
	return func(c echo.Context) error {
		facilityCode := c.Get("facilityCode").(string)

		items, err := h.systemsService.GetSystemsHierarchy(facilityCode)
		if err == nil {
			return c.JSON(http.StatusOK, items)
		}

		log.Error().Msg(err.Error())
		return echo.ErrInternalServerError
	}
}

// GetSystemLeavesByParentUID godoc
// @Summary Get leaf systems for a parent
// @Description Returns a paginated list of leaf systems (systems without subsystems) directly under the given parent system.
// @Tags Systems
// @Produce json
// @Security BearerAuth
// @Param uid path string true "Parent system UID"
// @Param pagination query string true "Pagination JSON (e.g. {\"page\":1,\"pageSize\":100})"
// @Param sorting query string false "Sorting JSON (array of {id, desc})"
// @Param search query string false "Search text"
// @Param columnFilter query string false "Column filter JSON"
// @Success 200 {object} helpers.PaginationResult[models.System]
// @Failure 500 "Internal server error"
// @Router /v1/system/{uid}/leaves [get]
func (h *SystemsHandlers) GetSystemLeavesByParentUID() echo.HandlerFunc {
	return func(c echo.Context) error {
		parentUID := c.Param("uid")
		facilityCode := c.Get("facilityCode").(string)

		pagination := c.QueryParam("pagination")
		sorting := c.QueryParam("sorting")
		search := c.QueryParam("search")

		pagingObject := &helpers.Pagination{Page: 1, PageSize: 100}
		if strings.TrimSpace(pagination) != "" {
			json.Unmarshal([]byte(pagination), &pagingObject)
		} else {
			if pageStr := strings.TrimSpace(c.QueryParam("page")); pageStr != "" {
				if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
					pagingObject.Page = page
				}
			}
			if pageSizeStr := strings.TrimSpace(c.QueryParam("pageSize")); pageSizeStr != "" {
				if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
					pagingObject.PageSize = pageSize
				}
			}
		}

		sortingObject := new([]helpers.Sorting)
		json.Unmarshal([]byte(sorting), &sortingObject)

		filterObject := new([]helpers.ColumnFilter)
		filter := c.QueryParam("columnFilter")
		json.Unmarshal([]byte(filter), &filterObject)

		items, err := h.systemsService.GetSystemLeavesByParentUID(parentUID, facilityCode, search, pagingObject, sortingObject, filterObject)
		if err == nil {
			return c.JSON(http.StatusOK, items)
		}

		log.Error().Msg(err.Error())
		return echo.ErrInternalServerError
	}
}

// Swagger documentation for GetSystemsForControlsSystems
// @Summary Get systems with system codes
// @Description Returns a simplified paginated list of systems where `systemCode` is filled; supports optional filtering by zone, system type, and search text.
// @Tags Systems
// @Produce json
// @Security BearerAuth
// @Param pagination query string true "Pagination JSON (e.g. {\"page\":1,\"pageSize\":100})"
// @Param sorting query string false "Sorting JSON (array of {id, desc})"
// @Param columnFilter query string false "Column filter JSON (may contain ids: zone, systemType, searchText; also accepts search)"
// @Param filter query string false "Alias for columnFilter (for older clients)"
// @Success 200 {object} helpers.PaginationResult[models.SystemCodesResult]
// @Failure 500 "Internal server error"
// @Router /v1/systems/system-codes [get]
func (h *SystemsHandlers) GetSystemsForControlsSystems() echo.HandlerFunc {

	return func(c echo.Context) error {

		pagination := c.QueryParam("pagination")
		sorting := c.QueryParam("sorting")
		facilityCode := c.Get("facilityCode").(string)

		pagingObject := &helpers.Pagination{Page: 1, PageSize: 100}
		if strings.TrimSpace(pagination) != "" {
			json.Unmarshal([]byte(pagination), &pagingObject)
		} else {
			if pageStr := strings.TrimSpace(c.QueryParam("page")); pageStr != "" {
				if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
					pagingObject.Page = page
				}
			}
			if pageSizeStr := strings.TrimSpace(c.QueryParam("pageSize")); pageSizeStr != "" {
				if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
					pagingObject.PageSize = pageSize
				}
			}
		}

		sortingObject := new([]helpers.Sorting)
		json.Unmarshal([]byte(sorting), &sortingObject)

		filterObject := new([]helpers.ColumnFilter)
		filter := c.QueryParam("columnFilter")
		if strings.TrimSpace(filter) == "" {
			filter = c.QueryParam("filter")
		}
		json.Unmarshal([]byte(filter), &filterObject)

		items, err := h.systemsService.GetSystemsForControlsSystems(facilityCode, pagingObject, sortingObject, filterObject)

		if err == nil {
			return c.JSON(http.StatusOK, items)
		}

		log.Error().Msg(err.Error())
		return echo.ErrInternalServerError
	}
}

// Swagger documentation for GetNewSystemCodesPreview
// @Summary Preview new system codes
// @Description Generates a preview of the next N system codes for a given system type and zone without creating any systems.
// @Tags Systems
// @Produce json
// @Security BearerAuth
// @Param systemTypeUid query string true "System type UID"
// @Param zoneUid query string true "Zone UID"
// @Param batch query int false "How many codes to generate (default 1)"
// @Success 200 {array} models.SystemCodesResult
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /v1/systems/system-codes/preview [get]
func (h *SystemsHandlers) GetNewSystemCodesPreview() echo.HandlerFunc {

	return func(c echo.Context) error {

		systemTypeUID := c.QueryParam("systemTypeUid")
		if systemTypeUID == "" {
			systemTypeUID = c.QueryParam("systemType")
		}
		zoneUID := c.QueryParam("zoneUid")
		if zoneUID == "" {
			zoneUID = c.QueryParam("zone")
		}
		if strings.TrimSpace(systemTypeUID) == "" || strings.TrimSpace(zoneUID) == "" {
			return helpers.BadRequest("missing systemTypeUid/zoneUid")
		}
		batchStr := c.QueryParam("batch")
		batch := 1
		if strings.TrimSpace(batchStr) != "" {
			parsed, err := strconv.Atoi(batchStr)
			if err != nil {
				return helpers.BadRequest("invalid batch")
			}
			batch = parsed
		}

		facilityCode := c.Get("facilityCode").(string)

		result, err := h.systemsService.GetNewSystemCodesPreview(systemTypeUID, zoneUID, batch, facilityCode)
		if err == nil {
			return c.JSON(http.StatusOK, result)
		}

		log.Error().Msg(err.Error())
		if err == helpers.ERR_INVALID_INPUT {
			return helpers.BadRequest(err.Error())
		}
		if strings.Contains(err.Error(), "missing default parent system") {
			return helpers.BadRequest(err.Error())
		}
		return echo.ErrInternalServerError
	}
}

// Swagger documentation for SaveNewSystemCodes
// @Summary Create new systems with generated system codes
// @Description Creates a batch of new systems with generated system codes for the given system type and zone. System name is set to the generated system code.
// @Tags Systems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.SystemCodesRequest true "System codes request"
// @Success 201 {array} models.SystemCodesResult
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /v1/systems/system-codes [post]
func (h *SystemsHandlers) SaveNewSystemCodes() echo.HandlerFunc {

	return func(c echo.Context) error {

		request := new(models.SystemCodesRequest)
		err := c.Bind(request)
		if err != nil {
			log.Error().Msg(err.Error())
			return helpers.BadRequest(err.Error())
		}

		facilityCode := c.Get("facilityCode").(string)
		userUID := c.Get("userUID").(string)

		result, err := h.systemsService.SaveNewSystemCodes(request, facilityCode, userUID)
		if err == nil {
			return c.JSON(http.StatusCreated, result)
		}

		log.Error().Msg(err.Error())
		if err == helpers.ERR_INVALID_INPUT {
			return helpers.BadRequest(err.Error())
		}
		return echo.ErrInternalServerError
	}
}

// GetSystemsForRelationship godoc
// @Summary Get systems for relationship
// @Description Returns a paginated list of systems for relationship creation; supports optional filtering.
// @Tags Systems
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search text"
// @Param pagination query string false "Pagination JSON (e.g. {\"page\":1,\"pageSize\":100})"
// @Param sorting query string false "Sorting JSON (array of {id, desc})"
// @Param columnFilter query string false "Column filter JSON"
// @Param systemFromUid query string false "Source system UID"
// @Param relationTypeCode query string false "Relation type code"
// @Success 200 {object} helpers.PaginationResult[models.System]
// @Failure 500 "Internal server error"
// @Router /v1/systems/for-relationship [get]
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

// GetSystemRelationships godoc
// @Summary Get system relationships
// @Description Returns relationships for a given system UID.
// @Tags Systems
// @Produce json
// @Security BearerAuth
// @Param uid path string true "System UID"
// @Success 200 {array} models.SystemRelationship
// @Failure 500 "Internal server error"
// @Router /v1/system/{uid}/relationships [get]
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

// DeleteSystemRelationship godoc
// @Summary Delete system relationship
// @Description Deletes a system relationship by relationship UID.
// @Tags Systems
// @Security BearerAuth
// @Param uid path string true "Relationship UID"
// @Success 204 "No content"
// @Failure 500 "Internal server error"
// @Router /v1/system/relationship/{uid} [delete]
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

// CreateNewSystemRelationship godoc
// @Summary Create new system relationship
// @Description Creates a new relationship between systems.
// @Tags Systems
// @Accept json
// @Produce plain
// @Security BearerAuth
// @Param uid path string true "Unused route parameter (reserved)"
// @Param body body models.SystemRelationshipRequest true "Relationship request"
// @Success 201 {string} string "Created relationship UID"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /v1/system/relationship/{uid} [post]
func (h *SystemsHandlers) CreateNewSystemRelationship() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		systemRelationshipRequest := new(models.SystemRelationshipRequest)
		err := c.Bind(systemRelationshipRequest)
		if err == nil {

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
		return helpers.BadRequest(err.Error())
	}
}

// GetSystemCode godoc
// @Summary Get next system code
// @Description Generates a new unique system code based on system type and zone/location/parent.
// @Tags Systems
// @Produce plain
// @Security BearerAuth
// @Param systemTypeUID query string true "System type UID"
// @Param zoneUID query string true "Zone UID"
// @Param locationUID query string false "Location UID"
// @Param parentUID query string false "Parent system UID"
// @Success 200 {string} string
// @Failure 400 {string} string "Bad request"
// @Failure 500 "Internal server error"
// @Router /v1/system/systemCode [get]
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

// GetPhysicalItemProperties godoc
// @Summary Get physical item properties
// @Description Returns physical item properties by physical item UID.
// @Tags Systems
// @Produce json
// @Security BearerAuth
// @Param uid path string true "Physical item UID"
// @Success 200 {array} models.PhysicalItemDetail
// @Failure 500 "Internal server error"
// @Router /v1/physical-item/{uid}/properties [get]
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

// UpdatePhysicalItemProperties godoc
// @Summary Update physical item properties
// @Description Updates physical item properties by physical item UID.
// @Tags Systems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uid path string true "Physical item UID"
// @Param body body []models.PhysicalItemDetail true "Physical item properties"
// @Success 200 {array} models.PhysicalItemDetail
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /v1/physical-item/{uid}/properties [put]
func (h *SystemsHandlers) UpdatePhysicalItemProperties() echo.HandlerFunc {

	return func(c echo.Context) error {

		uid := c.Param("uid")
		userUid := c.Get("userUID").(string)

		properties := new([]models.PhysicalItemDetail)
		err := c.Bind(properties)
		if err == nil {

			err := h.systemsService.UpdatePhysicalItemProperties(uid, *properties, userUid)

			if err == nil {
				return c.JSON(http.StatusOK, properties)
			}

			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError

		} else {
			log.Error().Msg(err.Error())
		}
		return helpers.BadRequest(err.Error())
	}
}

// GetSystemHistory godoc
// @Summary Get system history
// @Description Returns system history events by system UID.
// @Tags Systems
// @Produce json
// @Security BearerAuth
// @Param uid path string true "System UID"
// @Success 200 {array} models.SystemHistory
// @Failure 500 "Internal server error"
// @Router /v1/system/{uid}/history [get]
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

// GetSystemTypeGroups godoc
// @Summary Get system type groups
// @Description Returns system type groups for the current facility.
// @Tags Systems
// @Produce json
// @Security BearerAuth
// @Success 200 {array} codebookModels.Codebook
// @Failure 500 "Internal server error"
// @Router /v1/system/system-type-groups [get]
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

// @Summary Get system type groups as tree
// @Description Returns system type groups with their system types as a tree structure
// @Tags Systems
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search filter for name or code (case-insensitive)"
// @Success 200 {array} models.SystemTypeGroupTreeItem
// @Failure 500 {string} string "Internal server error"
// @Router /v1/system/system-type-groups/tree [get]
func (h *SystemsHandlers) GetSystemTypeGroupsTree() echo.HandlerFunc {

	return func(c echo.Context) error {

		facilityCode := c.Get("facilityCode").(string)
		search := c.QueryParam("search")

		items, err := h.systemsService.GetSystemTypeGroupsTree(facilityCode, search)

		if err == nil {
			return c.JSON(http.StatusOK, items)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

// GetSystemTypesBySystemTypeGroup godoc
// @Summary Get system types by group
// @Description Returns system types for the given system type group UID.
// @Tags Systems
// @Produce json
// @Security BearerAuth
// @Param uid path string true "System type group UID"
// @Success 200 {array} models.SystemType
// @Failure 500 "Internal server error"
// @Router /v1/system/system-type-group/{uid}/system-types [get]
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

// DeleteSystemTypeGroup godoc
// @Summary Delete system type group
// @Description Deletes a system type group.
// @Tags Systems
// @Security BearerAuth
// @Param uid path string true "System type group UID"
// @Success 204 "No content"
// @Failure 409 {object} helpers.ConflictErrorResponse
// @Failure 500 "Internal server error"
// @Router /v1/system/system-type-group/{uid} [delete]
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

// DeleteSystemType godoc
// @Summary Delete system type
// @Description Deletes a system type.
// @Tags Systems
// @Security BearerAuth
// @Param uid path string true "System type UID"
// @Success 204 "No content"
// @Failure 409 {object} helpers.ConflictErrorResponse
// @Failure 500 "Internal server error"
// @Router /v1/system/system-type/{uid} [delete]
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

// CreateSystemTypeGroup godoc
// @Summary Create system type group
// @Description Creates a new system type group.
// @Tags Systems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body codebookModels.Codebook true "System type group"
// @Success 201 {object} codebookModels.Codebook
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /v1/system/system-type-group [post]
func (h *SystemsHandlers) CreateSystemTypeGroup() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		systemTypeGroup := new(codebookModels.Codebook)
		err := c.Bind(systemTypeGroup)
		if err == nil {

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
		return helpers.BadRequest(err.Error())
	}
}

// UpdateSystemTypeGroup godoc
// @Summary Update system type group
// @Description Updates an existing system type group.
// @Tags Systems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uid path string true "System type group UID"
// @Param body body codebookModels.Codebook true "System type group"
// @Success 200 {object} codebookModels.Codebook
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /v1/system/system-type-group/{uid} [put]
func (h *SystemsHandlers) UpdateSystemTypeGroup() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		systemTypeGroup := new(codebookModels.Codebook)
		err := c.Bind(systemTypeGroup)
		if err == nil {

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

		return helpers.BadRequest(err.Error())
	}
}

// CreateSystemType godoc
// @Summary Create system type
// @Description Creates a new system type under a system type group.
// @Tags Systems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uid path string true "System type group UID"
// @Param body body models.SystemType true "System type"
// @Success 201 {object} models.SystemType
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /v1/system/system-type-group/{uid}/system-type [post]
func (h *SystemsHandlers) CreateSystemType() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		systemType := new(models.SystemType)
		err := c.Bind(systemType)
		if err == nil {

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
		return helpers.BadRequest(err.Error())
	}
}

// UpdateSystemType godoc
// @Summary Update system type
// @Description Updates an existing system type.
// @Tags Systems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param grpUid path string true "System type group UID"
// @Param uid path string true "System type UID"
// @Param body body models.SystemType true "System type"
// @Success 200 {object} models.SystemType
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /v1/system/system-type-group/{grpUid}/system-type/{uid} [put]
func (h *SystemsHandlers) UpdateSystemType() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		systemType := new(models.SystemType)
		err := c.Bind(systemType)
		if err == nil {

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

		return helpers.BadRequest(err.Error())
	}
}

// GetSystemByEun godoc
// @Summary Get system by EUN
// @Description Returns system by EUN.
// @Tags Systems
// @Produce json
// @Security BearerAuth
// @Param eun path string true "EUN"
// @Success 200 {object} models.System
// @Failure 500 "Internal server error"
// @Router /v1/system/by-eun/{eun} [get]
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

// GetSystemAsCsv godoc
// @Summary Export systems to CSV
// @Description Exports systems list to CSV using the same filtering/sorting/search as the systems table.
// @Tags Systems
// @Produce text/csv
// @Security BearerAuth
// @Param sorting query string false "Sorting JSON (array of {id, desc})"
// @Param search query string false "Search text"
// @Param columnFilter query string false "Column filter JSON"
// @Success 200 {file} file
// @Failure 500 "Internal server error"
// @Router /v1/systems/export-to-csv [get]
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
			writer.Comma = ','

			defer writer.Flush()
			// write header based on System struct
			// write header
			writer.Write([]string{
				"Name",
				"SystemType",
				"SystemCode",
				"Zone",
				"Location",
				"Importance",
				"EUN",
				"SerialNumber",
				"ItemUsage",
				"SparePartsCoverage",
				"MinimumSpareParts",
				"CatalogueNumber",
				"CatalogueCategory",
				"Supplier",
				"Owner",
				"Responsible",
				"ParentPath",
				"Description",
				"HasSubsystems",
				"SystemLevel",
				"OrderNumber",
				"OrderUid",
				"UID"})

			for _, item := range items.Data {
				// construct parent path string
				parentPath := ""
				for ip, parent := range item.ParentPath {
					if ip > 0 {
						parentPath += " > "
					}
					parentPath += parent.Name
				}
				systemType := ""
				if item.SystemType != nil {
					systemType = item.SystemType.Name
				}
				description := ""
				if item.Description != nil {
					description = *item.Description
				}
				systemCode := ""
				if item.SystemCode != nil {
					systemCode = *item.SystemCode
				}
				zone := ""
				if item.Zone != nil {
					zone = item.Zone.Code
				}
				location := ""
				if item.Location != nil {
					location = item.Location.Code
				}
				systemLevel := ""
				if item.SystemLevel != nil {
					systemLevel = *item.SystemLevel
				}
				owner := ""
				if item.Owner != nil {
					owner = item.Owner.Name
				}
				responsible := ""
				if item.Responsible != nil {
					responsible = item.Responsible.Name
				}
				importance := ""
				if item.Importance != nil {
					importance = item.Importance.Name
				}
				eun := ""
				if item.PhysicalItem != nil && item.PhysicalItem.EUN != nil {
					eun = *item.PhysicalItem.EUN
				}
				serialNumber := ""
				if item.PhysicalItem != nil && item.PhysicalItem.SerialNumber != nil {
					serialNumber = *item.PhysicalItem.SerialNumber
				}
				itemUsage := ""
				if item.PhysicalItem != nil && item.PhysicalItem.ItemUsage != nil {
					itemUsage = item.PhysicalItem.ItemUsage.Name
				}
				catalogueNumber := ""
				if item.PhysicalItem != nil {
					catalogueNumber = item.PhysicalItem.CatalogueItem.CatalogueNumber
				}
				catalogueCategory := ""
				if item.PhysicalItem != nil {
					catalogueCategory = item.PhysicalItem.CatalogueItem.Category.Name
				}
				orderNumber := ""
				if item.PhysicalItem != nil && item.PhysicalItem.OrderNumber != nil {
					orderNumber = *item.PhysicalItem.OrderNumber
				}
				orderUid := ""
				if item.PhysicalItem != nil && item.PhysicalItem.OrderUid != nil {
					orderUid = *item.PhysicalItem.OrderUid
				}
				supplier := ""
				if item.PhysicalItem != nil && item.PhysicalItem.CatalogueItem.Supplier != nil {
					supplier = item.PhysicalItem.CatalogueItem.Supplier.Name
				}
				sparePartsCoverage := ""
				if item.Statistics != nil && item.Statistics.Sp_coverage != nil {
					sparePartsCoverage = strconv.FormatFloat(float64(*item.Statistics.Sp_coverage), 'f', -1, 32)
				}
				minimumSpareParts := ""
				if item.Statistics != nil && item.Statistics.MinimalSpareParstCount != nil {
					minimumSpareParts = strconv.FormatFloat(float64(*item.Statistics.MinimalSpareParstCount), 'f', -1, 32)
				}
				writer.Write([]string{
					item.Name,
					systemType,
					systemCode,
					zone,
					location,
					importance,
					eun,
					serialNumber,
					itemUsage,
					sparePartsCoverage,
					minimumSpareParts,
					catalogueNumber,
					catalogueCategory,
					supplier,
					owner,
					responsible,
					parentPath,
					description,
					strconv.FormatBool(item.HasSubsystems),
					systemLevel,
					orderNumber,
					orderUid,
					item.UID})
			}

			return nil

		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

// GetEuns godoc
// @Summary Get EUNs
// @Description Returns list of physical item EUNs.
// @Tags Systems
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.EUN
// @Failure 500 "Internal server error"
// @Router /v1/physical-items/euns [get]
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

// Swagger documentation for SyncSystemLocationByEUNs
// @Summary Sync system locations by EUNs
// @Description Sync system locations by EUNs
// @Tags Systems
// @Accept json
// @Security BearerAuth
// @Param body body []models.EunLocation true "EUN with location UID"
// @Success 204 "No content"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /v1/systems/sync-locations-by-eun [post]
func (h *SystemsHandlers) SyncSystemLocationByEUNs() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		eunLocations := new([]models.EunLocation)
		err := c.Bind(eunLocations)
		if err == nil {

			userUID := c.Get("userUID").(string)

			err := h.systemsService.SyncSystemLocationByEUNs(*eunLocations, userUID)

			if err == nil {
				return c.NoContent(http.StatusNoContent)
			}

			return echo.ErrInternalServerError

		} else {
			log.Error().Msg(err.Error())
		}
		return helpers.BadRequest(err.Error())
	}
}

// Swagger documentation for GetAllLocationsFlat
// @Summary Get all locations flat list
// @Description Get all locations flat list
// @Tags Systems
// @Security BearerAuth
// @Success 200 {array} models.Codebook
// @Failure 500 "Internal server error"
// @Router /v1/systems/locations-flat [get]
func (h *SystemsHandlers) GetAllLocationsFlat() echo.HandlerFunc {

	return func(c echo.Context) error {

		facilityCode := c.Get("facilityCode").(string)

		locations, err := h.systemsService.GetAllLocationsFlat(facilityCode)

		if err == nil {
			return c.JSON(http.StatusOK, locations)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

// Swagger documentation for GetAllSystemTypes
// @Summary Get all system types
// @Description Get all system types
// @Tags Systems
// @Security BearerAuth
// @Success 200 {array} models.Codebook
// @Failure 500 "Internal server error"
// @Router /v1/systems/system-types [get]
func (h *SystemsHandlers) GetAllSystemTypes() echo.HandlerFunc {

	return func(c echo.Context) error {

		systemTypes, err := h.systemsService.GetAllSystemTypes()

		if err == nil {
			return c.JSON(http.StatusOK, systemTypes)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

// Swagger documentation for GetAllZones
// @Summary Get all zones
// @Description Get all zones
// @Tags Systems
// @Security BearerAuth
// @Success 200 {array} models.Codebook
// @Failure 500 "Internal server error"
// @Router /v1/systems/zones [get]
func (h *SystemsHandlers) GetAllZones() echo.HandlerFunc {

	return func(c echo.Context) error {

		facilityCode := c.Get("facilityCode").(string)
		zones, err := h.systemsService.GetAllZones(facilityCode)

		if err == nil {
			return c.JSON(http.StatusOK, zones)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

// Swagger documentation for CreateNewSystemCode
// @Summary Create new system with code
// @Description Create new system with new unique system code based on system type and zone
// @Tags Systems
// @Security BearerAuth
// @Param body body models.SystemCodeRequest true "System code request model"
// @Success 201 {object} models.System
// @Failure 400 "Bad request - missing required fields"
// @Failure 500 "Internal server error"
// @Router /v1/system/system-code [post]
func (h *SystemsHandlers) CreateNewSystemCode() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		systemCode := new(models.SystemCodeRequest)
		err := c.Bind(systemCode)
		if err == nil {

			facilityCode := c.Get("facilityCode").(string)
			userUID := c.Get("userUID").(string)

			newSystem, err := h.systemsService.CreateNewSystemCode(systemCode.ParentUID, systemCode.SystemTypeUID, systemCode.ZoneUID, facilityCode, userUID)

			if err == nil {
				return c.JSON(http.StatusCreated, newSystem)
			}

			if strings.Contains(err.Error(), "missing") {
				return c.String(http.StatusBadRequest, err.Error())
			}

			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError

		} else {
			log.Error().Msg(err.Error())
		}
		return helpers.BadRequest(err.Error())
	}
}

// Swagger documentation for RecalculateSpareParts
// @Summary Recalculate spare parts
// @Description Recalculate spare parts for all systems
// @Tags Systems
// @Security BearerAuth
// @Success 204 "No content"
// @Failure 500 "Internal server error"
// @Router /v1/systems/recalculate-spare-parts [post]
func (h *SystemsHandlers) RecalculateSpareParts() echo.HandlerFunc {

	return func(c echo.Context) error {

		err := h.systemsService.RecalculateSpareParts()

		if err == nil {
			return c.NoContent(http.StatusNoContent)
		}

		return echo.ErrInternalServerError
	}
}

// Swagger documentation for GetSystemsTreeByUids
// @Summary Get systems tree by UIDs
// @Description Get systems tree by UIDs
// @Tags Systems
// @Security BearerAuth
// @Param body body []models.SystemTreeUid true "Array of system tree UIDs"
// @Success 200 {array} models.System
// @Failure 500 "Internal server error"
// @Router /v1/systems/reload [post]
func (h *SystemsHandlers) GetSystemsTreeByUids() echo.HandlerFunc {

	return func(c echo.Context) error {

		// get array of SystemTreeUid from the body
		uids := new([]models.SystemTreeUid)

		if err := c.Bind(uids); err != nil {
			log.Error().Msg(err.Error())
			return helpers.BadRequest(err.Error())
		}

		systems, err := h.systemsService.GetSystemsTreeByUids(*uids)

		helpers.ProcessArrayResult(&systems, err)

		if err == nil {
			return c.JSON(http.StatusOK, systems)
		} else {
			log.Error().Msg(err.Error())
			return echo.ErrInternalServerError
		}
	}
}

// Swagger documentation for MovePhysicalItem
// @Summary Move physical item
// @Description Move physical item from one system to another
// @Tags Systems
// @Security BearerAuth
// @Param body body models.PhysicalItemMovement true "Move physical item request model"
// @Success 200 "Return destination system UID"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /v1/physical-item/move [post]
func (h *SystemsHandlers) MovePhysicalItem() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		movePhysicalItemRequest := new(models.PhysicalItemMovement)
		err := c.Bind(movePhysicalItemRequest)
		if err == nil {

			log.Info().Msgf("Move physical item request: %+v", movePhysicalItemRequest)

			userUID := c.Get("userUID").(string)
			facilityCode := c.Get("facilityCode").(string)

			destinationSystemUID, err := h.systemsService.MovePhysicalItem(movePhysicalItemRequest, userUID, facilityCode)

			if err == nil {
				return c.String(http.StatusOK, destinationSystemUID)
			}

			log.Error().Msg(err.Error())

			if strings.Contains(err.Error(), "missing") || strings.Contains(err.Error(), "destination system") {
				return c.String(http.StatusBadRequest, err.Error())
			}

			return echo.ErrInternalServerError

		} else {
			log.Error().Msg(err.Error())
		}
		return helpers.BadRequest(err.Error())
	}
}

// Swagger documentation for ReplacePhysicalItems
// @Summary Replace physical item
// @Description Replace physical items between two systems
// @Tags Systems
// @Security BearerAuth
// @Param body body models.PhysicalItemMovement true "Move physical item request model"
// @Success 200 "Return destination system UID"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /v1/physical-item/replace [post]
func (h *SystemsHandlers) ReplacePhysicalItems() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		movePhysicalItemRequest := new(models.PhysicalItemMovement)
		err := c.Bind(movePhysicalItemRequest)
		if err == nil {

			log.Info().Msgf("Move physical item request: %+v", movePhysicalItemRequest)

			userUID := c.Get("userUID").(string)
			facilityCode := c.Get("facilityCode").(string)

			destinationSystemUID, err := h.systemsService.ReplacePhysicalItems(movePhysicalItemRequest, userUID, facilityCode)

			if err == nil {
				return c.String(http.StatusOK, destinationSystemUID)
			}

			log.Error().Msg(err.Error())

			if strings.Contains(err.Error(), "missing") || strings.Contains(err.Error(), "system") {
				return c.String(http.StatusBadRequest, err.Error())
			}

			return echo.ErrInternalServerError

		} else {
			log.Error().Msg(err.Error())
		}
		return helpers.BadRequest(err.Error())
	}
}

// Swagger documentation for MoveSystems
// @Summary Move systems
// @Description Move systems to another parent system
// @Tags Systems
// @Security BearerAuth
// @Param body body models.SystemsMovement true "Move systems request model"
// @Success 200 "Return destination system UID"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /v1/systems/move [post]
func (h *SystemsHandlers) MoveSystems() echo.HandlerFunc {

	return func(c echo.Context) error {

		// lets bind catalogue category data from request body
		moveSystemsRequest := new(models.SystemsMovement)
		err := c.Bind(moveSystemsRequest)
		if err == nil {

			log.Info().Msgf("Move systems request: %+v", moveSystemsRequest)

			userUID := c.Get("userUID").(string)

			destinationSystemUID, err := h.systemsService.MoveSystems(moveSystemsRequest, userUID)

			if err == nil {
				return c.String(http.StatusOK, destinationSystemUID)
			}

			log.Error().Msg(err.Error())

			if strings.Contains(err.Error(), "missing") || strings.Contains(err.Error(), "system") {
				return c.String(http.StatusBadRequest, err.Error())
			}

			return echo.ErrInternalServerError

		} else {
			log.Error().Msg(err.Error())
		}
		return helpers.BadRequest(err.Error())
	}
}

// Swagger documentation for CopySystem
// @Summary Copy system(s)
// @Description Copies system (or its children) under an existing destination parent system. Copies only Name, SystemLevel, HAS_SUBSYSTEM and HAS_SYSTEM_TYPE.
// @Tags Systems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.SystemCopyRequest true "Copy system request model"
// @Success 201 {array} string "Created root system UID(s)"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Source or destination system not found"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/systems/copy [put]
func (h *SystemsHandlers) CopySystem() echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(models.SystemCopyRequest)
		err := c.Bind(request)
		if err != nil {
			log.Error().Msg(err.Error())
			return helpers.BadRequest(err.Error())
		}

		request.SourceSystemUID = strings.TrimSpace(request.SourceSystemUID)
		request.DestinationSystemUID = strings.TrimSpace(request.DestinationSystemUID)

		if request.SourceSystemUID == "" {
			return helpers.BadRequest("sourceSystemUid is required")
		}
		if request.DestinationSystemUID == "" {
			return helpers.BadRequest("destinationSystemUid is required")
		}

		userUID := c.Get("userUID").(string)
		facilityCode := c.Get("facilityCode").(string)

		createdRootUids, err := h.systemsService.CopySystem(request, facilityCode, userUID)
		if err == nil {
			return c.JSON(http.StatusCreated, createdRootUids)
		}

		log.Error().Msg(err.Error())
		if err == errSourceSystemNotFound || err == errDestinationSystemNotFound {
			return c.String(http.StatusNotFound, err.Error())
		}
		if strings.Contains(err.Error(), "missing") {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return echo.ErrInternalServerError
	}
}

// Swagger documentation for GetSystemSparePartsDetail
// @Summary Get system spare parts detail
// @Description Get comprehensive system and physical item information with all spare relations by system ID
// @Tags Systems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uid path string true "System UID"
// @Success 200 {object} models.SystemSparePartsDetail "Returns comprehensive system spare parts information"
// @Failure 400 {string} string "Bad request - invalid system ID"
// @Failure 404 {string} string "System not found"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/system/{uid}/spare-parts-detail [get]
func (h *SystemsHandlers) GetSystemSparePartsDetail() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get system UID from path parameter
		systemUid := c.Param("uid")
		facilityCode := c.Get("facilityCode").(string)

		// Call the service method
		result, err := h.systemsService.GetSystemSparePartsDetail(systemUid, facilityCode)

		if err == nil {
			return c.JSON(http.StatusOK, result)
		} else {
			log.Error().Err(err).Str("systemUid", systemUid).Msg("Error getting system spare parts detail")
			return echo.ErrInternalServerError
		}
	}
}
