package ordersService

import (
	"github.com/rs/zerolog/log"

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
	GetOrdersWithSearchAndPagination(search string, supplierUID string, orderStatusUID string, procurementResponsibleUID string, requestorUID string, facilityCode string, pagination *helpers.Pagination, sorting *[]helpers.Sorting) (result helpers.PaginationResult[models.OrderListItem], err error)
	GetOrderWithOrderLinesByUid(orderUid string, facilityCode string) (result models.OrderDetail, err error)
	InsertNewOrder(order *models.OrderDetail, facilityCode string, userUID string) (uid string, err error)
	UpdateOrder(order *models.OrderDetail, facilityCode string, userUID string) (err error)
	DeleteOrder(orderUid string, userUID string) (err error)
	UpdateOrderLineDelivery(itemUid string, isDelivered bool, serialNumber *string, eun *string, userUID string, facilityCode string) (result models.OrderLine, err error)
	GetItemsForEunPrint(euns []string) (result []models.ItemForEunPrint, err error)
	SetItemPrintEUN(eun string, printEUN bool) (err error)
	GetOrderUidByOrderNumber(orderNumber string) (result string, err error)
	GetOrdersForCatalogueItem(catalogueItemUid string, facilityCode string) (result []models.OrderListItem, err error)
	GetMinAndMaxOrderLinePrice(facilityCode string) (result models.OrderLineMinMax, err error)
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

func (svc *OrdersService) GetOrdersWithSearchAndPagination(search string, supplierUID string, orderStatusUID string, procurementResponsibleUID string, requestorUID string, facilityCode string, pagination *helpers.Pagination, sorting *[]helpers.Sorting) (result helpers.PaginationResult[models.OrderListItem], err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetOrdersBySearchTextFullTextQuery(search, supplierUID, orderStatusUID, procurementResponsibleUID, requestorUID, facilityCode, pagination, sorting)
	items, err := helpers.GetNeo4jArrayOfNodes[models.OrderListItem](session, query)
	totalCount, _ := helpers.GetNeo4jSingleRecordSingleValue[int64](session, GetOrdersBySearchTextFullTextCountQuery(search, supplierUID, orderStatusUID, procurementResponsibleUID, requestorUID, facilityCode))

	result = helpers.GetPaginationResult(items, int64(totalCount), err)

	return result, err
}

func (svc *OrdersService) GetOrderWithOrderLinesByUid(orderUid string, facilityCode string) (result models.OrderDetail, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetOrderWithOrderLinesByUidQuery(orderUid, facilityCode)
	result, err = helpers.GetNeo4jSingleRecordAndMapToStruct[models.OrderDetail](session, query)

	return result, err
}

func (svc *OrdersService) InsertNewOrder(order *models.OrderDetail, facilityCode string, userUID string) (uid string, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := InsertNewOrderQuery(order, facilityCode, userUID)
	uid, err = helpers.WriteNeo4jAndReturnSingleValue[string](session, query)

	return uid, err
}

func (svc *OrdersService) DeleteOrder(orderUid string, userUID string) (err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := DeleteOrderQuery(orderUid, userUID)
	err = helpers.WriteNeo4jAndReturnNothing(session, query)

	return err
}

func (svc *OrdersService) UpdateOrder(order *models.OrderDetail, facilityCode string, userUID string) (err error) {
	if order != nil {
		session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

		oldOrder, err := helpers.GetNeo4jSingleRecordAndMapToStruct[models.OrderDetail](session, GetOrderWithOrderLinesByUidQuery(order.UID, facilityCode))

		if err == nil {
			_, err = helpers.WriteNeo4jAndReturnSingleValue[string](session, UpdateOrderQuery(order, &oldOrder, facilityCode, userUID))
		}

		if err != nil {
			log.Error().Msg(err.Error())
			return err
		}

	} else {
		err = helpers.ERR_INVALID_INPUT
	}

	return err
}

func (svc *OrdersService) UpdateOrderLineDelivery(itemUid string, isDelivered bool, serialNumber *string, eun *string, userUID string, facilityCode string) (result models.OrderLine, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := UpdateOrderLineDeliveryQuery(itemUid, isDelivered, serialNumber, eun, userUID, facilityCode)
	result, err = helpers.WriteNeo4jReturnSingleRecordAndMapToStruct[models.OrderLine](session, query)

	return result, err
}

func (svc *OrdersService) GetItemsForEunPrint(euns []string) (result []models.ItemForEunPrint, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetItemsForEunPrintQuery(euns)
	result, err = helpers.GetNeo4jArrayOfNodes[models.ItemForEunPrint](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}

func (svc *OrdersService) SetItemPrintEUN(eun string, printEUN bool) (err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := SetItemPrintEUNQuery(eun, printEUN)
	err = helpers.WriteNeo4jAndReturnNothing(session, query)

	return err
}

func (svc *OrdersService) GetOrderUidByOrderNumber(orderNumber string) (result string, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetOrderUidByOrderNumberQuery(orderNumber)
	result, err = helpers.GetNeo4jSingleRecordSingleValue[string](session, query)

	return result, err
}

func (svc *OrdersService) GetOrdersForCatalogueItem(catalogueItemUid string, facilityCode string) (result []models.OrderListItem, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetOrdersForCatalogueItemQuery(catalogueItemUid, facilityCode)
	result, err = helpers.GetNeo4jArrayOfNodes[models.OrderListItem](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}

func (svc *OrdersService) GetMinAndMaxOrderLinePrice(facilityCode string) (result models.OrderLineMinMax, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetMinAndMaxOrderLinePriceQuery(facilityCode)
	result, err = helpers.GetNeo4jSingleRecordAndMapToStruct[models.OrderLineMinMax](session, query)

	return result, err
}
