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
	GetOrdersWithSearchAndPagination(search string, facilityCode string, pagination *helpers.Pagination, sorting *[]helpers.Sorting) (result helpers.PaginationResult[models.OrderListItem], err error)
	GetOrderWithOrderLinesByUid(orderUid string) (result models.OrderDetail, err error)
	InsertNewOrder(order *models.OrderDetail, facilityCode string, userUID string) (uid string, err error)
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

func (svc *OrdersService) GetOrdersWithSearchAndPagination(search string, facilityCode string, pagination *helpers.Pagination, sorting *[]helpers.Sorting) (result helpers.PaginationResult[models.OrderListItem], err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	//beacause of the full text search we need to modify the search string
	search = helpers.GetFullTextSearchString(search)

	query := GetOrdersBySearchTextFullTextQuery(search, facilityCode, pagination, sorting)
	items, err := helpers.GetNeo4jArrayOfNodes[models.OrderListItem](session, query)
	totalCount, _ := helpers.GetNeo4jSingleRecordSingleValue[int64](session, GetOrdersBySearchTextFullTextCountQuery(search, facilityCode))

	result = helpers.GetPaginationResult(items, int64(totalCount), err)

	return result, err
}

func (svc *OrdersService) GetOrderWithOrderLinesByUid(orderUid string) (result models.OrderDetail, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetOrderWithOrderLinesByUidQuery(orderUid)
	result, err = helpers.GetNeo4jSingleRecordAndMapToStruct[models.OrderDetail](session, query)

	return result, err
}

func (svc *OrdersService) InsertNewOrder(order *models.OrderDetail, facilityCode string, userUID string) (uid string, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := InsertNewOrderQuery(order, facilityCode, userUID)
	uid, err = helpers.WriteNeo4jAndReturnSingleValue[string](session, query)

	return uid, err
}
