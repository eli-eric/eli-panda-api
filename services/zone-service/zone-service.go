package zoneservice

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"panda/apigateway/helpers"
	"panda/apigateway/services/zone-service/models"
	"strings"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

var (
	ErrNotFound = errors.New("NOT_FOUND:zone not found")
)

type ZoneService struct {
	neo4jDriver *neo4j.Driver
}

type IZoneService interface {
	GetAllZones(facilityCode, search string, page, pageSize int, sorting *[]helpers.Sorting) ([]models.Zone, int64, error)
	GetZoneByUID(uid, facilityCode string) (models.Zone, error)
	CreateZone(facilityCode, userUID string, req *models.ZoneCreateRequest) (models.Zone, error)
	UpdateZone(uid, facilityCode, userUID string, req *models.ZoneUpdateRequest) (models.Zone, error)
	DeleteZone(uid, facilityCode, userUID string) error
	ImportZones(facilityCode, userUID string, file io.Reader) (models.ZoneImportResult, error)
}

func NewZoneService(driver *neo4j.Driver) IZoneService {
	return &ZoneService{neo4jDriver: driver}
}

func (svc *ZoneService) GetAllZones(facilityCode, search string, page, pageSize int, sorting *[]helpers.Sorting) (result []models.Zone, totalCount int64, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	limit := pageSize
	if limit <= 0 {
		limit = 50
	}
	skip := 0
	if page > 1 {
		skip = (page - 1) * limit
	}

	query := GetAllZonesQuery(facilityCode, search, skip, limit, sorting)
	result, err = helpers.GetNeo4jArrayOfNodes[models.Zone](session, query)
	helpers.ProcessArrayResult(&result, err)

	if err != nil {
		return result, 0, err
	}

	countQuery := GetAllZonesCountQuery(facilityCode, search)
	totalCount, err = helpers.GetNeo4jSingleRecordSingleValue[int64](session, countQuery)

	return result, totalCount, err
}

func (svc *ZoneService) GetZoneByUID(uid, facilityCode string) (result models.Zone, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	query := GetZoneByUIDQuery(uid, facilityCode)
	result, err = helpers.GetNeo4jSingleRecordAndMapToStruct[models.Zone](session, query)
	return result, err
}

func (svc *ZoneService) CreateZone(facilityCode, userUID string, req *models.ZoneCreateRequest) (result models.Zone, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	// validate code uniqueness per facility
	codeCount, err := helpers.GetNeo4jSingleRecordSingleValue[int64](session, CheckZoneCodeExistsQuery(req.Code, facilityCode, ""))
	if err != nil {
		return result, err
	}
	if codeCount > 0 {
		return result, fmt.Errorf("zone with code '%s' already exists in this facility", req.Code)
	}

	uid := uuid.New().String()

	if req.ParentUID != nil && *req.ParentUID != "" {
		if err := svc.validateParentIsRoot(*req.ParentUID, facilityCode); err != nil {
			return result, err
		}

		query := CreateSubZoneQuery(uid, req.Name, req.Code, facilityCode, *req.ParentUID, userUID)
		result, err = helpers.WriteNeo4jReturnSingleRecordAndMapToStruct[models.Zone](session, query)
	} else {
		query := CreateRootZoneQuery(uid, req.Name, req.Code, facilityCode, userUID)
		result, err = helpers.WriteNeo4jReturnSingleRecordAndMapToStruct[models.Zone](session, query)
	}

	return result, err
}

func (svc *ZoneService) UpdateZone(uid, facilityCode, userUID string, req *models.ZoneUpdateRequest) (result models.Zone, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	// check zone exists in this facility
	_, err = svc.GetZoneByUID(uid, facilityCode)
	if err != nil {
		if strings.Contains(err.Error(), "no more records") {
			return result, ErrNotFound
		}
		return result, err
	}

	// validate code uniqueness (exclude current zone)
	codeCount, err := helpers.GetNeo4jSingleRecordSingleValue[int64](session, CheckZoneCodeExistsQuery(req.Code, facilityCode, uid))
	if err != nil {
		return result, err
	}
	if codeCount > 0 {
		return result, fmt.Errorf("zone with code '%s' already exists in this facility", req.Code)
	}

	// validate parent before any writes
	if req.ParentUID != nil && *req.ParentUID != "" {
		if *req.ParentUID == uid {
			return result, fmt.Errorf("zone cannot be its own parent")
		}
		if err := svc.validateParentIsRoot(*req.ParentUID, facilityCode); err != nil {
			return result, err
		}
	}

	// update properties
	updateQuery := UpdateZoneQuery(uid, req.Name, req.Code, facilityCode, userUID)
	err = helpers.WriteNeo4jAndReturnNothing(session, updateQuery)
	if err != nil {
		return result, err
	}

	// remove existing parent rel
	removeQuery := RemoveParentRelQuery(uid, facilityCode)
	err = helpers.WriteNeo4jAndReturnNothing(session, removeQuery)
	if err != nil {
		return result, err
	}

	// set new parent if provided
	if req.ParentUID != nil && *req.ParentUID != "" {
		setQuery := SetParentRelQuery(uid, *req.ParentUID, facilityCode)
		err = helpers.WriteNeo4jAndReturnNothing(session, setQuery)
		if err != nil {
			return result, err
		}
	}

	return svc.GetZoneByUID(uid, facilityCode)
}

