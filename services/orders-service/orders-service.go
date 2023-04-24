package ordersService

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"

	"panda/apigateway/helpers"
	codebookModels "panda/apigateway/services/codebook-service/models"
)

type OrdersService struct {
	neo4jDriver *neo4j.Driver
}

type IOrdersService interface {
	GetOrderStatusesCodebook() (result []codebookModels.Codebook, err error)
	GetSuppliersAutocompleteCodebook(searchString string, limit int) (result []codebookModels.Codebook, err error)
}

func NewOrdersService(driver *neo4j.Driver) IOrdersService {
	return &OrdersService{neo4jDriver: driver}
}

func (svc *OrdersService) GetOrderStatusesCodebook() (result []codebookModels.Codebook, err error) {
	// Open a new Session
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	//get all categories by parent path
	query := GetOrderStatusesCodebookQuery()
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}

func (svc *OrdersService) GetSuppliersAutocompleteCodebook(searchString string, limit int) (result []codebookModels.Codebook, err error) {
	// Open a new Session
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	//get all categories by parent path
	query := GetSuppliersAutoCompleteQuery(searchString, limit)
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}
