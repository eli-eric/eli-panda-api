package systemsService

import (
	"errors"
	"fmt"
	"panda/apigateway/config"
	"panda/apigateway/helpers"
	codebookModels "panda/apigateway/services/codebook-service/models"
	"panda/apigateway/services/systems-service/models"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type SystemsService struct {
	neo4jDriver *neo4j.Driver
	jwtSecret   string
}

type ISystemsService interface {
	GetSystemTypesCodebook(facilityCode string) (result []codebookModels.Codebook, err error)
	GetSystemImportancesCodebook() (result []codebookModels.Codebook, err error)
	GetSystemCriticalitiesCodebook() (result []codebookModels.Codebook, err error)
	GetItemUsagesCodebook() (result []codebookModels.Codebook, err error)
	GetItemConditionsCodebook() (result []codebookModels.Codebook, err error)
	GetLocationAutocompleteCodebook(searchText string, limit int, facilityCode string) (result []codebookModels.Codebook, err error)
	GetZonesCodebook(facilityCode string, searchString string, filter *[]helpers.Filter) (result []codebookModels.Codebook, err error)
	GetSubSystemsByParentUID(parentUID string, facilityCode string) (result []models.System, err error)
	GetSystemImageByUid(uid string) (imageBase64 string, err error)
	GetSystemDetail(uid string, facilityCode string) (result models.System, err error)
	CreateNewSystem(system *models.System, facilityCode string, userUID string) (uid string, err error)
	CreateNewSystemFromJira(facilityCode string, userUID string, userRoles []string, request *models.JiraSystemImportRequest) (result string, err error)
	UpdateSystem(newSystem *models.System, facilityCode string, userUID string) (err error)
	DeleteSystemRecursive(uid, userUid string) (err error)
	GetSystemsAutocompleteCodebook(searchText string, limit int, facilityCode string, filter *[]helpers.Filter) (result []codebookModels.Codebook, err error)
	GetSystemsWithSearchAndPagination(search string, facilityCode string, pagination *helpers.Pagination, sorting *[]helpers.Sorting, filering *[]helpers.ColumnFilter) (result helpers.PaginationResult[models.System], err error)
	GetSystemsForRelationship(search string, facilityCode string, pagination *helpers.Pagination, sorting *[]helpers.Sorting, filering *[]helpers.ColumnFilter, systemFromUid string, relationTypeCode string) (result helpers.PaginationResult[models.System], err error)
	GetSystemRelationships(uid string) (result []models.SystemRelationship, err error)
	DeleteSystemRelationship(uid int64, userUID string) (err error)
	CreateNewSystemRelationship(newRelationship *models.SystemRelationshipRequest, facilityCode string, userUID string) (relId int64, err error)
	GetSystemCode(systemTypeUid, zoneUID, locationUID, parentUID, facilityCode string) (result string, err error)
	GetPhysicalItemProperties(physicalItemUid string) (result []models.PhysicalItemDetail, err error)
	UpdatePhysicalItemProperties(physicalItemUid string, details []models.PhysicalItemDetail, userUID string) (err error)
	GetSystemHistory(uid string) (result []models.SystemHistory, err error)
	GetSystemSparePartsDetail(systemId string, facilityCode string) (result models.SystemSparePartsDetail, err error)
	GetSystemTypeGroups(facilityCode string) (result []codebookModels.Codebook, err error)
	GetSystemTypesBySystemTypeGroup(systemTypeGroupUid, facilityCode string) (result []models.SystemType, err error)
	DeleteSystemTypeGroup(systemTypeGroupUid string) (err error, relatedNodeLabels []helpers.RelatedNodeLabelAmount)
	DeleteSystemType(systemTypeUid string) (err error, relatedNodeLabels []helpers.RelatedNodeLabelAmount)
	CreateSystemTypeGroup(systemTypeGroup *codebookModels.Codebook, facilityCode, userUID string) (err error)
	UpdateSystemTypeGroup(systemTypeGroup *codebookModels.Codebook, userUID string) (err error)
	CreateSystemType(systemType *models.SystemType, facilityCode, userUID, systemTypeGroupUID string) (err error)
	UpdateSystemType(systemType *models.SystemType, facilityCode, userUID string) (err error)
	GetSystemByEun(eun string) (result models.System, err error)
	GetSystemAttributesCodebook(facilityCode string) (result []codebookModels.Codebook, err error)
	GetEuns(facilityCode string) (result []models.EUN, err error)
	SyncSystemLocationByEUNs(euns []models.EunLocation, userUID string) (errs []error)
	GetAllLocationsFlat(facilityCode string) (result []codebookModels.Codebook, err error)
	GetAllSystemTypes() (result []codebookModels.Codebook, err error)
	GetAllZones(facilityCode string) (result []codebookModels.Codebook, err error)
	CreateNewSystemCode(parentUID, systemTypeUID, zoneUID, facilityCode, userUID string) (result models.System, err error)
	GetSystemsForControlsSystems(facilityCode string, pagination *helpers.Pagination, sorting *[]helpers.Sorting, filering *[]helpers.ColumnFilter) (result helpers.PaginationResult[models.SystemCodesResult], err error)
	GetNewSystemCodesPreview(systemTypeUID string, zoneUID string, batch int, facilityCode string) (result []models.SystemCodesResult, err error)
	SaveNewSystemCodes(request *models.SystemCodesRequest, facilityCode string, userUID string) (result []models.SystemCodesResult, err error)
	RecalculateSpareParts() (err error)
	GetSystemsByUids(uids []string) (result []models.System, err error)
	GetSystemsTreeByUids(trees []models.SystemTreeUid) (result []models.System, err error)
	BuildSystemHierarchy(tree models.SystemTreeUid) (*models.System, error)
	MovePhysicalItem(movement *models.PhysicalItemMovement, userUID, facilityCode string) (destinationSystemUid string, err error)
	ReplacePhysicalItems(movement *models.PhysicalItemMovement, userUID, facilityCode string) (destinationSystemUid string, err error)
	MoveSystems(movement *models.SystemsMovement, userUID string) (destinationSystemUid string, err error)
	GetPhysicalItemsBySystemUidRecursive(systemUid string) (result []models.SystemPhysicalItemInfo, err error)
	AssignSpareItem(request models.AssignSpareRequest, userUID string) (models.AssignSpareResponse, error)
}

// Create new security service instance
func NewSystemsService(settings *config.Config, driver *neo4j.Driver) ISystemsService {

	return &SystemsService{neo4jDriver: driver, jwtSecret: settings.JwtSecret}
}

func (svc *SystemsService) GetSystemTypesCodebook(facilityCode string) (result []codebookModels.Codebook, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetSystemTypesCodebookQuery(facilityCode)
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	return result, err
}

func (svc *SystemsService) GetSystemImportancesCodebook() (result []codebookModels.Codebook, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetSystemImportancesCodebookQuery()
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	return result, err
}

func (svc *SystemsService) GetSystemCriticalitiesCodebook() (result []codebookModels.Codebook, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetSystemCriticalityCodebookQuery()
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	return result, err
}

func (svc *SystemsService) GetItemUsagesCodebook() (result []codebookModels.Codebook, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetItemUsagesCodebookQuery()
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	return result, err
}

func (svc *SystemsService) GetItemConditionsCodebook() (result []codebookModels.Codebook, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetItemConditionsCodebookQuery()
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	return result, err
}

func (svc *SystemsService) GetLocationAutocompleteCodebook(searchText string, limit int, facilityCode string) (result []codebookModels.Codebook, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetLocationsBySearchTextQuery(searchText, limit, facilityCode)
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	return result, err
}

func (svc *SystemsService) GetZonesCodebook(facilityCode string, searchString string, filter *[]helpers.Filter) (result []codebookModels.Codebook, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	onlyRootElements := false
	if filter != nil {
		for _, f := range *filter {
			if f.Key != "onlyRootElements" {
				continue
			}

			switch v := f.Value.(type) {
			case bool:
				onlyRootElements = v
			case string:
				onlyRootElements = strings.EqualFold(v, "true") || v == "1"
			case float64:
				onlyRootElements = v != 0
			}
			break
		}
	}

	query := GetZonesCodebookQuery(facilityCode, searchString, onlyRootElements)
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	return result, err
}

func (svc *SystemsService) GetSubSystemsByParentUID(parentUID string, facilityCode string) (result []models.System, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetSubSystemsQuery(parentUID, facilityCode)
	result, err = helpers.GetNeo4jArrayOfNodes[models.System](session, query)

	return result, err
}

func (svc *SystemsService) GetSystemImageByUid(uid string) (result string, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := SystemImageByUidQuery(uid)
	result, err = helpers.GetNeo4jSingleRecordSingleValue[string](session, query)

	return result, err
}

func (svc *SystemsService) GetSystemDetail(uid string, facilityCode string) (result models.System, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	result, err = helpers.GetNeo4jSingleRecordAndMapToStruct[models.System](session, SystemDetailQuery(uid, facilityCode))

	return result, err
}

func (svc *SystemsService) CreateNewSystem(system *models.System, facilityCode string, userUID string) (uid string, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	uid, err = helpers.WriteNeo4jAndReturnSingleValue[string](session, CreateNewSystemQuery(system, facilityCode, userUID))

	if err != nil {
		log.Info().Msg(err.Error())
	}
	// we will do that later, more precise
	// else {
	// 	go func() {
	// 		helpers.LogDBHistory(session, uid, nil, system, userUID, helpers.DB_LOG_CREATE)
	// 	}()
	// }

	return uid, err
}

func (svc *SystemsService) CreateNewSystemFromJira(facilityCode string, userUID string, userRoles []string, request *models.JiraSystemImportRequest) (result string, err error) {

	// check if system code already exists
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	exists, err := helpers.GetNeo4jSingleRecordSingleValue[bool](session, checkSystemCodeExistsQuery(request.Code))
	defer session.Close()

	if err != nil {
		return "", err
	}

	if exists {
		return "", helpers.ERR_DUPLICATE_SYSTEM_CODE
	}

	// create new system
	result, err = helpers.WriteNeo4jAndReturnSingleValue[string](session, CreateNewSystemFromJiraQuery(request, facilityCode, userUID))

	return result, err
}

func (svc *SystemsService) UpdateSystem(system *models.System, facilityCode string, userUID string) (err error) {

	if system != nil {
		session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

		oldSystem, err := helpers.GetNeo4jSingleRecordAndMapToStruct[models.System](session, SystemDetailQuery(system.UID, facilityCode))

		if err == nil {
			_, err = helpers.WriteNeo4jAndReturnSingleValue[string](session, UpdateSystemQuery(system, &oldSystem, facilityCode, userUID))
		}

		if err != nil {
			log.Info().Msg(err.Error())
		}
		// we will do that later, more precise
		// else {
		// 	go func() {
		// 		helpers.LogDBHistory(session, oldSystem.UID, oldSystem, system, userUID, helpers.DB_LOG_UPDATE)
		// 	}()
		// }

	} else {
		err = helpers.ERR_INVALID_INPUT
	}
	return err
}

func (svc *SystemsService) DeleteSystemRecursive(uid, userUid string) (err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := DeleteSystemByUidQuery(uid, userUid)
	err = helpers.WriteNeo4jAndReturnNothing(session, query)

	return err
}

func (svc *SystemsService) GetPhysicalItemsBySystemUidRecursive(systemUid string) (result []models.SystemPhysicalItemInfo, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetPhysicalItemsBySystemUidRecursiveQuery(systemUid)
	result, err = helpers.GetNeo4jArrayOfNodes[models.SystemPhysicalItemInfo](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}

func (svc *SystemsService) GetSystemsAutocompleteCodebook(searchText string, limit int, facilityCode string, filter *[]helpers.Filter) (result []codebookModels.Codebook, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	onlyTechnologicalUnits := true

	if filter != nil {
		for _, f := range *filter {
			if f.Key == "technologicalUnits" {
				onlyTechnologicalUnits = f.Value.(bool)
			}
		}
	}

	query := GetSystemsForAutocomplete(searchText, limit, facilityCode, onlyTechnologicalUnits)
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}

func (svc *SystemsService) GetSystemsWithSearchAndPagination(search string, facilityCode string, pagination *helpers.Pagination, sorting *[]helpers.Sorting, filering *[]helpers.ColumnFilter) (result helpers.PaginationResult[models.System], err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetSystemsBySearchTextFullTextQuery(search, facilityCode, pagination, sorting, filering)
	items, err := helpers.GetNeo4jArrayOfNodes[models.System](session, query)
	totalCount, _ := helpers.GetNeo4jSingleRecordSingleValue[int64](session, GetSystemsBySearchTextFullTextCountQuery(search, facilityCode, filering))

	result = helpers.GetPaginationResult(items, int64(totalCount), err)

	return result, err
}

func (svc *SystemsService) GetSystemsForRelationship(search string, facilityCode string, pagination *helpers.Pagination, sorting *[]helpers.Sorting, filering *[]helpers.ColumnFilter, systemFromUid string, relationTypeCode string) (result helpers.PaginationResult[models.System], err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetSystemsBySearchTextFullTextQuery(search, facilityCode, pagination, sorting, filering)
	items, err := helpers.GetNeo4jArrayOfNodes[models.System](session, query)
	totalCount, _ := helpers.GetNeo4jSingleRecordSingleValue[int64](session, GetSystemsBySearchTextFullTextCountQuery(search, facilityCode, filering))

	result = helpers.GetPaginationResult(items, int64(totalCount), err)

	return result, err
}

func (svc *SystemsService) GetSystemRelationships(uid string) (result []models.SystemRelationship, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetSystemRelationshipsQuery(uid)
	result, err = helpers.GetNeo4jArrayOfNodes[models.SystemRelationship](session, query)

	return result, err
}

func (svc *SystemsService) DeleteSystemRelationship(uid int64, userUID string) (err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := DeleteSystemRelationshipQuery(uid, userUID)
	err = helpers.WriteNeo4jAndReturnNothing(session, query)

	if err != nil {
		log.Info().Msg(err.Error())
	}

	return err
}

func (svc *SystemsService) CreateNewSystemRelationship(newRelationship *models.SystemRelationshipRequest, facilityCode string, userUID string) (relId int64, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	relId, err = helpers.WriteNeo4jAndReturnSingleValue[int64](session, CreateNewSystemRelationshipQuery(newRelationship, facilityCode, userUID))

	if err != nil {
		log.Info().Msg(err.Error())
	}

	return relId, err
}

func (svc *SystemsService) GetSystemCode(systemTypeUid, zoneUID, locationUID, parentUID, facilityCode string) (result string, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	if systemTypeUid == "" {
		err = errors.New("missing system type")
		return "", err
	}

	mask, err := helpers.GetNeo4jSingleRecordSingleValue[string](session, GetSystemTypeMask(systemTypeUid, facilityCode))

	if err != nil {
		log.Error().Msg(err.Error())
		return "", err
	}

	if !strings.Contains(mask, SYSTEM_CODE_GENERATE_SERIAL_PREFIX) {
		err = errors.New("missing serial number information")
		return "", err
	}

	systemTypeCode, err := helpers.GetNeo4jSingleRecordSingleValue[string](session, GetSystemTypeCode(systemTypeUid))

	if err != nil {
		log.Error().Msg(err.Error())
		return "", err
	}

	mask = strings.ReplaceAll(mask, SYSTEM_CODE_GENERATE_SYSTEM_TYPE_CODE, systemTypeCode)

	if strings.Contains(mask, SYSTEM_CODE_GENERATE_ZONE_CODE) || strings.Contains(mask, SYSTEM_CODE_GENERATE_ZONE_NAME) || strings.Contains(mask, SYSTEM_CODE_GENERATE_SUB_ZONE_CODE) {
		if zoneUID == "" {
			err = errors.New("missing zone")
			return "", err
		} else {
			zoneCode, err := helpers.GetNeo4jSingleRecordSingleValue[string](session, GetZoneCode(zoneUID))
			if err != nil {
				log.Error().Msg(err.Error())
				return "", err
			}
			zoneName, err := helpers.GetNeo4jSingleRecordSingleValue[string](session, GetZoneName(zoneUID))
			if err != nil {
				log.Error().Msg(err.Error())
				return "", err
			}
			subZoneCode, err := helpers.GetNeo4jSingleRecordSingleValue[string](session, GetSubZoneCode(zoneUID))
			if err != nil {
				log.Error().Msg(err.Error())
				return "", err
			}

			mask = strings.ReplaceAll(mask, SYSTEM_CODE_GENERATE_ZONE_CODE, zoneCode)
			mask = strings.ReplaceAll(mask, SYSTEM_CODE_GENERATE_ZONE_NAME, zoneName)
			mask = strings.ReplaceAll(mask, SYSTEM_CODE_GENERATE_SUB_ZONE_CODE, subZoneCode)
		}
	}

	if strings.Contains(mask, SYSTEM_CODE_GENERATE_LOCATION_CODE) || strings.Contains(mask, SYSTEM_CODE_GENERATE_LOCATION_NAME) {
		if locationUID == "" {
			err = errors.New("missing location")
			return "", err
		} else {
			locationCode, err := helpers.GetNeo4jSingleRecordSingleValue[string](session, GetLocationCode(locationUID))
			if err != nil {
				log.Error().Msg(err.Error())
				return "", err
			}
			locationName, err := helpers.GetNeo4jSingleRecordSingleValue[string](session, GetLocationName(locationUID))
			if err != nil {
				log.Error().Msg(err.Error())
				return "", err
			}

			mask = strings.ReplaceAll(mask, SYSTEM_CODE_GENERATE_LOCATION_CODE, locationCode)
			mask = strings.ReplaceAll(mask, SYSTEM_CODE_GENERATE_LOCATION_NAME, locationName)
		}
	}

	if strings.Contains(mask, SYSTEM_CODE_GENERATE_FACILITY_CODE) {
		if facilityCode == "" {
			err = errors.New("missing facility")
			return "", err
		} else {
			mask = strings.ReplaceAll(mask, SYSTEM_CODE_GENERATE_FACILITY_CODE, facilityCode)
		}
	}

	// now take the mask and split it by serial info number which starts with {serial( - could be for example {serial(3)}

	serialPrefix := SYSTEM_CODE_GENERATE_SERIAL_PREFIX
	serialPrefixIndex := strings.Index(mask, serialPrefix)
	maskWithoutSerial := mask[:serialPrefixIndex]
	serialNumberLengthString := mask[serialPrefixIndex+len(serialPrefix):]
	serialNumberLengthString = serialNumberLengthString[:len(serialNumberLengthString)-2]
	serialNumberLength, err := strconv.ParseInt(serialNumberLengthString, 10, 64)

	if err != nil {
		log.Error().Msg(err.Error())
		return "", err
	}

	result, err = helpers.GetNeo4jSingleRecordSingleValue[string](session, GetNewSystemCode(maskWithoutSerial, int(serialNumberLength), facilityCode))

	if err != nil {
		log.Error().Msg(err.Error())
		return "", err
	}

	return result, err
}

func (svc *SystemsService) getSystemCodePrefixAndSerialLength(session neo4j.Session, systemTypeUid, zoneUID, locationUID, facilityCode string) (systemCodePrefix string, serialNumberLength int, err error) {

	if systemTypeUid == "" {
		return "", 0, errors.New("missing system type")
	}

	mask, err := helpers.GetNeo4jSingleRecordSingleValue[string](session, GetSystemTypeMask(systemTypeUid, facilityCode))
	if err != nil {
		return "", 0, err
	}

	if !strings.Contains(mask, SYSTEM_CODE_GENERATE_SERIAL_PREFIX) {
		return "", 0, errors.New("missing serial number information")
	}

	systemTypeCode, err := helpers.GetNeo4jSingleRecordSingleValue[string](session, GetSystemTypeCode(systemTypeUid))
	if err != nil {
		return "", 0, err
	}

	mask = strings.ReplaceAll(mask, SYSTEM_CODE_GENERATE_SYSTEM_TYPE_CODE, systemTypeCode)

	if strings.Contains(mask, SYSTEM_CODE_GENERATE_ZONE_CODE) || strings.Contains(mask, SYSTEM_CODE_GENERATE_ZONE_NAME) || strings.Contains(mask, SYSTEM_CODE_GENERATE_SUB_ZONE_CODE) {
		if zoneUID == "" {
			return "", 0, errors.New("missing zone")
		}
		zoneCode, err := helpers.GetNeo4jSingleRecordSingleValue[string](session, GetZoneCode(zoneUID))
		if err != nil {
			return "", 0, err
		}
		zoneName, err := helpers.GetNeo4jSingleRecordSingleValue[string](session, GetZoneName(zoneUID))
		if err != nil {
			return "", 0, err
		}
		subZoneCode, err := helpers.GetNeo4jSingleRecordSingleValue[string](session, GetSubZoneCode(zoneUID))
		if err != nil {
			return "", 0, err
		}

		mask = strings.ReplaceAll(mask, SYSTEM_CODE_GENERATE_ZONE_CODE, zoneCode)
		mask = strings.ReplaceAll(mask, SYSTEM_CODE_GENERATE_ZONE_NAME, zoneName)
		mask = strings.ReplaceAll(mask, SYSTEM_CODE_GENERATE_SUB_ZONE_CODE, subZoneCode)
	}

	if strings.Contains(mask, SYSTEM_CODE_GENERATE_LOCATION_CODE) || strings.Contains(mask, SYSTEM_CODE_GENERATE_LOCATION_NAME) {
		if locationUID == "" {
			return "", 0, errors.New("missing location")
		}
		locationCode, err := helpers.GetNeo4jSingleRecordSingleValue[string](session, GetLocationCode(locationUID))
		if err != nil {
			return "", 0, err
		}
		locationName, err := helpers.GetNeo4jSingleRecordSingleValue[string](session, GetLocationName(locationUID))
		if err != nil {
			return "", 0, err
		}

		mask = strings.ReplaceAll(mask, SYSTEM_CODE_GENERATE_LOCATION_CODE, locationCode)
		mask = strings.ReplaceAll(mask, SYSTEM_CODE_GENERATE_LOCATION_NAME, locationName)
	}

	if strings.Contains(mask, SYSTEM_CODE_GENERATE_FACILITY_CODE) {
		if facilityCode == "" {
			return "", 0, errors.New("missing facility")
		}
		mask = strings.ReplaceAll(mask, SYSTEM_CODE_GENERATE_FACILITY_CODE, facilityCode)
	}

	serialPrefixIndex := strings.Index(mask, SYSTEM_CODE_GENERATE_SERIAL_PREFIX)
	maskWithoutSerial := mask[:serialPrefixIndex]
	serialNumberLengthString := mask[serialPrefixIndex+len(SYSTEM_CODE_GENERATE_SERIAL_PREFIX):]
	serialNumberLengthString = serialNumberLengthString[:len(serialNumberLengthString)-2]
	parsedSerialLength, err := strconv.ParseInt(serialNumberLengthString, 10, 64)
	if err != nil {
		return "", 0, err
	}

	return maskWithoutSerial, int(parsedSerialLength), nil
}

func (svc *SystemsService) GetSystemsForControlsSystems(facilityCode string, pagination *helpers.Pagination, sorting *[]helpers.Sorting, filering *[]helpers.ColumnFilter) (result helpers.PaginationResult[models.SystemCodesResult], err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetSystemsForControlsSystemsQuery(facilityCode, pagination, sorting, filering)
	items, err := helpers.GetNeo4jArrayOfNodes[models.SystemCodesResult](session, query)
	totalCount, _ := helpers.GetNeo4jSingleRecordSingleValue[int64](session, GetSystemsForControlsSystemsCountQuery(facilityCode, filering))

	result = helpers.GetPaginationResult(items, int64(totalCount), err)
	return result, err
}

func (svc *SystemsService) GetNewSystemCodesPreview(systemTypeUID string, zoneUID string, batch int, facilityCode string) (result []models.SystemCodesResult, err error) {

	if batch <= 0 {
		return nil, helpers.ERR_INVALID_INPUT
	}

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	prefix, serialLen, err := svc.getSystemCodePrefixAndSerialLength(session, systemTypeUID, zoneUID, "", facilityCode)
	if err != nil {
		return nil, err
	}

	query := GetNewSystemCodesPreviewQuery(systemTypeUID, zoneUID, prefix, serialLen, batch, facilityCode)
	result, err = helpers.GetNeo4jArrayOfNodes[models.SystemCodesResult](session, query)

	helpers.ProcessArrayResult(&result, err)
	return result, err
}

func (svc *SystemsService) SaveNewSystemCodes(request *models.SystemCodesRequest, facilityCode string, userUID string) (result []models.SystemCodesResult, err error) {

	if request == nil || request.SystemType == nil || request.Zone == nil {
		return nil, helpers.ERR_INVALID_INPUT
	}
	if request.Batch <= 0 {
		return nil, helpers.ERR_INVALID_INPUT
	}

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	prefix, serialLen, err := svc.getSystemCodePrefixAndSerialLength(session, request.SystemType.UID, request.Zone.UID, "", facilityCode)
	if err != nil {
		return nil, err
	}

	query := SaveNewSystemCodesQuery(request.SystemType.UID, request.Zone.UID, prefix, serialLen, request.Batch, facilityCode, userUID)
	result, err = helpers.WriteNeo4jAndReturnArrayOfNodes[models.SystemCodesResult](session, query)
	if err != nil {
		//log the error
		log.Error().Msg(fmt.Sprintf("Error saving new system codes: %v", err))
		// translate APOC validation error into a stable message for API client
		if strings.Contains(err.Error(), "missing default parent system for selected zone") {
			return nil, errors.New("missing default parent system for selected zone")
		}
	}

	helpers.ProcessArrayResult(&result, err)
	return result, err
}

const SYSTEM_CODE_GENERATE_ZONE_CODE = "{ZC}"
const SYSTEM_CODE_GENERATE_SUB_ZONE_CODE = "{SZC}"
const SYSTEM_CODE_GENERATE_ZONE_NAME = "{ZN}"
const SYSTEM_CODE_GENERATE_LOCATION_CODE = "{LC}"
const SYSTEM_CODE_GENERATE_LOCATION_NAME = "{LN}"
const SYSTEM_CODE_GENERATE_FACILITY_CODE = "{FC}"
const SYSTEM_CODE_GENERATE_SYSTEM_TYPE_CODE = "{STC}"
const SYSTEM_CODE_GENERATE_SERIAL_PREFIX = "{serial("

func (svc *SystemsService) GetPhysicalItemProperties(physicalItemUid string) (result []models.PhysicalItemDetail, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetPhysicalItemPropertiesQuery(physicalItemUid)
	result, err = helpers.GetNeo4jArrayOfNodes[models.PhysicalItemDetail](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}

func (svc *SystemsService) UpdatePhysicalItemProperties(physicalItemUid string, details []models.PhysicalItemDetail, userUID string) (err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	err = helpers.WriteNeo4jAndReturnNothing(session, UpdatePhysicalItemDetailsQuery(physicalItemUid, details, userUID))

	return err
}

func (svc *SystemsService) GetSystemHistory(uid string) (result []models.SystemHistory, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetSystemHistoryQuery(uid)
	result, err = helpers.GetNeo4jArrayOfNodes[models.SystemHistory](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}

func (svc *SystemsService) GetSystemTypeGroups(facilityCode string) (result []codebookModels.Codebook, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetSystemTypeGroupsQuery(facilityCode)
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}

func (svc *SystemsService) GetSystemTypesBySystemTypeGroup(systemTypeGroupUid, facilityCode string) (result []models.SystemType, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetSystemTypesBySystemTypeGroupQuery(systemTypeGroupUid, facilityCode)
	result, err = helpers.GetNeo4jArrayOfNodes[models.SystemType](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}

func (svc *SystemsService) DeleteSystemTypeGroup(systemTypeGroupUid string) (err error, relatedNodeLabels []helpers.RelatedNodeLabelAmount) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	relatedLabels, err := helpers.GetNeo4jArrayOfNodes[helpers.RelatedNodeLabelAmount](session, GetSystemTypeGroupRelatedNodeLabelsCountQuery(systemTypeGroupUid))

	if err != nil {
		return err, relatedLabels
	}

	if len(relatedLabels) > 0 {
		return helpers.ERR_DELETE_RELATED_ITEMS, relatedLabels
	}

	err = helpers.WriteNeo4jAndReturnNothing(session, DeleteSystemTypeGroupQuery(systemTypeGroupUid))

	return err, relatedLabels
}

func (svc *SystemsService) DeleteSystemType(systemTypeUid string) (err error, relatedNodeLabels []helpers.RelatedNodeLabelAmount) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	relatedLabels, err := helpers.GetNeo4jArrayOfNodes[helpers.RelatedNodeLabelAmount](session, GetSystemTypeRelatedNodeLabelsCountQuery(systemTypeUid))

	if err != nil {
		return err, relatedLabels
	}

	if len(relatedLabels) > 0 {
		return helpers.ERR_DELETE_RELATED_ITEMS, relatedLabels
	}

	err = helpers.WriteNeo4jAndReturnNothing(session, DeleteSystemTypeQuery(systemTypeUid))

	return err, relatedLabels
}

func (svc *SystemsService) CreateSystemTypeGroup(systemTypeGroup *codebookModels.Codebook, facilityCode, userUID string) (err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	err = helpers.WriteNeo4jAndReturnNothing(session, CreateSystemTypeGroupQuery(systemTypeGroup, facilityCode, userUID))

	return err
}

func (svc *SystemsService) UpdateSystemTypeGroup(systemTypeGroup *codebookModels.Codebook, userUID string) (err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	err = helpers.WriteNeo4jAndReturnNothing(session, UpdateSystemTypeGroupQuery(systemTypeGroup, userUID))

	return err
}

func (svc *SystemsService) CreateSystemType(systemType *models.SystemType, facilityCode, userUID, systemTypeGroupUID string) (err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	err = helpers.WriteNeo4jAndReturnNothing(session, CreateSystemTypeQuery(systemType, facilityCode, userUID, systemTypeGroupUID))

	return err
}

func (svc *SystemsService) UpdateSystemType(systemType *models.SystemType, facilityCode, userUID string) (err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	err = helpers.WriteNeo4jAndReturnNothing(session, UpdateSystemTypeQuery(systemType, facilityCode, userUID))

	return err
}

func (svc *SystemsService) GetSystemByEun(eun string) (result models.System, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	result, err = helpers.GetNeo4jSingleRecordAndMapToStruct[models.System](session, GetSystemByEunQuery(eun))

	return result, err
}

func (svc *SystemsService) GetSystemAttributesCodebook(facilityCode string) (result []codebookModels.Codebook, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetSystemAttributesCodebookQuery(facilityCode)
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}

func (svc *SystemsService) GetEuns(facilityCode string) (result []models.EUN, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetEunsQuery(facilityCode)
	result, err = helpers.GetNeo4jArrayOfNodes[models.EUN](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}

func (svc *SystemsService) SyncSystemLocationByEUNs(euns []models.EunLocation, userUID string) (errs []error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	for _, eun := range euns {
		err := helpers.WriteNeo4jAndReturnNothing(session, SyncSystemLocationByEUNQuery(eun.EUN, eun.LocationUID, userUID))
		if err != nil {
			errs = append(errs, err)
			log.Error().Msg(err.Error())
		}
	}

	return errs
}

func (svc *SystemsService) GetAllLocationsFlat(facilityCode string) (result []codebookModels.Codebook, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetAllLocationsFlatQuery(facilityCode)
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}

func (svc *SystemsService) GetAllSystemTypes() (result []codebookModels.Codebook, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetAllSystemTypesQuery()
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}

func (svc *SystemsService) GetAllZones(facilityCode string) (result []codebookModels.Codebook, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetAllZonesQuery(facilityCode)
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}

func (svc *SystemsService) CreateNewSystemCode(parentUID, systemTypeUID, zoneUID, facilityCode, userUID string) (result models.System, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	if parentUID == "" {
		err = errors.New("missing parent")
		return result, err
	}

	newCode, err := svc.GetSystemCode(systemTypeUID, zoneUID, "", parentUID, facilityCode)

	if err != nil {
		log.Error().Msg(err.Error())
		return result, err
	}

	newSystemUid, err := helpers.WriteNeo4jAndReturnSingleValue[string](session, CreateNewSystemForNewSystemCodeQuery(parentUID, newCode, systemTypeUID, zoneUID, facilityCode, userUID))

	if err != nil {
		log.Error().Msg(err.Error())
		return result, err
	}

	result, err = helpers.GetNeo4jSingleRecordAndMapToStruct[models.System](session, SystemDetailQuery(newSystemUid, facilityCode))

	return result, err
}

func (svc *SystemsService) RecalculateSpareParts() (err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	err = helpers.WriteNeo4jAndReturnNothing(session, RecalculateSparePartsCoverageCoeficientQuery())

	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	err = helpers.WriteNeo4jAndReturnNothing(session, RecalculateSystemSparePartsCoverageQuery())

	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	err = helpers.WriteNeo4jAndReturnNothing(session, RecalculateSystemSparePartsEmptyCoverageQuery())

	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	return err
}

func (svc *SystemsService) GetSystemsByUids(uids []string) (result []models.System, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetSystemsByUidsQuery(uids)
	result, err = helpers.GetNeo4jArrayOfNodes[models.System](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}

func (svc *SystemsService) GetSystemsTreeByUids(trees []models.SystemTreeUid) (result []models.System, err error) {

	// Step 1: Collect all UIDs
	uids := collectUids(trees)

	// Step 2: Retrieve all systems from the database in one query
	systems, err := svc.GetSystemsByUids(uids)
	if err != nil {
		return nil, err
	}

	// Step 3: Create a map of UID -> System for quick lookup
	systemMap := make(map[string]models.System)
	for _, system := range systems {
		systemMap[system.UID] = system
	}

	// Step 4: Build the system hierarchy
	for _, tree := range trees {
		system, err := buildSystemHierarchy(tree, systemMap)
		if err != nil {
			return nil, err
		}
		result = append(result, *system)
	}

	return result, err
}

// Recursive function to map SystemTreeUid to System structure
func (svc *SystemsService) BuildSystemHierarchy(tree models.SystemTreeUid) (*models.System, error) {
	// Retrieve the system for the current node's UID
	systems, err := svc.GetSystemsByUids([]string{tree.UID})
	if err != nil || len(systems) == 0 {
		log.Error().Msg(err.Error())
		return nil, err
	}
	system := systems[0]

	// Recursively build children
	if tree.Children != nil {
		var subSystems []models.System
		for _, childTree := range *tree.Children {
			childSystem, err := svc.BuildSystemHierarchy(childTree)
			if err != nil {
				return nil, err
			}
			subSystems = append(subSystems, *childSystem)
		}
		system.SubSystems = &subSystems
	}

	return &system, nil
}

// Recursive function to map SystemTreeUid to System structure using preloaded Systems
func buildSystemHierarchy(tree models.SystemTreeUid, systemMap map[string]models.System) (*models.System, error) {
	// Retrieve the system for the current node's UID from the preloaded map
	system, exists := systemMap[tree.UID]
	if !exists {
		return nil, fmt.Errorf("system with UID %s not found", tree.UID)
	}

	// Recursively build children
	if tree.Children != nil {
		var subSystems []models.System
		for _, childTree := range *tree.Children {
			childSystem, err := buildSystemHierarchy(childTree, systemMap)
			if err != nil {
				return nil, err
			}
			subSystems = append(subSystems, *childSystem)
		}
		system.SubSystems = &subSystems
	}

	return &system, nil
}

// Helper function to collect all UIDs from the tree (including root and children)
func collectUids(trees []models.SystemTreeUid) []string {
	var uids []string
	var collect func(tree models.SystemTreeUid)

	collect = func(tree models.SystemTreeUid) {
		uids = append(uids, tree.UID)
		if tree.Children != nil {
			for _, child := range *tree.Children {
				collect(child)
			}
		}
	}

	for _, tree := range trees {
		collect(tree)
	}

	return uids
}

func (svc *SystemsService) MovePhysicalItem(movement *models.PhysicalItemMovement, userUID, facilityCode string) (destinationSystemUid string, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	err = ValidatePhysicalItemMovement(movement)

	if err != nil {
		return "", err
	}

	// check if the destination system, if it already exists, does not have physical item assigned,
	// because in this version we can move physical item only to the system that does not have any physical item assigned
	if movement.DestinationSystemUID != "" {
		containsItem, err := helpers.GetNeo4jSingleRecordSingleValue[bool](session, HasSystemPhysicalItemQuery(movement.DestinationSystemUID))
		if err != nil {
			return "", err
		}
		if containsItem {
			return "", errors.New("destination system already contains physical item")
		}
	}

	destinationSystemUid, err = helpers.WriteNeo4jAndReturnSingleValue[string](session, MovePhysicalItemQuery(movement, userUID))

	if err == nil {
		queries := []helpers.DatabaseQuery{
			CopySystemTypeQuery(movement.SourceSystemUID, destinationSystemUid),
			CopySystemResponsibleQuery(movement.SourceSystemUID, destinationSystemUid),
			CopySystemResponsibleTeamQuery(movement.SourceSystemUID, destinationSystemUid),
			RecordItemMoveHistoryQuery(userUID, movement.SourceSystemUID, destinationSystemUid),
			SetMissingFacilityToSystems(facilityCode),
		}
		err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session, queries...)
	}

	return destinationSystemUid, err
}

func ValidatePhysicalItemMovement(movement *models.PhysicalItemMovement) error {

	// if the destination system is new there has to be parent system specified
	if movement.DestinationSystemUID == "" && movement.ParentSystemUID == "" {
		return errors.New("missing parent system")
	}

	if movement.SystemName == "" {
		return errors.New("missing system name")
	}

	return nil
}

func (svc *SystemsService) ReplacePhysicalItems(movement *models.PhysicalItemMovement, userUID, facilityCode string) (destinationSystemUid string, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	err = ValidatePhysicalItemReplacement(movement)

	if err != nil {
		return "", err
	}

	oldSystemUid, err := helpers.WriteNeo4jAndReturnSingleValue[string](session, ReplacePhysicalItemsQuery(movement, userUID, facilityCode))

	if err == nil {
		queries := []helpers.DatabaseQuery{
			CopySystemTypeQuery(movement.SourceSystemUID, destinationSystemUid),
			CopySystemTypeQuery(movement.DestinationSystemUID, oldSystemUid),
			CopySystemResponsibleQuery(movement.SourceSystemUID, destinationSystemUid),
			CopySystemResponsibleQuery(movement.DestinationSystemUID, oldSystemUid),
			CopySystemResponsibleTeamQuery(movement.SourceSystemUID, destinationSystemUid),
			CopySystemResponsibleTeamQuery(movement.DestinationSystemUID, oldSystemUid),
			RecordItemMoveHistoryQuery(userUID, movement.SourceSystemUID, destinationSystemUid),
			RecordItemMoveHistoryQuery(userUID, movement.DestinationSystemUID, oldSystemUid),
		}
		err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session, queries...)
	}

	return destinationSystemUid, err
}

func ValidatePhysicalItemReplacement(movement *models.PhysicalItemMovement) error {

	erros := make([]string, 0)
	if movement.DestinationSystemUID == "" {
		erros = append(erros, "missing destination system")
	}

	if movement.ParentSystemUID == "" {
		erros = append(erros, "missing parent system")
	}

	if movement.SourceSystemUID == "" {
		erros = append(erros, "missing source system")
	}

	if len(erros) > 0 {
		resultErr := strings.Join(erros[:], ", ")
		return errors.New(resultErr)
	}

	return nil
}

func (svc *SystemsService) MoveSystems(movement *models.SystemsMovement, userUID string) (destinationSystemUid string, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	err = ValidateSystemsMovement(movement)

	if err != nil {
		return "", err
	}

	destinationSystemUid, err = helpers.WriteNeo4jAndReturnSingleValue[string](session, MoveSystemsQuery(movement, userUID))

	return destinationSystemUid, err
}

func ValidateSystemsMovement(movement *models.SystemsMovement) error {

	erros := make([]string, 0)
	if movement.TargetParentSystemUid == "" {
		erros = append(erros, "missing destination parent system")
	}

	if len(movement.SystemsToMoveUids) == 0 {
		erros = append(erros, "missing systems to move")
	}

	if len(erros) > 0 {
		resultErr := strings.Join(erros[:], ", ")
		return errors.New(resultErr)
	}

	return nil
}

func (svc *SystemsService) AssignSpareItem(request models.AssignSpareRequest, userUID string) (models.AssignSpareResponse, error) {
	session, err := helpers.NewNeo4jSession(*svc.neo4jDriver)
	if err != nil {
		return models.AssignSpareResponse{}, err
	}
	defer session.Close()

	var response models.AssignSpareResponse
	response.UpdatedSystemUid = request.SystemUid

	// Step 1: Validation - Check if system exists and has at most one item
	validateSystemQuery := `
        MATCH (s:System {uid: $systemUid})
        OPTIONAL MATCH (s)-[:CONTAINS_ITEM]->(currentItem:Item)
        RETURN s.uid as systemUid, count(currentItem) as itemCount, collect(currentItem.uid) as currentItems
    `

	result, err := session.Run(validateSystemQuery, map[string]interface{}{
		"systemUid": request.SystemUid,
	})
	if err != nil {
		return response, err
	}

	record, err := result.Single()
	if err != nil {
		return response, errors.New("system not found")
	}

	itemCount := record.Values[1].(int64)
	if itemCount > 1 {
		return response, errors.New("system contains more than one item")
	}

	var currentItemUid string
	if itemCount == 1 {
		currentItems := record.Values[2].([]interface{})
		if len(currentItems) > 0 {
			currentItemUid = currentItems[0].(string)
		}
	}

	// Step 2: Validate spare item relationship and get spare system UID
	validateSpareQuery := `
        MATCH (spareItem:Item {uid: $spareItemUid})<-[:CONTAINS_ITEM]-(spareSystem:System)-[:IS_SPARE_FOR]->(targetSystem:System {uid: $systemUid})
        RETURN spareItem.uid as spareUid, spareSystem.uid as spareSystemUid
    `

	result, err = session.Run(validateSpareQuery, map[string]interface{}{
		"spareItemUid": request.SpareItemUid,
		"systemUid":    request.SystemUid,
	})
	if err != nil {
		return response, err
	}

	if !result.Next() {
		return response, errors.New("spare item does not have IS_SPARE_FOR relationship with target system")
	}

	spareSystemRecord := result.Record()
	spareSystemUid := spareSystemRecord.Values[1].(string)

	// Step 3: Execute the assignment in a transaction
	queries := []helpers.DatabaseQuery{}

	// Get current item's usage code if it exists
	var currentItemUsageCode string
	if currentItemUid != "" {
		getCurrentItemUsageQuery := `
			MATCH (item:Item {uid: $itemUid})
			OPTIONAL MATCH (item)-[:HAS_ITEM_USAGE]->(usage:ItemUsage)
			RETURN usage.code as usageCode
		`
		usageResult, err := session.Run(getCurrentItemUsageQuery, map[string]interface{}{
			"itemUid": currentItemUid,
		})
		if err == nil && usageResult.Next() {
			usageRecord := usageResult.Record()
			if usageRecord.Values[0] != nil {
				currentItemUsageCode = usageRecord.Values[0].(string)
			}
		}
	}

	// FIRST: Detach spare item from spare system (CRITICAL FIX)
	detachSpareItemQuery := helpers.DatabaseQuery{
		Query: `MATCH (spareSystem:System {uid: $spareSystemUid})-[r:CONTAINS_ITEM]->(spareItem:Item {uid: $spareItemUid}) DELETE r`,
		Parameters: map[string]interface{}{
			"spareSystemUid": spareSystemUid,
			"spareItemUid":   request.SpareItemUid,
		},
	}
	queries = append(queries, detachSpareItemQuery)

	// THEN: Assign spare item to target system
	assignSpareQuery := helpers.DatabaseQuery{
		Query: `MATCH (s:System {uid: $systemUid})
                MATCH (spareItem:Item {uid: $spareItemUid})
                CREATE (s)-[:CONTAINS_ITEM]->(spareItem)
                SET spareItem.status = 'in_system', spareItem.updatedAt = datetime()`,
		Parameters: map[string]interface{}{
			"systemUid":    request.SystemUid,
			"spareItemUid": request.SpareItemUid,
		},
	}
	queries = append(queries, assignSpareQuery)

	// Update spare item usage to the current item's usage (if current item exists)
	if currentItemUid != "" && currentItemUsageCode != "" {
		updateSpareItemUsageQuery := helpers.DatabaseQuery{
			Query: `MATCH (spareItem:Item {uid: $spareItemUid})
                    OPTIONAL MATCH (spareItem)-[oldUsageRel:HAS_ITEM_USAGE]->(:ItemUsage)
                    DELETE oldUsageRel
                    WITH spareItem
                    MATCH (newUsage:ItemUsage {code: $usageCode})
                    CREATE (spareItem)-[:HAS_ITEM_USAGE]->(newUsage)`,
			Parameters: map[string]interface{}{
				"spareItemUid": request.SpareItemUid,
				"usageCode":    currentItemUsageCode,
			},
		}
		queries = append(queries, updateSpareItemUsageQuery)
	}

	// Remove ALL IS_SPARE_FOR relationships from the spare system (since it's no longer a spare for any system)
	removeAllSpareRelQuery := helpers.DatabaseQuery{
		Query: `MATCH (spareSystem:System {uid: $spareSystemUid})-[r:IS_SPARE_FOR]->(:System) DELETE r`,
		Parameters: map[string]interface{}{
			"spareSystemUid": spareSystemUid,
		},
	}
	queries = append(queries, removeAllSpareRelQuery)

	// If current item exists, relocate it to the spare system
	if currentItemUid != "" {
		response.RelocatedItemUid = currentItemUid

		// Detach current item from target system
		detachQuery := helpers.DatabaseQuery{
			Query: `MATCH (s:System {uid: $systemUid})-[r:CONTAINS_ITEM]->(oldItem:Item {uid: $oldItemUid}) DELETE r`,
			Parameters: map[string]interface{}{
				"systemUid":  request.SystemUid,
				"oldItemUid": currentItemUid,
			},
		}
		queries = append(queries, detachQuery)

		// Update old item condition
		updateItemConditionQuery := helpers.DatabaseQuery{
			Query: `MATCH (item:Item {uid: $oldItemUid})
                    OPTIONAL MATCH (item)-[oldCondRel:HAS_CONDITION_STATUS]->(:ItemCondition)
                    DELETE oldCondRel
                    WITH item
                    MATCH (newCondition:ItemCondition {uid: $conditionUid})
                    CREATE (item)-[:HAS_CONDITION_STATUS]->(newCondition)
                    SET item.updatedAt = datetime()`,
			Parameters: map[string]interface{}{
				"oldItemUid":   currentItemUid,
				"conditionUid": request.OldItemCondition.UID,
			},
		}
		queries = append(queries, updateItemConditionQuery)

		// Update old item usage to "stock-item"
		updateOldItemUsageQuery := helpers.DatabaseQuery{
			Query: `MATCH (item:Item {uid: $oldItemUid})
                    OPTIONAL MATCH (item)-[oldUsageRel:HAS_ITEM_USAGE]->(:ItemUsage)
                    DELETE oldUsageRel
                    WITH item
                    MATCH (stockUsage:ItemUsage {code: 'stock-item'})
                    CREATE (item)-[:HAS_ITEM_USAGE]->(stockUsage)`,
			Parameters: map[string]interface{}{
				"oldItemUid": currentItemUid,
			},
		}
		queries = append(queries, updateOldItemUsageQuery)

		// Assign old item to spare system (reuse the spare system)
		assignOldItemToSpareSystemQuery := helpers.DatabaseQuery{
			Query: `MATCH (spareSystem:System {uid: $spareSystemUid})
                    MATCH (oldItem:Item {uid: $oldItemUid})
                    CREATE (spareSystem)-[:CONTAINS_ITEM]->(oldItem)`,
			Parameters: map[string]interface{}{
				"spareSystemUid": spareSystemUid,
				"oldItemUid":     currentItemUid,
			},
		}
		queries = append(queries, assignOldItemToSpareSystemQuery)

		// Update spare system status to active (it now contains the old item)
		updateSpareSystemStatusQuery := helpers.DatabaseQuery{
			Query: `MATCH (spareSystem:System {uid: $spareSystemUid})
                    SET spareSystem.status = 'active', spareSystem.updatedAt = datetime()`,
			Parameters: map[string]interface{}{
				"spareSystemUid": spareSystemUid,
			},
		}
		queries = append(queries, updateSpareSystemStatusQuery)

		// Set location for the spare system if provided
		if request.NewItemLocation.UID != "" {
			setLocationQuery := helpers.DatabaseQuery{
				Query: `MATCH (spareSystem:System {uid: $spareSystemUid})
                        OPTIONAL MATCH (spareSystem)-[oldLocRel:HAS_LOCATION]->(:Location)
                        DELETE oldLocRel
                        WITH spareSystem
                        MATCH (location:Location {uid: $locationUid})
                        CREATE (spareSystem)-[:HAS_LOCATION]->(location)`,
				Parameters: map[string]interface{}{
					"spareSystemUid": spareSystemUid,
					"locationUid":    request.NewItemLocation.UID,
				},
			}
			queries = append(queries, setLocationQuery)
		}

		// Determine parent system for spare system placement
		newParentSystemUid := request.NewParentSystemUid
		if newParentSystemUid == "" {
			// Use hardcoded unassigned items system as requested
			newParentSystemUid = "c6a8e743-f348-4e4f-b5cc-c38da8bf4641"
		}

		response.NewParentSystemUid = newParentSystemUid

		// Move spare system to parent system (remove old parent relationships first)
		removeOldParentQuery := helpers.DatabaseQuery{
			Query: `MATCH (parentSys:System)-[r:HAS_SUBSYSTEM]->(spareSystem:System {uid: $spareSystemUid}) DELETE r`,
			Parameters: map[string]interface{}{
				"spareSystemUid": spareSystemUid,
			},
		}
		queries = append(queries, removeOldParentQuery)

		// Add new parent relationship
		addNewParentQuery := helpers.DatabaseQuery{
			Query: `MATCH (parentSys:System {uid: $parentSystemUid})
                    MATCH (spareSystem:System {uid: $spareSystemUid})
                    CREATE (parentSys)-[:HAS_SUBSYSTEM]->(spareSystem)`,
			Parameters: map[string]interface{}{
				"parentSystemUid": newParentSystemUid,
				"spareSystemUid":  spareSystemUid,
			},
		}
		queries = append(queries, addNewParentQuery)
	} else {
		// No current item exists, mark spare system as deleted since it's now empty
		markSpareSystemDeletedQuery := helpers.DatabaseQuery{
			Query: `MATCH (spareSystem:System {uid: $spareSystemUid})
                    SET spareSystem.status = 'deleted', spareSystem.updatedAt = datetime()`,
			Parameters: map[string]interface{}{
				"spareSystemUid": spareSystemUid,
			},
		}
		queries = append(queries, markSpareSystemDeletedQuery)
	}

	// Add history log
	historyLog := helpers.HistoryLogQuery(request.SystemUid, "ASSIGN_SPARE", userUID)
	queries = append(queries, historyLog)

	// Execute all queries in transaction
	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session, queries...)
	if err != nil {
		return response, err
	}

	// Recalculate spare parts coverage since IS_SPARE_FOR relationships have changed
	err = svc.RecalculateSpareParts()
	if err != nil {
		log.Error().Err(err).Msg("Failed to recalculate spare parts coverage after assignment")
		// Don't fail the entire operation if recalculation fails, just log the error
	}

	response.Success = true
	response.Message = "Spare item assigned successfully"
	return response, nil
}

// GetSystemSparePartsDetail returns comprehensive system and physical item information with all spare relations
func (s *SystemsService) GetSystemSparePartsDetail(systemId string, facilityCode string) (result models.SystemSparePartsDetail, err error) {
	driver := *s.neo4jDriver
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	// Execute multiple queries to avoid nested aggregations
	// Query 1: Get basic system information
	systemQuery := `
		MATCH (s:System {uid: $systemId})-[:BELONGS_TO_FACILITY]->(facility:Facility {code: $facilityCode})
		OPTIONAL MATCH (s)-[:HAS_ZONE]->(zone:Zone)
		OPTIONAL MATCH (s)-[:HAS_LOCATION]->(location:Location)
		OPTIONAL MATCH (s)-[:HAS_SYSTEM_TYPE]->(systemType:SystemType)
		OPTIONAL MATCH (s)-[:HAS_RESPONSIBLE]->(responsible:Employee)
		OPTIONAL MATCH (s)-[:HAS_OPERATOR]->(operator:Employee)
		OPTIONAL MATCH (s)-[:IS_MAINTAINED_BY]->(maintainer:Employee)
		OPTIONAL MATCH (s)-[:HAS_RESPONSIBLE_TEAM]->(team:Team)
		
		RETURN {
			system: {
				uid: s.uid,
				name: s.name,
				systemCode: s.systemCode,
				description: s.description,
				status: s.status,
				isTechnologicalUnit: s.isTechnologicalUnit,
				isCritical: s.isCritical,
				minimalSpareParstCount: s.minimalSpareParstCount,
				sparePartsCoverageSum: s.sparePartsCoverageSum,
				systemLevel: s.systemLevel,
				systemAlias: s.systemAlias,
				image: s.image,
				miniImageUrl: s.miniImageUrl,
				lastUpdateTime: toString(s.lastUpdateTime),
				lastUpdateBy: s.lastUpdateBy
			},
			location: CASE WHEN location IS NOT NULL THEN {
				uid: location.uid,
				name: location.name,
				code: location.code
			} ELSE null END,
			zone: CASE WHEN zone IS NOT NULL THEN {
				uid: zone.uid,
				name: zone.name,
				code: zone.code
			} ELSE null END,
			systemType: CASE WHEN systemType IS NOT NULL THEN {
				uid: systemType.uid,
				name: systemType.name,
				code: systemType.code
			} ELSE null END,
			responsiblePersons: {
				responsible: CASE WHEN responsible IS NOT NULL THEN {
					uid: responsible.uid,
					firstName: responsible.firstName,
					lastName: responsible.lastName,
					email: responsible.email,
					phone: responsible.phone
				} ELSE null END,
				operator: CASE WHEN operator IS NOT NULL THEN {
					uid: operator.uid,
					firstName: operator.firstName,
					lastName: operator.lastName,
					email: operator.email,
					phone: operator.phone
				} ELSE null END,
				maintainer: CASE WHEN maintainer IS NOT NULL THEN {
					uid: maintainer.uid,
					firstName: maintainer.firstName,
					lastName: maintainer.lastName,
					email: maintainer.email,
					phone: maintainer.phone
				} ELSE null END
			},
			team: CASE WHEN team IS NOT NULL THEN {
				uid: team.uid,
				name: team.name
			} ELSE null END
		} AS result
	`

	systemResult, err := session.Run(systemQuery, map[string]interface{}{
		"systemId":     systemId,
		"facilityCode": facilityCode,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error executing system basic info query")
		return result, err
	}

	var basicInfo map[string]interface{}
	if systemResult.Next() {
		record := systemResult.Record()
		if resultValue, ok := record.Get("result"); ok {
			basicInfo = resultValue.(map[string]interface{})
		}
	}

	// Query 2: Get system attributes
	attributesQuery := `
		MATCH (s:System {uid: $systemId})-[:HAS_SYSTEM_ATTRIBUTE]->(attr:SystemAttribute)
		RETURN COLLECT({
			uid: attr.uid,
			name: attr.name
		}) AS systemAttributes
	`

	attrResult, err := session.Run(attributesQuery, map[string]interface{}{
		"systemId": systemId,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error executing system attributes query")
		return result, err
	}

	var systemAttributes []interface{}
	if attrResult.Next() {
		record := attrResult.Record()
		if attrs, ok := record.Get("systemAttributes"); ok {
			systemAttributes = attrs.([]interface{})
		}
	}

	// Query 3: Get physical items
	itemsQuery := `
		MATCH (s:System {uid: $systemId})-[:CONTAINS_ITEM]->(item:Item)
		OPTIONAL MATCH (item)-[:IS_BASED_ON]->(catalogueItem:CatalogueItem)
		OPTIONAL MATCH (catalogueItem)-[:BELONGS_TO_CATEGORY]->(category:CatalogueCategory)
		OPTIONAL MATCH (catalogueItem)-[:HAS_MANUFACTURER]->(manufacturer:Manufacturer)
		OPTIONAL MATCH (item)-[:HAS_CONDITION_STATUS]->(condition:ItemCondition)
		OPTIONAL MATCH (item)-[:HAS_ITEM_USAGE]->(usage:ItemUsage)
		
		RETURN COLLECT({
			uid: item.uid,
			name: item.name,
			serialNumber: item.serialNumber,
			eun: item.eun,
			price: item.price,
			currency: item.currency,
			status: item.status,
			notes: item.notes,
			printEUN: item.printEUN,
			lastUpdateTime: toString(item.lastUpdateTime),
			condition: CASE WHEN condition IS NOT NULL THEN {
				uid: condition.uid,
				name: condition.name,
				code: condition.code
			} ELSE null END,
			usage: CASE WHEN usage IS NOT NULL THEN {
				uid: usage.uid,
				name: usage.name,
				code: usage.code
			} ELSE null END,
			catalogueItem: CASE WHEN catalogueItem IS NOT NULL THEN {
				uid: catalogueItem.uid,
				name: catalogueItem.name,
				description: catalogueItem.description,
				catalogueNumber: catalogueItem.catalogueNumber,
				image: catalogueItem.image,
				miniImageUrl: catalogueItem.miniImageUrl,
				manufacturerUrl: catalogueItem.manufacturerUrl,
				lastUpdateTime: toString(catalogueItem.lastUpdateTime),
				category: CASE WHEN category IS NOT NULL THEN {
					uid: category.uid,
					name: category.name,
					code: category.code,
					image: category.image,
					miniImageUrl: category.miniImageUrl
				} ELSE null END,
				manufacturer: CASE WHEN manufacturer IS NOT NULL THEN {
					uid: manufacturer.uid,
					name: manufacturer.name
				} ELSE null END
			} ELSE null END
		}) AS physicalItems
	`

	itemsResult, err := session.Run(itemsQuery, map[string]interface{}{
		"systemId": systemId,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error executing physical items query")
		return result, err
	}

	var physicalItems []interface{}
	if itemsResult.Next() {
		record := itemsResult.Record()
		if items, ok := record.Get("physicalItems"); ok {
			physicalItems = items.([]interface{})
		}
	}

	// Query 4: Get spare systems with their physical items
	spareSystemsQuery := `
		MATCH (spareSystem:System)-[:IS_SPARE_FOR]->(s:System {uid: $systemId})
		OPTIONAL MATCH (spareSystem)-[:CONTAINS_ITEM]->(spareItem:Item)
		OPTIONAL MATCH (spareItem)-[:IS_BASED_ON]->(spareCatalogueItem:CatalogueItem)
		
		WITH spareSystem, COLLECT({
			uid: spareItem.uid,
			name: spareItem.name,
			serialNumber: spareItem.serialNumber,
			eun: spareItem.eun,
			catalogueItem: CASE WHEN spareCatalogueItem IS NOT NULL THEN {
				uid: spareCatalogueItem.uid,
				name: spareCatalogueItem.name,
				catalogueNumber: spareCatalogueItem.catalogueNumber
			} ELSE null END
		}) AS spareItems
		
		RETURN COLLECT({
			uid: spareSystem.uid,
			name: spareSystem.name,
			systemCode: spareSystem.systemCode,
			description: spareSystem.description,
			status: spareSystem.status,
			physicalItems: [item IN spareItems WHERE item.uid IS NOT NULL | item]
		}) AS spareSystems
	`

	spareSystemsResult, err := session.Run(spareSystemsQuery, map[string]interface{}{
		"systemId": systemId,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error executing spare systems query")
		return result, err
	}

	var spareSystems []interface{}
	if spareSystemsResult.Next() {
		record := spareSystemsResult.Record()
		if systems, ok := record.Get("spareSystems"); ok {
			spareSystems = systems.([]interface{})
		}
	}

	// Query 5: Get parent systems (systems this system is spare for)
	parentSystemsQuery := `
		MATCH (s:System {uid: $systemId})-[:IS_SPARE_FOR]->(parentSystem:System)
		RETURN COLLECT({
			uid: parentSystem.uid,
			name: parentSystem.name,
			systemCode: parentSystem.systemCode,
			description: parentSystem.description,
			status: parentSystem.status
		}) AS parentSystems
	`

	parentSystemsResult, err := session.Run(parentSystemsQuery, map[string]interface{}{
		"systemId": systemId,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error executing parent systems query")
		return result, err
	}

	var parentSystems []interface{}
	if parentSystemsResult.Next() {
		record := parentSystemsResult.Record()
		if systems, ok := record.Get("parentSystems"); ok {
			parentSystems = systems.([]interface{})
		}
	}

	// Combine all results
	if basicInfo != nil {
		basicInfo["systemAttributes"] = systemAttributes
		basicInfo["physicalItems"] = physicalItems
		basicInfo["sparePartsRelations"] = map[string]interface{}{
			"spareSystems":  spareSystems,
			"parentSystems": parentSystems,
		}

		// Use MapStruct to convert the result map to our model
		mappedResult, err := helpers.MapStruct[models.SystemSparePartsDetail](basicInfo)
		if err != nil {
			log.Error().Err(err).Msg("Error mapping system spare parts detail result")
			return result, err
		}

		result = mappedResult
		log.Info().Str("systemId", systemId).Msg("Successfully retrieved system spare parts detail")
	}

	return result, nil
}
