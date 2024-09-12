package general

import (
	"panda/apigateway/helpers"
	"panda/apigateway/services/general/models"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type GeneralService struct {
	neo4jDriver *neo4j.Driver
}

type IGeneralService interface {
	GetGraphNodesByUid(uid string) (result []models.GraphNode, err error)
	GetGraphLinksByUid(uid string) (result []models.GraphLink, err error)
}

func NewGeneralService(driver *neo4j.Driver) IGeneralService {
	return &GeneralService{neo4jDriver: driver}
}

func (svc *GeneralService) GetGraphNodesByUid(uid string) (result []models.GraphNode, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetGraphNodesByUidQuery(uid)
	result, err = helpers.GetNeo4jArrayOfNodes[models.GraphNode](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}

func (svc *GeneralService) GetGraphLinksByUid(uid string) (result []models.GraphLink, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetGraphLinksByUidQuery(uid)
	result, err = helpers.GetNeo4jArrayOfNodes[models.GraphLink](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}
