package systemsService

import (
	"errors"
	"panda/apigateway/config"
	"panda/apigateway/helpers"
	codebookModels "panda/apigateway/services/codebook-service/models"
	"panda/apigateway/services/systems-service/models"
	systemsModels "panda/apigateway/services/systems-service/models"
	"strconv"
	"strings"

	//"strings"

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
	GetZonesCodebook(facilityCode string) (result []codebookModels.Codebook, err error)
	GetSubSystemsByParentUID(parentUID string, facilityCode string) (result []systemsModels.System, err error)
	GetSystemImageByUid(uid string) (imageBase64 string, err error)
	GetSystemDetail(uid string, facilityCode string) (result models.System, err error)
	CreateNewSystem(system *models.System, facilityCode string, userUID string) (uid string, err error)
	UpdateSystem(newSystem *models.System, facilityCode string, userUID string) (err error)
	DeleteSystemRecursive(uid string) (err error)
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

	if searchText != "" {
		query := GetLocationsBySearchTextQuery(searchText, limit, facilityCode)
		result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)
	} else {
		result = make([]codebookModels.Codebook, 0)
	}

	return result, err
}

func (svc *SystemsService) GetZonesCodebook(facilityCode string) (result []codebookModels.Codebook, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetZonesCodebookQuery(facilityCode)
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

func (svc *SystemsService) DeleteSystemRecursive(uid string) (err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := DeleteSystemByUidQuery(uid)
	err = helpers.WriteNeo4jAndReturnNothing(session, query)

	return err
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

	if strings.Contains(mask, SYSTEM_CODE_GENERATE_ZONE_CODE) || strings.Contains(mask, SYSTEM_CODE_GENERATE_ZONE_NAME) {
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

			mask = strings.ReplaceAll(mask, SYSTEM_CODE_GENERATE_ZONE_CODE, zoneCode)
			mask = strings.ReplaceAll(mask, SYSTEM_CODE_GENERATE_ZONE_NAME, zoneName)
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