func (svc *ZoneService) DeleteZone(uid, facilityCode, userUID string) error {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	// check zone exists in this facility
	_, err := svc.GetZoneByUID(uid, facilityCode)
	if err != nil {
		if strings.Contains(err.Error(), "no more records") {
			return ErrNotFound
		}
		return err
	}

	// check subzones
	subCount, err := helpers.GetNeo4jSingleRecordSingleValue[int64](session, CheckZoneHasSubzonesQuery(uid, facilityCode))
	if err != nil {
		return err
	}
	if subCount > 0 {
		return errors.New("CONFLICT:zone has subzones")
	}

	// check system references
	sysCount, err := helpers.GetNeo4jSingleRecordSingleValue[int64](session, CheckZoneHasSystemRefsQuery(uid, facilityCode))
	if err != nil {
		return err
	}
	if sysCount > 0 {
		return errors.New("CONFLICT:zone is referenced by systems")
	}

	// soft delete + history log
	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session,
		helpers.SoftDeleteNodeQuery(uid),
		helpers.HistoryLogQuery(uid, "DELETE", userUID))

	return err
}

func (svc *ZoneService) ImportZones(facilityCode, userUID string, file io.Reader) (result models.ZoneImportResult, err error) {
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// read header
	header, err := reader.Read()
	if err != nil {
		return result, fmt.Errorf("failed to read CSV header: %w", err)
	}

	colIndex := mapCSVColumns(header)
	if colIndex["name"] == -1 || colIndex["code"] == -1 {
		return result, fmt.Errorf("CSV must have 'name' and 'code' columns")
	}

	// single session for entire import
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	rowNum := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: %s", rowNum+1, err.Error()))
			rowNum++
			continue
		}
		rowNum++

		// bounds check
		if colIndex["name"] >= len(record) || colIndex["code"] >= len(record) {
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: too few columns", rowNum))
			continue
		}

		name := strings.TrimSpace(record[colIndex["name"]])
		code := strings.TrimSpace(record[colIndex["code"]])

		if name == "" || code == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: name and code are required", rowNum))
			continue
		}

		// check if zone exists
		existing, existErr := helpers.GetNeo4jSingleRecordSingleValue[int64](session,
			CheckZoneCodeExistsQuery(code, facilityCode, ""))

		if existErr != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: error checking existence: %s", rowNum, existErr.Error()))
			continue
		}

		if existing > 0 {
			result.Skipped++
			continue
		}

		// handle parent
		var parentUID *string
		if colIndex["parentCode"] != -1 && colIndex["parentCode"] < len(record) {
			parentCode := strings.TrimSpace(record[colIndex["parentCode"]])
			if parentCode != "" {
				parentZone, pErr := helpers.GetNeo4jSingleRecordAndMapToStruct[models.Zone](session,
					GetZoneByCodeAndFacilityQuery(parentCode, facilityCode))

				if pErr != nil {
					result.Errors = append(result.Errors, fmt.Sprintf("row %d: parent zone '%s' not found", rowNum, parentCode))
					continue
				}

				if parentZone.ParentZone != nil {
					result.Errors = append(result.Errors, fmt.Sprintf("row %d: parent zone '%s' is already a subzone (max 2 levels)", rowNum, parentCode))
					continue
				}

				parentUID = &parentZone.UID
			}
		}

		// create zone
		req := &models.ZoneCreateRequest{
			Name:      name,
			Code:      code,
			ParentUID: parentUID,
		}
		_, createErr := svc.CreateZone(facilityCode, userUID, req)
		if createErr != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: %s", rowNum, createErr.Error()))
			continue
		}

		result.Created++
	}

	return result, nil
}

func (svc *ZoneService) validateParentIsRoot(parentUID, facilityCode string) error {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	type parentCheck struct {
		UID       string `json:"uid"`
		HasParent bool   `json:"hasParent"`
	}

	check, err := helpers.GetNeo4jSingleRecordAndMapToStruct[parentCheck](session, CheckParentIsRootQuery(parentUID, facilityCode))
	if err != nil {
		if strings.Contains(err.Error(), "no more records") {
			return fmt.Errorf("parent zone not found")
		}
		return fmt.Errorf("error validating parent zone: %w", err)
	}
	if check.HasParent {
		return fmt.Errorf("parent zone is already a subzone (max 2 levels of nesting)")
	}

	return nil
}

func mapCSVColumns(header []string) map[string]int {
	colIndex := map[string]int{
		"name":       -1,
		"code":       -1,
		"parentCode": -1,
	}
	for i, col := range header {
		normalized := strings.TrimSpace(strings.ToLower(col))
		switch normalized {
		case "name":
			colIndex["name"] = i
		case "code":
			colIndex["code"] = i
		case "parentcode", "parent_code":
			colIndex["parentCode"] = i
		}
	}
	return colIndex
}
