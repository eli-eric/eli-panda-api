package systemsService

import (
	"errors"
	"fmt"
	"panda/apigateway/config"
	"panda/apigateway/helpers"
	codebookModels "panda/apigateway/services/codebook-service/models"
	"panda/apigateway/services/systems-service/models"
	systemsModels "panda/apigateway/services/systems-service/models"
	"strconv"
	"strings"

	//"strings"

	"github.com/google/uuid"
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
	GetZonesCodebook(facilityCode string, searchString string) (result []codebookModels.Codebook, err error)
	GetSubSystemsByParentUID(parentUID string, facilityCode string) (result []systemsModels.System, err error)
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

func (svc *SystemsService) GetZonesCodebook(facilityCode string, searchString string) (result []codebookModels.Codebook, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetZonesCodebookQuery(facilityCode, searchString)
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

func ValidatePhysicalItemMovement(movement *systemsModels.PhysicalItemMovement) error {

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

func ValidatePhysicalItemReplacement(movement *systemsModels.PhysicalItemMovement) error {

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
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
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

	// Step 2: Validate spare item relationship
	validateSpareQuery := `
        MATCH (spareItem:Item {uid: $spareItemUid})<-[:CONTAINS_ITEM]-(spareSystem:System)-[:IS_SPARE_FOR]->(targetSystem:System {uid: $systemUid})
        RETURN spareItem.uid as spareUid
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

	// Step 3: Execute the assignment in a transaction
	queries := []helpers.DatabaseQuery{}

	// If current item exists, prepare relocation
	if currentItemUid != "" {
		response.RelocatedItemUid = currentItemUid

		// Detach current item from system
		detachQuery := helpers.DatabaseQuery{
			Query: `MATCH (s:System {uid: $systemUid})-[r:CONTAINS_ITEM]->(oldItem:Item {uid: $oldItemUid}) DELETE r`,
			Parameters: map[string]interface{}{
				"systemUid":  request.SystemUid,
				"oldItemUid": currentItemUid,
			},
		}
		queries = append(queries, detachQuery)

		// Update old item condition using relationship
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

		// Create new system for old item with location
		newSystemUid := uuid.New().String()
		createNewSystemQuery := helpers.DatabaseQuery{
			Query: `CREATE (newSys:System {
                            uid: $newSystemUid,
                            name: 'Relocated System',
                            type: 'relocated',
                            status: 'active',
                            createdAt: datetime(),
                            updatedAt: datetime()
                        })
                        WITH newSys
                        MATCH (location:Location {uid: $locationUid})
                        CREATE (newSys)-[:HAS_LOCATION]->(location)`,
			Parameters: map[string]interface{}{
				"newSystemUid": newSystemUid,
				"locationUid":  request.NewItemLocation.UID,
			},
		}
		queries = append(queries, createNewSystemQuery)
	}

	// Assign spare item to system
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

	// Remove IS_SPARE_FOR relationship
	removeSpareRelQuery := helpers.DatabaseQuery{
		Query: `MATCH (spareItem:Item {uid: $spareItemUid})-[:CONTAINS]-(spareSystem:System)-[r:IS_SPARE_FOR]->(s:System {uid: $systemUid}) DELETE r`,
		Parameters: map[string]interface{}{
			"spareItemUid": request.SpareItemUid,
			"systemUid":    request.SystemUid,
		},
	}
	queries = append(queries, removeSpareRelQuery)

	// Mark spare item's system as deleted
	markSpareSystemDeletedQuery := helpers.DatabaseQuery{
		Query: `MATCH (spareItem:Item {uid: $spareItemUid})-[:CONTAINS]-(spareSystem:System)
                SET spareSystem.status = 'deleted', spareSystem.updatedAt = datetime()`,
		Parameters: map[string]interface{}{
			"spareItemUid": request.SpareItemUid,
		},
	}
	queries = append(queries, markSpareSystemDeletedQuery)

	// If current item exists, handle relocation
	if currentItemUid != "" {
		newParentSystemUid := request.NewParentSystemUid

		// Auto-detect parent system if not provided
		if newParentSystemUid == "" {
			autoDetectQuery := helpers.DatabaseQuery{
				Query: `MATCH (currentSystem:System {uid: $systemUid})
                        MATCH (currentSystem)-[:HAS_SUBSYSTEM*]->(parentSystem:System)
                        WHERE parentSystem.type = 'trash' OR parentSystem.name CONTAINS 'trash'
                        RETURN parentSystem.uid as parentUid
                        ORDER BY length((currentSystem)-[:HAS_SUBSYSTEM*]->(parentSystem)) ASC
                        LIMIT 1`,
				Parameters: map[string]interface{}{
					"systemUid": request.SystemUid,
				},
			}

			result, err := session.Run(autoDetectQuery.Query, autoDetectQuery.Parameters)
			if err == nil && result.Next() {
				newParentSystemUid = result.Record().Values[0].(string)
			}
		}

		if newParentSystemUid != "" {
			response.NewParentSystemUid = newParentSystemUid

			// Create new system for old item
			newSystemUid := uuid.New().String()
			createNewSystemQuery := helpers.DatabaseQuery{
				Query: `CREATE (newSys:System {
                            uid: $newSystemUid,
                            name: 'Relocated System',
                            type: 'relocated',
                            status: 'active',
                            createdAt: datetime(),
                            updatedAt: datetime()
                        })`,
				Parameters: map[string]interface{}{
					"newSystemUid": newSystemUid,
				},
			}
			queries = append(queries, createNewSystemQuery)

			// Link old item to new system
			linkItemToNewSystemQuery := helpers.DatabaseQuery{
				Query: `MATCH (newSys:System {uid: $newSystemUid})
                        MATCH (oldItem:Item {uid: $oldItemUid})
                        CREATE (newSys)-[:CONTAINS]->(oldItem)`,
				Parameters: map[string]interface{}{
					"newSystemUid": newSystemUid,
					"oldItemUid":   currentItemUid,
				},
			}
			queries = append(queries, linkItemToNewSystemQuery)

			// Move new system to parent system
			moveSystemQuery := helpers.DatabaseQuery{
				Query: `MATCH (parentSys:System {uid: $parentSystemUid})
                        MATCH (newSys:System {uid: $newSystemUid})
                        CREATE (parentSys)-[:HAS_SUBSYSTEM]->(newSys)`,
				Parameters: map[string]interface{}{
					"parentSystemUid": newParentSystemUid,
					"newSystemUid":    newSystemUid,
				},
			}
			queries = append(queries, moveSystemQuery)
		}
	}

	// Add history log
	historyLog := helpers.HistoryLogQuery(request.SystemUid, "ASSIGN_SPARE", userUID)
	queries = append(queries, historyLog)

	// Execute all queries
	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session, queries...)
	if err != nil {
		return response, err
	}

	response.Success = true
	response.Message = "Spare item assigned successfully"
	return response, nil
}
