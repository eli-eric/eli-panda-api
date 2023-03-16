package systemsService

import (
	"log"
	"panda/apigateway/config"
	"panda/apigateway/helpers"
	codebookModels "panda/apigateway/services/codebook-service/models"
	"panda/apigateway/services/systems-service/models"
	systemsModels "panda/apigateway/services/systems-service/models"

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
	GetSubSystemsByParentUID(parentUID string, facilityCode string) (result []systemsModels.SystemSimpleResponse, err error)
	GetSystemImageByUid(uid string) (imageBase64 string, err error)
	GetSystemDetail(uid string, facilityCode string) (result models.SystemResponse, err error)
	CreateNewSystem(system *models.SystemForm, facilityCode string, userUID string) (uid string, err error)
	UpdateSystem(newSystem *models.SystemForm, facilityCode string, userUID string) (err error)
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

func (svc *SystemsService) GetSubSystemsByParentUID(parentUID string, facilityCode string) (result []models.SystemSimpleResponse, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetSubSystemsQuery(parentUID, facilityCode)
	result, err = helpers.GetNeo4jArrayOfNodes[models.SystemSimpleResponse](session, query)

	return result, err
}

func (svc *SystemsService) GetSystemImageByUid(uid string) (result string, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := SystemImageByUidQuery(uid)
	result, err = helpers.GetNeo4jSingleRecordSingleValue[string](session, query)

	return result, err
}

func (svc *SystemsService) GetSystemDetail(uid string, facilityCode string) (result models.SystemResponse, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	result, err = helpers.GetNeo4jSingleRecordAndMapToStruct[models.SystemResponse](session, SystemDetailQuery(uid, facilityCode))

	return result, err
}

func (svc *SystemsService) CreateNewSystem(system *models.SystemForm, facilityCode string, userUID string) (uid string, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	uid, err = helpers.WriteNeo4jAndReturnSingleValue[string](session, CreateNewSystemQuery(system, facilityCode))

	if err != nil {
		log.Println(err.Error())
	} else {
		go func() {
			// we dont want to log image data
			system.Image = nil
			helpers.LogDBHistory(session, uid, nil, system, userUID, helpers.DB_LOG_CREATE)
		}()
	}

	return uid, err
}

func (svc *SystemsService) UpdateSystem(system *models.SystemForm, facilityCode string, userUID string) (err error) {

	if system != nil {
		session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

		oldSystem, err := helpers.GetNeo4jSingleRecordAndMapToStruct[models.SystemForm](session, SystemFormDetailQuery(*system.UID, facilityCode))

		if err == nil {
			_, err = helpers.WriteNeo4jAndReturnSingleValue[string](session, UpdateSystemQuery(system, &oldSystem, facilityCode))
		}

		if err != nil {
			log.Println(err.Error())
		} else {
			go func() {
				// we dont want to log image data
				system.Image = nil
				helpers.LogDBHistory(session, *oldSystem.UID, oldSystem, system, userUID, helpers.DB_LOG_UPDATE)
			}()
		}

	} else {
		err = helpers.ERR_INVALID_INPUT
	}
	return err
}
