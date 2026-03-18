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

type ZoneService struct {
	neo4jDriver *neo4j.Driver
}

type IZoneService interface {
	GetAllZones(facilityCode string) ([]models.Zone, error)
	GetZoneByUID(uid, facilityCode string) (models.Zone, error)
	CreateZone(facilityCode, userUID string, req *models.ZoneCreateRequest) (models.Zone, error)
	UpdateZone(uid, facilityCode, userUID string, req *models.ZoneUpdateRequest) (models.Zone, error)
	DeleteZone(uid, facilityCode, userUID string) error
	ImportZones(facilityCode, userUID string, file io.Reader) (models.ZoneImportResult, error)
}

func NewZoneService(driver *neo4j.Driver) IZoneService {
	return &ZoneService{neo4jDriver: driver}
}

func (svc *ZoneService) GetAllZones(facilityCode string) (result []models.Zone, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	query := GetAllZonesQuery(facilityCode)
	result, err = helpers.GetNeo4jArrayOfNodes[models.Zone](session, query)
	helpers.ProcessArrayResult(&result, err)
	return result, err
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
		// validate parent is root zone (no incoming HAS_SUBZONE)
		if err := svc.validateParentIsRoot(*req.ParentUID); err != nil {
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

	// validate code uniqueness (exclude current zone)
	codeCount, err := helpers.GetNeo4jSingleRecordSingleValue[int64](session, CheckZoneCodeExistsQuery(req.Code, facilityCode, uid))
	if err != nil {
		return result, err
	}
	if codeCount > 0 {
		return result, fmt.Errorf("zone with code '%s' already exists in this facility", req.Code)
	}

	// update properties
	updateQuery := UpdateZoneQuery(uid, req.Name, req.Code, userUID)
	err = helpers.WriteNeo4jAndReturnNothing(session, updateQuery)
	if err != nil {
		return result, err
	}

	// handle parent reassignment
	// first remove existing parent rel
	removeQuery := RemoveParentRelQuery(uid)
	_ = helpers.WriteNeo4jAndReturnNothing(session, removeQuery)

	// set new parent if provided
	if req.ParentUID != nil && *req.ParentUID != "" {
		if err := svc.validateParentIsRoot(*req.ParentUID); err != nil {
			return result, err
		}
		setQuery := SetParentRelQuery(uid, *req.ParentUID)
		err = helpers.WriteNeo4jAndReturnNothing(session, setQuery)
		if err != nil {
			return result, err
		}
	}

	// return updated zone
	return svc.GetZoneByUID(uid, facilityCode)
}

func (svc *ZoneService) DeleteZone(uid, facilityCode, userUID string) error {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	// check subzones
	subCount, err := helpers.GetNeo4jSingleRecordSingleValue[int64](session, CheckZoneHasSubzonesQuery(uid))
	if err != nil {
		return err
	}
	if subCount > 0 {
		return errors.New("CONFLICT:zone has subzones")
	}

	// check system references
	sysCount, err := helpers.GetNeo4jSingleRecordSingleValue[int64](session, CheckZoneHasSystemRefsQuery(uid))
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

		name := strings.TrimSpace(record[colIndex["name"]])
		code := strings.TrimSpace(record[colIndex["code"]])

		if name == "" || code == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: name and code are required", rowNum))
			continue
		}

		// check if zone exists
		session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
		existing, existErr := helpers.GetNeo4jSingleRecordSingleValue[int64](session,
			CheckZoneCodeExistsQuery(code, facilityCode, ""))
		session.Close()

		if existErr == nil && existing > 0 {
			result.Skipped++
			continue
		}

		// handle parent
		var parentUID *string
		if colIndex["parentCode"] != -1 && colIndex["parentCode"] < len(record) {
			parentCode := strings.TrimSpace(record[colIndex["parentCode"]])
			if parentCode != "" {
				pSession, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
				parentZone, pErr := helpers.GetNeo4jSingleRecordAndMapToStruct[models.Zone](pSession,
					GetZoneByCodeAndFacilityQuery(parentCode, facilityCode))
				pSession.Close()

				if pErr != nil {
					result.Errors = append(result.Errors, fmt.Sprintf("row %d: parent zone '%s' not found", rowNum, parentCode))
					continue
				}

				// check parent is root
				if parentZone.ParentUID != nil {
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

func (svc *ZoneService) validateParentIsRoot(parentUID string) error {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	type parentCheck struct {
		UID       string `json:"uid"`
		HasParent bool   `json:"hasParent"`
	}

	check, err := helpers.GetNeo4jSingleRecordAndMapToStruct[parentCheck](session, CheckParentIsRootQuery(parentUID))
	if err != nil {
		return fmt.Errorf("parent zone not found")
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

