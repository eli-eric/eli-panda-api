package systemsService

import (
	"panda/apigateway/config"
	"panda/apigateway/helpers"
	codebookModels "panda/apigateway/services/codebook-service/models"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type SystemsService struct {
	neo4jDriver *neo4j.Driver
	jwtSecret   string
}

type ISystemsService interface {
	GetSystemTypesCodebook() (result []codebookModels.Codebook, err error)
}

// Create new security service instance
func NewSystemsService(settings *config.Config, driver *neo4j.Driver) ISystemsService {

	return &SystemsService{neo4jDriver: driver, jwtSecret: settings.JwtSecret}
}


func (svc *SystemsService) GetSystemTypesCodebook() (result []codebookModels.Codebook, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetSystemTypesCodebookQuery()
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	return result, err
}
