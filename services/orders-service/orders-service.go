package ordersService

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"

	"panda/apigateway/helpers"
	codebookModels "panda/apigateway/services/codebook-service/models"
	"panda/apigateway/services/orders-service/models"
)

type OrdersService struct {
	neo4jDriver *neo4j.Driver
}

type IOrdersService interface {
	GetOrderStatusesCodebook() (result []codebookModels.Codebook, err error)
	GetSuppliersAutocompleteCodebook(searchString string, limit int) (result []codebookModels.Codebook, err error)
	GetOrdersWithSearchAndPagination(search string, pagination *helpers.Pagination, sorting *[]helpers.Sorting) (result helpers.PaginationResult[models.OrderListItem], err error)
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

func (svc *OrdersService) GetOrdersWithSearchAndPagination(search string, pagination *helpers.Pagination, sorting *[]helpers.Sorting) (result helpers.PaginationResult[models.OrderListItem], err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetOrdersBySearchTextFullTextQuery(search, pagination, sorting)
	items, err := helpers.GetNeo4jArrayOfNodes[models.OrderListItem](session, query)
	totalCount, _ := helpers.GetNeo4jSingleRecordSingleValue[int64](session, GetOrdersBySearchTextFullTextCountQuery(search))

	result = helpers.GetPaginationResult(items, int64(totalCount), err)

	return result, err
}
