package systemsService

import (
	"panda/apigateway/config"
	"panda/apigateway/helpers"
	codebookModels "panda/apigateway/services/codebook-service/models"
	"panda/apigateway/services/systems-service/models"
	systemsModels "panda/apigateway/services/systems-service/models"

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
