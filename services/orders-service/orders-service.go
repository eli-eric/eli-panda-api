package ordersService

import (
	"encoding/json"

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
	GetOrdersWithSearchAndPagination(search string, facilityCode string, pagination *helpers.Pagination, sorting *[]helpers.Sorting, filtering *[]helpers.ColumnFilter) (result helpers.PaginationResult[models.OrderListItem], err error)
	GetOrderWithOrderLinesByUid(orderUid string, facilityCode string) (result models.OrderDetail, err error)
	InsertNewOrder(order *models.OrderDetail, facilityCode string, userUID string) (uid string, err error)
	UpdateOrder(order *models.OrderDetail, facilityCode string, userUID string) (err error)
	DeleteOrder(orderUid string, userUID string) (err error)
	UpdateOrderLineDelivery(itemUid string, isDelivered bool, serialNumber *string, eun *string, userUID string, facilityCode string) (result models.OrderLine, err error)
	UpdateMultipleOrderLineDelivery(itemUids []string, userUID string, facilityCode string) (result []models.OrderLine, err error)
	GetItemsForEunPrint(euns []string) (result []models.ItemForEunPrint, err error)
	SetItemPrintEUN(eun string, printEUN bool) (err error)
	GetOrderUidByOrderNumber(orderNumber string) (result string, err error)
	GetOrdersForCatalogueItem(catalogueItemUid string, facilityCode string) (result []models.OrderListItem, err error)
	GetMinAndMaxOrderLinePrice(facilityCode string) (result models.OrderLineMinMax, err error)
	UpdateServiceLineDelivery(serviceItemUid string, isDelivered bool, userUID string, facilityCode string) (result models.ServiceLine, err error)
	UpdateMultipleServiceLineDelivery(serviceItemUids []string, userUID string, facilityCode string) (result []models.ServiceLine, err error)
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

func (svc *OrdersService) GetOrdersWithSearchAndPagination(search string, facilityCode string, pagination *helpers.Pagination, sorting *[]helpers.Sorting, filtering *[]helpers.ColumnFilter) (result helpers.PaginationResult[models.OrderListItem], err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetOrdersBySearchTextFullTextQuery(search, facilityCode, pagination, sorting, filtering)
	items, err := helpers.GetNeo4jArrayOfNodes[models.OrderListItem](session, query)
	totalCount, _ := helpers.GetNeo4jSingleRecordSingleValue[int64](session, GetOrdersBySearchTextFullTextCountQuery(search, facilityCode, filtering))

	result = helpers.GetPaginationResult(items, int64(totalCount), err)

	return result, err
}

func (svc *OrdersService) GetOrderWithOrderLinesByUid(orderUid string, facilityCode string) (result models.OrderDetail, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetOrderWithOrderLinesByUidQuery(orderUid, facilityCode)
	result, err = helpers.GetNeo4jSingleRecordAndMapToStruct[models.OrderDetail](session, query)

	// Konverze string hodnot zpět na objekty pro service lines
	if err == nil && result.ServiceLines != nil {
		for i, serviceLine := range result.ServiceLines {
			for j, detail := range serviceLine.Details {
				if strValue, ok := detail.Value.(string); ok {
					// Pokus o deserializaci JSON stringu
					var jsonValue interface{}
					if err := json.Unmarshal([]byte(strValue), &jsonValue); err == nil {
						result.ServiceLines[i].Details[j].Value = jsonValue
					}
				}
			}
		}
	}

	return result, err
}

func (svc *OrdersService) InsertNewOrder(order *models.OrderDetail, facilityCode string, userUID string) (uid string, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	// Validate system existence for all order lines
	if order.OrderLines != nil {
		for i := range order.OrderLines {
			if err := svc.validateOrderLineSystemExists(&order.OrderLines[i], facilityCode, session); err != nil {
				return "", err
			}
		}
	}

	queries := []helpers.DatabaseQuery{}
	generalQuery, newUid := InsertNewOrderQuery(order, facilityCode, userUID)
	queries = append(queries, generalQuery)

	if order.OrderLines != nil && len(order.OrderLines) > 0 {
		for _, orderLine := range order.OrderLines {
			orderLineQuery := InsertNewOrderOrderLineQuery(newUid, &orderLine, facilityCode, userUID)
			queries = append(queries, orderLineQuery)
		}
	}

	if order.ServiceLines != nil && len(order.ServiceLines) > 0 {
		for _, serviceLine := range order.ServiceLines {
			serviceLineQuery := InsertNewServiceLineQuery(newUid, &serviceLine, facilityCode, userUID)
			queries = append(queries, serviceLineQuery)
		}
	}

	deliveryStatusQuery := UpdateOrderDeliveryStatusQuery(newUid, facilityCode)
	queries = append(queries, deliveryStatusQuery)

	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session, queries...)

	return newUid, err
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

		if err != nil {
			log.Error().Msg(err.Error())
			return err
		}

		// temporary disabled
		// if oldOrder.LastUpdateTime != order.LastUpdateTime {
		// 	log.Err(helpers.ERR_CONFLICT).Msg("Order was updated by another user")
		// 	return helpers.ERR_CONFLICT
		// }

		// Validate order lines
		if len(order.OrderLines) > 0 {
			for i := range order.OrderLines {
				newLine := &order.OrderLines[i]
				if newLine.UID != "" {
					// Existing order line - prevent system changes
					oldLine := findOrderLineByUID(oldOrder.OrderLines, newLine.UID)
					if oldLine != nil {
						if !systemsEqual(oldLine.System, newLine.System) {
							return helpers.BadRequest("Cannot change system for existing order line with UID " + newLine.UID)
						}
						if !systemsEqual(oldLine.ParentSystem, newLine.ParentSystem) {
							return helpers.BadRequest("Cannot change parent system for existing order line with UID " + newLine.UID)
						}
					}
				} else {
					// New order line - validate system existence
					if err := svc.validateOrderLineSystemExists(newLine, facilityCode, session); err != nil {
						return err
					}
				}
			}
		}

		queries := []helpers.DatabaseQuery{}
		genralUpdateQuery, additionalQueries := UpdateOrderQuery(order, &oldOrder, facilityCode, userUID)
		queries = append(queries, genralUpdateQuery)
		queries = append(queries, additionalQueries...)

		if len(order.OrderLines) > 0 {
			for _, orderLine := range order.OrderLines {
				orderLineQuery := UpdateOrderLineQuery(order.UID, &orderLine, facilityCode, userUID)
				queries = append(queries, orderLineQuery)
			}
		}

		if len(order.ServiceLines) > 0 {
			for _, serviceLine := range order.ServiceLines {
				serviceLineQuery := UpdateServiceLineQuery(order.UID, &serviceLine, facilityCode, userUID)
				queries = append(queries, serviceLineQuery)
			}
		}

		deleteOrderLinesQuery := DeleteOrderLinesQuery(order, &oldOrder, facilityCode, userUID)
		queries = append(queries, deleteOrderLinesQuery)

		deleteServiceLinesQuery := DeleteServiceLinesQuery(order, &oldOrder, facilityCode, userUID)
		queries = append(queries, deleteServiceLinesQuery)

		deliveryStatusQuery := UpdateOrderDeliveryStatusQuery(order.UID, facilityCode)
		queries = append(queries, deliveryStatusQuery)

		return helpers.WriteNeo4jAndReturnNothingMultipleQueries(session, queries...)

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

func (svc *OrdersService) UpdateMultipleOrderLineDelivery(itemUids []string, userUID string, facilityCode string) (result []models.OrderLine, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	for _, itemUid := range itemUids {
		query := UpdateOrderLineDeliveryQuery(itemUid, true, nil, nil, userUID, facilityCode)
		orderLine, err := helpers.WriteNeo4jReturnSingleRecordAndMapToStruct[models.OrderLine](session, query)
		if err != nil {
			return nil, err
		}
		result = append(result, orderLine)
	}

	return result, err
}

func (svc *OrdersService) UpdateServiceLineDelivery(serviceItemUid string, isDelivered bool, userUID string, facilityCode string) (result models.ServiceLine, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := UpdateServiceLineDeliveryQuery(serviceItemUid, isDelivered, nil, nil, userUID, facilityCode)
	result, err = helpers.WriteNeo4jReturnSingleRecordAndMapToStruct[models.ServiceLine](session, query)

	return result, err
}

func (svc *OrdersService) UpdateMultipleServiceLineDelivery(serviceItemUids []string, userUID string, facilityCode string) (result []models.ServiceLine, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	for _, serviceItemUid := range serviceItemUids {
		query := UpdateServiceLineDeliveryQuery(serviceItemUid, true, nil, nil, userUID, facilityCode)
		serviceLine, err := helpers.WriteNeo4jReturnSingleRecordAndMapToStruct[models.ServiceLine](session, query)
		if err != nil {
			return nil, err
		}
		result = append(result, serviceLine)
	}

	return result, err
}

// checkSystemExists validates if a system exists in the specified facility
func (svc *OrdersService) checkSystemExists(systemUID string, facilityCode string, session neo4j.Session) (bool, error) {
	query := helpers.DatabaseQuery{
		Query: `
			MATCH (s:System {uid: $systemUID, deleted: false})-[:BELONGS_TO_FACILITY]->(f:Facility {code: $facilityCode})
			RETURN count(s) > 0 as exists
		`,
		Parameters:  map[string]interface{}{"systemUID": systemUID, "facilityCode": facilityCode},
		ReturnAlias: "exists",
	}

	exists, err := helpers.GetNeo4jSingleRecordSingleValue[bool](session, query)
	return exists, err
}

// validateOrderLineSystemExists validates system existence for a single order line.
// It checks System first (new logic), then ParentSystem (old logic).
// Returns nil if validation passes, or an error if the system doesn't exist.
func (svc *OrdersService) validateOrderLineSystemExists(orderLine *models.OrderLine, facilityCode string, session neo4j.Session) error {
	// Priority: check System first (new logic), then ParentSystem (old logic)
	if orderLine.System != nil && orderLine.System.UID != "" {
		exists, err := svc.checkSystemExists(orderLine.System.UID, facilityCode, session)
		if err != nil {
			log.Error().Err(err).Msg("Error checking system existence")
			return err
		}
		if !exists {
			return helpers.BadRequest("System with UID " + orderLine.System.UID + " not found")
		}
	} else if orderLine.ParentSystem != nil && orderLine.ParentSystem.UID != "" {
		exists, err := svc.checkSystemExists(orderLine.ParentSystem.UID, facilityCode, session)
		if err != nil {
			log.Error().Err(err).Msg("Error checking parent system existence")
			return err
		}
		if !exists {
			return helpers.BadRequest("Parent system with UID " + orderLine.ParentSystem.UID + " not found")
		}
	}
	return nil
}

// systemsEqual compares two Codebook pointers for equality
func systemsEqual(a, b *codebookModels.Codebook) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.UID == b.UID
}

// findOrderLineByUID finds an order line in a slice by its UID
func findOrderLineByUID(orderLines []models.OrderLine, uid string) *models.OrderLine {
	for i := range orderLines {
		if orderLines[i].UID == uid {
			return &orderLines[i]
		}
	}
	return nil
}
