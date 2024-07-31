package ordersService

import (
	"fmt"
	"panda/apigateway/helpers"
	"panda/apigateway/services/orders-service/models"
	"strings"

	"github.com/google/uuid"
)

const CATALOGUE_CATEGORY_GENERAL_UID string = "97598f04-948f-4da5-95b6-b2a44e0076db"

func GetOrderStatusesCodebookQuery() (result helpers.DatabaseQuery) {
	result.Query = `MATCH(r:OrderStatus) RETURN {uid: r.uid,name:r.name} as orderStatuses ORDER BY orderStatuses.sortOrder ASC`
	result.ReturnAlias = "orderStatuses"
	result.Parameters = make(map[string]interface{})
	return result
}

func GetSuppliersAutoCompleteQuery(searchString string, limit int) (result helpers.DatabaseQuery) {
	result.Query = `MATCH(s:Supplier) WHERE apoc.text.clean(s.name) STARTS WITH apoc.text.clean($searchString) RETURN {uid: s.uid,name:s.name + case when s.CIN is not null then " (" + s.CIN + ")" else "" end } as suppliers ORDER BY suppliers.name ASC LIMIT $limit`
	result.ReturnAlias = "suppliers"
	result.Parameters = make(map[string]interface{})
	result.Parameters["searchString"] = searchString
	result.Parameters["limit"] = limit
	return result
}

func GetOrdersBySearchTextFullTextQuery(searchString string, supplierUID string, orderStatusUID string, procurementResponsibleUID string, requestorUID string, facilityCode string, pagination *helpers.Pagination, sorting *[]helpers.Sorting) (result helpers.DatabaseQuery) {

	result.Parameters = make(map[string]interface{})

	//beacause of the full text search we need to modify the search string
	searchFulltext := helpers.GetFullTextSearchString(searchString)

	if searchString == "" {
		result.Query = "MATCH(f:Facility{code: $facilityCode}) WITH f MATCH(o:Order{deleted:false})-[:BELONGS_TO_FACILITY]->(f) WITH o "
	} else {
		result.Query = `
		CALL {
			CALL db.index.fulltext.queryNodes('searchIndexOrders', $searchFulltext) YIELD node AS o WHERE o:Order AND o.deleted = false return o
			UNION
			MATCH(o)-[:HAS_ORDER_LINE]->(itm) where o.deleted = false AND toLower(itm.eun) STARTS WITH $searchString return o
		}
		WITH o
		MATCH(f:Facility{code: $facilityCode}) WITH f, o
		MATCH(o)-[:BELONGS_TO_FACILITY]->(f)
		WITH o `
	}

	if supplierUID != "" {
		result.Query += `MATCH(o)-[:HAS_SUPPLIER]->(ss:Supplier{uid: $supplierUID}) `
		result.Parameters["supplierUID"] = supplierUID
	}

	if orderStatusUID != "" {
		result.Query += `MATCH(o)-[:HAS_ORDER_STATUS]->(oss:OrderStatus{uid: $orderStatusUID}) `
		result.Parameters["orderStatusUID"] = orderStatusUID
	}

	if procurementResponsibleUID != "" {
		result.Query += `MATCH(o)-[:HAS_PROCUREMENT_RESPONSIBLE]->(procs:Employee{uid: $procurementResponsibleUID}) `
		result.Parameters["procurementResponsibleUID"] = procurementResponsibleUID
	}

	if requestorUID != "" {
		result.Query += `MATCH(o)-[:HAS_REQUESTOR]->(reqs:Employee{uid: $requestorUID}) `
		result.Parameters["requestorUID"] = requestorUID
	}

	result.Query += `	
	OPTIONAL MATCH (o)-[:HAS_SUPPLIER]->(s)  
	OPTIONAL MATCH (o)-[:HAS_ORDER_STATUS]->(os)
	OPTIONAL MATCH (o)-[:HAS_REQUESTOR]->(req)
	OPTIONAL MATCH (o)-[:HAS_PROCUREMENT_RESPONSIBLE]->(proc)
	RETURN DISTINCT {  
	uid: o.uid,
	name: o.name,
	orderNumber: o.orderNumber,
	requestNumber: o.requestNumber,
	contractNumber: o.contractNumber,
	orderDate: o.orderDate,
	supplier: s.name,
	orderStatus: case when os is not null then { uid: os.uid,name: os.name, code: os.code} else null end ,
	deliveryStatus: o.deliveryStatus,
	requestor: req.lastName + " " + req.firstName,
	procurementResponsible: proc.lastName + " " + proc.firstName,
	notes: o.notes,
	lastUpdateTime: o.lastUpdateTime,
	lastUpdateBy: o.lastUpdateBy
} AS orders

` + GetOrdersOrderByClauses(sorting) + `

	SKIP $skip
	LIMIT $limit

`
	result.ReturnAlias = "orders"

	result.Parameters["searchFulltext"] = strings.ToLower(searchFulltext)
	result.Parameters["searchString"] = strings.ToLower(searchString)
	result.Parameters["limit"] = pagination.PageSize
	result.Parameters["skip"] = (pagination.Page - 1) * pagination.PageSize
	result.Parameters["facilityCode"] = facilityCode

	return result
}

func GetOrdersOrderByClauses(sorting *[]helpers.Sorting) string {

	if sorting == nil || len(*sorting) == 0 {
		return `ORDER BY orders.lastUpdateTime DESC `
	}

	var result string = ` ORDER BY `

	for i, sort := range *sorting {
		if i > 0 {
			result += ", "
		}
		result += "orders." + sort.ID
		if sort.DESC {
			result += " DESC "
		}
	}

	return result
}

func GetOrdersBySearchTextFullTextCountQuery(searchString string, supplierUID string, orderStatusUID string, procurementResponsibleUID string, requestorUID string, facilityCode string) (result helpers.DatabaseQuery) {

	result.Parameters = make(map[string]interface{})

	//beacause of the full text search we need to modify the search string
	searchFulltext := helpers.GetFullTextSearchString(searchString)

	if searchString == "" {
		result.Query = "MATCH(f:Facility{code: $facilityCode}) WITH f MATCH(o:Order{deleted:false})-[:BELONGS_TO_FACILITY]->(f) WITH o "
	} else {
		result.Query = `
		CALL {
			CALL db.index.fulltext.queryNodes('searchIndexOrders', $searchFulltext) YIELD node AS o WHERE o:Order AND o.deleted = false return o
			UNION
			MATCH(o)-[:HAS_ORDER_LINE]->(itm) where o.deleted = false AND toLower(itm.eun) = $searchString return o
		}
		WITH o
		MATCH(f:Facility{code: $facilityCode}) WITH f, o
		MATCH(o)-[:BELONGS_TO_FACILITY]->(f)
		WITH o `
	}

	if supplierUID != "" {
		result.Query += `MATCH(o)-[:HAS_SUPPLIER]->(ss:Supplier{uid: $supplierUID}) `
		result.Parameters["supplierUID"] = supplierUID
	}

	if orderStatusUID != "" {
		result.Query += `MATCH(o)-[:HAS_ORDER_STATUS]->(oss:OrderStatus{uid: $orderStatusUID}) `
		result.Parameters["orderStatusUID"] = orderStatusUID
	}

	if procurementResponsibleUID != "" {
		result.Query += `MATCH(o)-[:HAS_PROCUREMENT_RESPONSIBLE]->(procs:User{uid: $procurementResponsibleUID}) `
		result.Parameters["procurementResponsibleUID"] = procurementResponsibleUID
	}

	if requestorUID != "" {
		result.Query += `MATCH(o)-[:HAS_REQUESTOR]->(reqs:User{uid: $requestorUID}) `
		result.Parameters["requestorUID"] = requestorUID
	}

	result.Query += `	
	OPTIONAL MATCH (o)-[:HAS_SUPPLIER]->(s)  
	OPTIONAL MATCH (o)-[:HAS_ORDER_STATUS]->(os)
	OPTIONAL MATCH (o)-[:HAS_REQUESTOR]->(req)
	OPTIONAL MATCH (o)-[:HAS_PROCUREMENT_RESPONSIBLE]->(proc)
		
    return count(o) as count
`
	result.ReturnAlias = "count"

	result.Parameters["searchFulltext"] = strings.ToLower(searchFulltext)
	result.Parameters["searchString"] = strings.ToLower(searchString)
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func GetOrderWithOrderLinesByUidQuery(uid string, facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = `
	MATCH(o:Order {uid: $uid, deleted: false})-[:BELONGS_TO_FACILITY]->(f:Facility{code: $facilityCode})
	WITH o
	OPTIONAL MATCH (o)-[:HAS_SUPPLIER]->(s)  
	OPTIONAL MATCH (o)-[:HAS_ORDER_STATUS]->(os)	
	OPTIONAL MATCH (o)-[:HAS_REQUESTOR]->(req)
	OPTIONAL MATCH (o)-[:HAS_PROCUREMENT_RESPONSIBLE]->(proc)
	OPTIONAL MATCH (o)-[ol:HAS_ORDER_LINE]->(itm)-[:IS_BASED_ON]->(ci)	
	WITH o, s,os, itm, ci, req, proc, ol order by ol.isDelivered desc, ol.name
	OPTIONAL MATCH (parentSystem)-[:HAS_SUBSYSTEM]->(sys)-[:CONTAINS_ITEM]->(itm)
	OPTIONAL MATCH (itm)-[:HAS_ITEM_USAGE]->(itemUsage)
	OPTIONAL MATCH (sys)-[:HAS_LOCATION]->(loc)
	WITH o, s, os, req, proc, CASE WHEN itm IS NOT NULL THEN collect({ uid: itm.uid,  
		price: ol.price,
		currency: ol.currency, 
		notes: ol.notes,
		name: itm.name, 
		eun: itm.eun,
		serialNumber: itm.serialNumber,
		isDelivered: ol.isDelivered,
		deliveredTime: ol.deliveredTime,		
		catalogueNumber: ci.catalogueNumber, 
		catalogueUid: ci.uid, 		
		system: CASE WHEN parentSystem IS NOT NULL THEN {uid: parentSystem.uid,name: parentSystem.name} ELSE NULL END,
		location: CASE WHEN loc IS NOT NULL THEN {uid: loc.uid,name: loc.name} ELSE NULL END,
		itemUsage: CASE WHEN itemUsage IS NOT NULL THEN {uid: itemUsage.uid,name: itemUsage.name} ELSE NULL END   }) ELSE NULL END as orderLines
	RETURN DISTINCT {  
	uid: o.uid,
	name: o.name,
	orderNumber: o.orderNumber,
	requestNumber: o.requestNumber,
	contractNumber: o.contractNumber,	
	notes: o.notes,
	deliveryStatus: o.deliveryStatus,
	supplier: CASE WHEN s IS NOT NULL THEN {uid: s.uid,name: s.name} ELSE NULL END,
	orderStatus: CASE WHEN os IS NOT NULL THEN {uid: os.uid,name: os.name} ELSE NULL END,
	requestor: CASE WHEN req IS NOT NULL THEN {uid: req.uid,name: req.lastName + " " + req.firstName} ELSE NULL END,
	procurementResponsible: CASE WHEN proc IS NOT NULL THEN {uid: proc.uid,name: proc.lastName + " " + proc.firstName} ELSE NULL END,
	orderDate: o.orderDate,
	orderLines:  orderLines,
	lastUpdateTime: o.lastUpdateTime
} AS order 
	`
	result.ReturnAlias = "order"
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func InsertNewOrderQuery(newOrder *models.OrderDetail, facilityCode string, userUID string) (result helpers.DatabaseQuery, uid string) {
	result.Parameters = make(map[string]interface{})

	result.Query = `
	MATCH(f:Facility{code: $facilityCode}) 
	MATCH(u:User{uid: $lastUpdateBy})
	WITH f, u
	CREATE(o:Order { 
		uid: $uid,
		name: $name,
		orderNumber: $orderNumber,
		requestNumber: $requestNumber,
		contractNumber: $contractNumber,
		notes: $notes,
		orderDate: $orderDate,
		lastUpdateTime: datetime(),
		lastUpdateBy: u.username,
		deleted: false
	})-[:BELONGS_TO_FACILITY]->(f)
	with o,u
	CREATE(o)-[:WAS_UPDATED_BY{ at: datetime(), action: "INSERT" }]->(u)	
	`
	if newOrder.Supplier != nil {
		result.Query += `WITH o MATCH(s:Supplier{uid: $supplierUID}) MERGE (o)-[:HAS_SUPPLIER]->(s) `
		result.Parameters["supplierUID"] = newOrder.Supplier.UID
	}

	if newOrder.OrderStatus != nil {
		result.Query += `WITH o MATCH(os:OrderStatus{uid: $orderStatusUID}) MERGE (o)-[:HAS_ORDER_STATUS]->(os) `
		result.Parameters["orderStatusUID"] = newOrder.OrderStatus.UID
	}

	if newOrder.Requestor != nil {
		result.Query += `WITH o MATCH(req:Employee{uid: $requestorUID}) MERGE (o)-[:HAS_REQUESTOR]->(req) `
		result.Parameters["requestorUID"] = newOrder.Requestor.UID
	}

	if newOrder.ProcurementResponsible != nil {
		result.Query += `WITH o MATCH(proc:Employee{uid: $procurementResponsibleUID}) MERGE (o)-[:HAS_PROCUREMENT_RESPONSIBLE]->(proc) `
		result.Parameters["procurementResponsibleUID"] = newOrder.ProcurementResponsible.UID
	}

	result.Query += `
	WITH o
	MATCH(o)-[olAll:HAS_ORDER_LINE]->()
	WITH count(olAll) as totalLines, o
	OPTIONAL MATCH(o)-[olDelivered:HAS_ORDER_LINE{isDelivered: true}]->()
	WITH totalLines, count(olDelivered) as deliveredLines, o
	SET o.deliveryStatus = case when deliveredLines = 0 then 0 when deliveredLines = totalLines then 2 else 1 end
	WITH o

	RETURN DISTINCT o.uid as uid
	`
	result.ReturnAlias = "uid"

	newUid := uuid.New().String()
	result.Parameters["uid"] = newUid
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["name"] = strings.TrimSpace(newOrder.Name)
	result.Parameters["orderNumber"] = newOrder.OrderNumber
	result.Parameters["requestNumber"] = newOrder.RequestNumber
	result.Parameters["contractNumber"] = newOrder.ContractNumber
	result.Parameters["notes"] = newOrder.Notes
	result.Parameters["orderDate"] = newOrder.OrderDate.Local()
	result.Parameters["lastUpdateBy"] = userUID

	return result, newUid
}

func InsertNewOrderOrderLineQuery(orderUID string, orderLine *models.OrderLine, facilityCode string, userUID string) (result helpers.DatabaseQuery) {

	result.Parameters = make(map[string]interface{}, 0)

	result.Query = `
	MATCH (o:Order{uid: $uid})-[:BELONGS_TO_FACILITY]->(f:Facility{code: $facilityCode}) 
	MATCH(ccg:CatalogueCategory{uid: $catalogueCategoryGeneralUID}) WITH o, ccg 
	`

	result.Parameters["catalogueCategoryGeneralUID"] = CATALOGUE_CATEGORY_GENERAL_UID

	// the item is everytime new so we create a new one and the edge HAS_ORDER_LINE will have the price and lastUpdateTime
	result.Query += `
	CREATE (o)-[:HAS_ORDER_LINE{price: $price, currency: $currency, lastUpdateTime: datetime() }]->(itm:Item{uid: $itemUID, name: $itemName, serialNumber: $serialNumber, notes: $notes , lastUpdateTime: datetime() }) 
	WITH o,ccg, itm `

	result.Parameters["price"] = orderLine.Price
	result.Parameters["currency"] = orderLine.Currency
	result.Parameters["itemUID"] = uuid.New().String()
	result.Parameters["itemName"] = strings.TrimSpace(orderLine.Name)
	result.Parameters["serialNumber"] = orderLine.SerialNumber
	result.Parameters["notes"] = orderLine.Notes

	// assign system to the item only  if system(techn. unit) is set
	if orderLine.System != nil {
		result.Query += `MATCH(parentSystem:System{uid: $systemUID})  MERGE(parentSystem)-[:HAS_SUBSYSTEM]->(sys:System{ uid: $newSystemUID, deleted: false, name: $itemName, systemLevel: 'SUBSYSTEMS_AND_PARTS'  })-[:CONTAINS_ITEM]->(itm)  WITH o, ccg, itm, sys `
		result.Query += `MATCH(usr:User{uid: $lastUpdateBy}) CREATE(sys)-[:WAS_UPDATED_BY{at: datetime(), action: "CREATE" }]->(usr)  WITH o, ccg, itm, sys `
		result.Query += `MATCH(f:Facility{code: $facilityCode})  MERGE(sys)-[:BELONGS_TO_FACILITY]->(f) WITH o, ccg, itm `

		result.Parameters["systemUID"] = orderLine.System.UID
		result.Parameters["newSystemUID"] = uuid.New().String()
	}

	// assign item usage to the item only  if item usage is set
	if orderLine.ItemUsage != nil {
		result.Query += `MATCH(itemUsage:ItemUsage{uid: $itemUsageUID}) MERGE(itm)-[:HAS_ITEM_USAGE]->(itemUsage) WITH o, ccg, itm `

		result.Parameters["itemUsageUID"] = orderLine.ItemUsage.UID
	}

	newCatalogueItemUIDs := make(map[string]string, 0)
	// if catalogue item is not set, create new catalogue item
	if orderLine.CatalogueUID == "" {
		if newCatalogueItemUIDs[orderLine.CatalogueNumber] == "" {
			newCatalogueItemUIDs[orderLine.CatalogueNumber] = uuid.New().String()
		}
		result.Query += `MERGE (ci:CatalogueItem{ name: $itemName, catalogueNumber: $catalogueNumber }) WITH o, itm, ci, ccg `
		result.Query += `SET ci.uid = $catalogueItemUID, ci.lastUpdateTime = datetime() WITH o, itm, ci, ccg `
		result.Query += `MERGE (itm)-[:IS_BASED_ON]->(ci) WITH o, itm, ci, ccg `
		result.Query += `MERGE (ci)-[:BELONGS_TO_CATEGORY]->(ccg) WITH o,ccg, itm, ci `

		result.Parameters["catalogueItemUID"] = newCatalogueItemUIDs[orderLine.CatalogueNumber]
		result.Parameters["catalogueNumber"] = strings.TrimSpace(orderLine.CatalogueNumber)

	} else {
		result.Query += `MATCH(ci:CatalogueItem{uid: $catalogueItemUID}) WITH o,ccg, itm, ci `

		result.Parameters["catalogueItemUID"] = orderLine.CatalogueUID
	}

	result.Query += `MERGE (itm)-[:IS_BASED_ON]->(ci) `

	result.Query += `	
	RETURN o.uid as uid
	`

	result.Parameters["uid"] = orderUID
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["lastUpdateBy"] = userUID
	result.ReturnAlias = "uid"

	return result
}

func InsertNewOrderDeliveryStatusQuery(orderUID string, facilityCode string) (result helpers.DatabaseQuery) {
	result.Parameters = make(map[string]interface{}, 0)

	result.Query = `
	MATCH (o:Order{uid: $uid})-[:BELONGS_TO_FACILITY]->(f:Facility{code: $facilityCode}) 
	WITH o
	MATCH(o)-[olAll:HAS_ORDER_LINE]->()
	WITH count(olAll) as totalLines, o
	OPTIONAL MATCH(o)-[olDelivered:HAS_ORDER_LINE{isDelivered: true}]->()
	WITH totalLines, count(olDelivered) as deliveredLines, o
	SET o.deliveryStatus = case when deliveredLines = 0 then 0 when deliveredLines = totalLines then 2 else 1 end
	RETURN o.uid as uid
	`

	result.Parameters["uid"] = orderUID
	result.Parameters["facilityCode"] = facilityCode
	result.ReturnAlias = "uid"

	return result
}

// update order query
func UpdateOrderQuery(newOrder *models.OrderDetail, oldOrder *models.OrderDetail, facilityCode string, userUID string) (result helpers.DatabaseQuery) {
	result.Parameters = make(map[string]interface{}, 0)

	result.Query = `
	MATCH (o:Order{uid: $uid})-[:BELONGS_TO_FACILITY]->(f:Facility{code: $facilityCode}) 
	`

	helpers.AutoResolveObjectToUpdateQuery(&result, *newOrder, *oldOrder, "o")

	result.Query += `	
	WITH o
	MATCH(u:User{uid: $lastUpdateBy})
	WITH o, u
	SET o.lastUpdateTime = datetime(), o.lastUpdateBy = u.username
	WITH o, u
	CREATE(o)-[:WAS_UPDATED_BY{at: datetime(), action: "UPDATE" }]->(u)
	RETURN o.uid as uid
	`

	result.Parameters["uid"] = oldOrder.UID
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["lastUpdateBy"] = userUID
	result.ReturnAlias = "uid"

	return result
}

func UpdateOrderLineQuery(orderUid string, orderLine *models.OrderLine, facilityCode string, userUID string) (result helpers.DatabaseQuery) {

	result.Parameters = make(map[string]interface{}, 0)

	result.Query = `
	MATCH (o:Order{uid: $uid})-[:BELONGS_TO_FACILITY]->(f:Facility{code: $facilityCode}) 
	MATCH(ccg:CatalogueCategory{uid: $catalogueCategoryGeneralUID}) WITH o, ccg 
	`

	result.Parameters["catalogueCategoryGeneralUID"] = CATALOGUE_CATEGORY_GENERAL_UID

	if orderLine.UID == "" {
		// the item is everytime new so we create a new one and the edge HAS_ORDER_LINE will have the price and lastUpdateTime
		result.Query += `
			CREATE (o)-[:HAS_ORDER_LINE{price: $price, currency: $currency, lastUpdateTime: datetime() }]->(itm:Item{uid: $itemUID, name: $itemName, serialNumber: $serialNumber, notes: $notes , lastUpdateTime: datetime() }) 
			WITH o,ccg, itm `

		result.Parameters["price"] = orderLine.Price
		result.Parameters["currency"] = orderLine.Currency
		result.Parameters["itemUID"] = uuid.New().String()
		result.Parameters["itemName"] = strings.TrimSpace(orderLine.Name)
		result.Parameters["serialNumber"] = orderLine.SerialNumber
		result.Parameters["notes"] = orderLine.Notes

		// assign system to the item only  if system(techn. unit) is set
		if orderLine.System != nil {
			result.Query += `MATCH(parentSystem:System{uid: $systemUID})  MERGE(parentSystem)-[:HAS_SUBSYSTEM]->(sys:System{ uid: $newSystemUID, deleted: false, name: $itemName, systemLevel: 'SUBSYSTEMS_AND_PARTS'  })-[:CONTAINS_ITEM]->(itm)  WITH o, ccg, itm, sys `
			result.Query += `MATCH(usr:User{uid: $lastUpdateBy}) CREATE(sys)-[:WAS_UPDATED_BY{at: datetime(), action: "CREATE" }]->(usr)  WITH o, ccg, itm, sys `
			result.Query += `MATCH(f:Facility{code: $facilityCode})  MERGE(sys)-[:BELONGS_TO_FACILITY]->(f) WITH o, ccg, itm `

			result.Parameters["systemUID"] = orderLine.System.UID
			result.Parameters["newSystemUID"] = uuid.New().String()
		}

		// assign item usage to the item only  if item usage is set
		if orderLine.ItemUsage != nil {
			result.Query += `MATCH(itemUsage:ItemUsage{uid: $itemUsageUID}) MERGE(itm)-[:HAS_ITEM_USAGE]->(itemUsage) WITH o, ccg, itm `

			result.Parameters["itemUsageUID"] = orderLine.ItemUsage.UID
		}

		newCatalogueItemUIDs := make(map[string]string, 0)
		// if catalogue item is not set, create new catalogue item
		if orderLine.CatalogueUID == "" {
			if newCatalogueItemUIDs[orderLine.CatalogueNumber] == "" {
				newCatalogueItemUIDs[orderLine.CatalogueNumber] = uuid.New().String()
			}
			result.Query += `MERGE (ci:CatalogueItem{ name: $itemName, catalogueNumber: $catalogueNumber }) WITH o, itm, ci, ccg `
			result.Query += `SET ci.uid = $catalogueItemUID, ci.lastUpdateTime = datetime() WITH o, itm, ci, ccg `
			result.Query += `MERGE (itm)-[:IS_BASED_ON]->(ci) WITH o, itm, ci, ccg `
			result.Query += `MERGE (ci)-[:BELONGS_TO_CATEGORY]->(ccg) WITH o,ccg, itm, ci `

			result.Parameters["catalogueItemUID"] = newCatalogueItemUIDs[orderLine.CatalogueNumber]
			result.Parameters["catalogueNumber"] = strings.TrimSpace(orderLine.CatalogueNumber)

		} else {
			result.Query += `MATCH(ci:CatalogueItem{uid: $catalogueItemUID}) WITH o,ccg, itm, ci `

			result.Parameters["catalogueItemUID"] = orderLine.CatalogueUID
		}

		result.Query += `MERGE (itm)-[:IS_BASED_ON]->(ci) `
	} else {
		//update existing order lines
		result.Query += `WITH o, ccg MATCH (o)-[ol:HAS_ORDER_LINE]->(itm:Item{uid: $itemUID}) SET ol.price = $price, ol.currency = $currency, ol.lastUpdateTime = datetime(), itm.serialNumber = $serialNumber, ol.notes = $notes WITH o, ccg, itm `
		result.Parameters["price"] = orderLine.Price
		result.Parameters["currency"] = orderLine.Currency
		result.Parameters["itemUID"] = orderLine.UID
		result.Parameters["itemName"] = strings.TrimSpace(orderLine.Name)
		result.Parameters["serialNumber"] = orderLine.SerialNumber
		result.Parameters["notes"] = orderLine.Notes

		if orderLine.System != nil {
			//delete existing system
			result.Query += `OPTIONAL MATCH(sys)-[:CONTAINS_ITEM]->(itm) WITH o, ccg, itm, sys `

			if orderLine.Location != nil && orderLine.Location.UID != "" {
				//delete existing location
				result.Query += `OPTIONAL MATCH()<-[oldLocation:HAS_LOCATION]-(sys) DELETE oldLocation WITH o, ccg, itm, sys  `
				//then create new one
				result.Query += `MATCH(loc:Location{uid: $locationUID}) MERGE(sys)-[:HAS_LOCATION]->(loc) WITH o, ccg, itm, sys  `

				result.Parameters["locationUID"] = orderLine.Location.UID
			} else {
				//delete existing location
				result.Query += `OPTIONAL MATCH()<-[oldLocation:HAS_LOCATION]-(sys) DELETE oldLocation WITH o, ccg, itm, sys  `
			}
		}

		// assign item usage to the item only  if item usage is set
		if orderLine.ItemUsage != nil {
			//delete existing item usage relationship
			result.Query += `OPTIONAL MATCH(itm)-[itemUsageRel:HAS_ITEM_USAGE]->() DELETE itemUsageRel WITH o, ccg, itm `
			result.Query += `MATCH(itemUsage:ItemUsage{uid: $itemUsageUID}) MERGE(itm)-[:HAS_ITEM_USAGE]->(itemUsage) WITH o, ccg, itm `

			result.Parameters["itemUsageUID"] = orderLine.ItemUsage.UID
		}
	}

	result.Query += `	
	RETURN o.uid as uid
	`

	result.Parameters["uid"] = orderUid
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["lastUpdateBy"] = userUID
	result.ReturnAlias = "uid"

	return result
}

func DeleteOrderLinesQuery(newOrder *models.OrderDetail, oldOrder *models.OrderDetail, facilityCode string, userUID string) (result helpers.DatabaseQuery) {
	result.Parameters = make(map[string]interface{}, 0)

	result.Query = `
	MATCH (o:Order{uid: $uid})-[:BELONGS_TO_FACILITY]->(f:Facility{code: $facilityCode}) 
	`
	// compare new and old order lines and delete the ones that are not in the new order
	if newOrder.OrderLines != nil && len(newOrder.OrderLines) > 0 {
		for idxDelete, oldOrderLine := range oldOrder.OrderLines {
			found := false
			for _, newOrderLine := range newOrder.OrderLines {
				if oldOrderLine.UID == newOrderLine.UID {
					found = true
					break
				}
			}
			if !found {
				result.Query += fmt.Sprintf(` WITH o MATCH (o)-[:HAS_ORDER_LINE]->(itmForDelete%[1]v:Item{uid: $itemUIDForDelete%[1]v}) DETACH DELETE itmForDelete%[1]v `, idxDelete)
				result.Parameters[fmt.Sprintf("itemUIDForDelete%v", idxDelete)] = oldOrderLine.UID
			}
		}
	}

	result.Query += `
	RETURN o.uid as uid`

	result.Parameters["uid"] = oldOrder.UID
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["lastUpdateBy"] = userUID
	result.ReturnAlias = "uid"

	return result
}

func UpdateOrderDeliveryStatusQuery(orderUID string, facilityCode string) (result helpers.DatabaseQuery) {
	result.Parameters = make(map[string]interface{}, 0)

	result.Query = `
	MATCH (o:Order{uid: $uid})-[:BELONGS_TO_FACILITY]->(f:Facility{code: $facilityCode}) 
	WITH o
	MATCH(o)-[olAll:HAS_ORDER_LINE]->()
	WITH count(olAll) as totalLines, o
	OPTIONAL MATCH(o)-[olDelivered:HAS_ORDER_LINE{isDelivered: true}]->()
	WITH totalLines, count(olDelivered) as deliveredLines, o
	SET o.deliveryStatus = case when deliveredLines = 0 then 0 when deliveredLines = totalLines then 2 else 1 end
	WITH o

	RETURN o.uid as uid
	`

	result.Parameters["uid"] = orderUID
	result.Parameters["facilityCode"] = facilityCode
	result.ReturnAlias = "uid"

	return result
}

func DeleteOrderQuery(uid string, userUID string) (result helpers.DatabaseQuery) {

	result.Query = `
	MATCH(u:User{uid: $userUID})
	WITH u
	MATCH (o:Order{uid: $uid}) 	
	SET o.deleted = true, o.lastUpdateTime = datetime(), o.lastUpdateBy = u.username
	WITH o, u
	CREATE(o)-[:WAS_UPDATED_BY{at: datetime(), action: "DELETE" }]->(u)
	RETURN o.uid as uid`

	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid
	result.Parameters["userUID"] = userUID

	return result
}

func UpdateOrderLineDeliveryQuery(itemUID string, isDelivered bool, serialNumber *string, eun *string, userUID string, facilityCode string) (result helpers.DatabaseQuery) {
	result.Parameters = make(map[string]interface{})

	result.Query = `
	WITH apoc.text.split(toString(date.truncate('month', date())), '-') as dateParts
	WITH $facilityCode + right(dateParts[0],2) + dateParts[1] as eunPrefix
	MATCH(u:User{uid: $userUID})
	WITH u, eunPrefix
	MATCH(o)-[ol:HAS_ORDER_LINE]->(itm:Item{uid: $itemUID})
	WITH u, eunPrefix , o, ol, itm 
	OPTIONAL MATCH (parentSystem)-[:HAS_SUBSYSTEM]->(sys)-[:CONTAINS_ITEM]->(itm)
	OPTIONAL MATCH (itm)-[:HAS_ITEM_USAGE]->(itemUsage)
    WITH u, eunPrefix , o, ol, itm , parentSystem, itemUsage
    OPTIONAL MATCH (itm)-[:IS_BASED_ON]->(ci)	
	SET ol.isDelivered = $isDelivered, 
	    ol.deliveredTime = datetime(), 
		ol.lastUpdateTime = datetime(), 
		ol.lastUpdateBy = u.username, 
		o.lastUpdateTime = datetime(), 
		o.lastUpdateBy = u.username
	WITH o, ol, u, itm, ci, parentSystem , itemUsage, eunPrefix `

	if eun == nil {
		result.Query += `
	OPTIONAL MATCH(maxEuns:Item) WHERE maxEuns.eun STARTS WITH eunPrefix
	WITH max(maxEuns.eun) as maxEun, o, ol, u, itm, ci, parentSystem, itemUsage, eunPrefix 
	SET itm.eun = case when (itm.eun is null and ol.isDelivered = true) then
					case when maxEun is not null then 
						eunPrefix + apoc.text.lpad(toString(toInteger(right(maxEun, 4)) + 1), 4, '0') 
					else 
						eunPrefix + '0001' 
					end 
				  else 
					case when ol.isDelivered = true then 
						itm.eun 
					else 
						null 
					end 
				  end
	`
	} else {
		result.Query += `SET itm.eun = $eun `
		result.Parameters["eun"] = eun
	}

	if serialNumber != nil && *serialNumber != "" {
		result.Query += `, itm.serialNumber = $serialNumber `
		result.Parameters["serialNumber"] = serialNumber
	}

	result.Query += `
	WITH o, ol, u, itm, ci, parentSystem , itemUsage
	MATCH(o)-[olAll:HAS_ORDER_LINE]->()
	WITH count(olAll) as totalLines,o, ol, u, itm, ci, parentSystem , itemUsage 
	OPTIONAL MATCH(o)-[olDelivered:HAS_ORDER_LINE{isDelivered: true}]->()
	WITH totalLines, count(olDelivered) as deliveredLines,o, ol, u, itm, ci, parentSystem , itemUsage 
	SET o.deliveryStatus = case when deliveredLines = 0 then 0 when deliveredLines = totalLines then 2 else 1 end
	WITH o, ol, u, itm, ci, parentSystem , itemUsage 
	CREATE(o)-[:WAS_UPDATED_BY{at: datetime(), action: "UPDATE" }]->(u)
	RETURN DISTINCT { uid: itm.uid,  
			isDelivered: ol.isDelivered,
			deliveredTime: ol.deliveredTime,
			price: ol.price,
			currency: ol.currency, 
			name: itm.name, 
			catalogueNumber: ci.catalogueNumber, 
			catalogueUid: ci.uid, 
			eun:  itm.eun,
			serialNumber: itm.serialNumber,
			system: CASE WHEN parentSystem IS NOT NULL THEN {uid: parentSystem.uid,name: parentSystem.name} ELSE NULL END,
			itemUsage: CASE WHEN itemUsage IS NOT NULL THEN {uid: itemUsage.uid,name: itemUsage.name} ELSE NULL END  } as orderLine;
	 `

	result.ReturnAlias = "orderLine"

	result.Parameters["itemUID"] = itemUID
	result.Parameters["userUID"] = userUID
	result.Parameters["isDelivered"] = isDelivered
	result.Parameters["facilityCode"] = facilityCode

	return result
}

func GetItemsForEunPrintQuery(euns []string) (result helpers.DatabaseQuery) {

	result.ReturnAlias = "items"
	result.Parameters = make(map[string]interface{})

	if len(euns) == 0 {
		result.Query = `
		MATCH (o:Order)-[:HAS_ORDER_LINE]->(itm:Item)-[:IS_BASED_ON]->(ci:CatalogueItem) WHERE itm.printEUN = true
	WITH o, itm, ci
	OPTIONAL MATCH (o)-[:HAS_SUPPLIER]->(supplier)
	WITH o, itm, ci, supplier
	RETURN DISTINCT {
		eun: itm.eun,
		name: itm.name,
		catalogueNumber: ci.catalogueNumber,
		serialNumber: itm.serialNumber,
		manufacturer: supplier.name,
		quantity: 1
	} as items `

	} else {
		result.Query = `	
	MATCH (o:Order)-[:HAS_ORDER_LINE]->(itm:Item)-[:IS_BASED_ON]->(ci:CatalogueItem) WHERE itm.eun IN $euns
	WITH o, itm, ci
	OPTIONAL MATCH (o)-[:HAS_SUPPLIER]->(supplier)
	WITH o, itm, ci, supplier
	RETURN DISTINCT {
		eun: itm.eun,
		name: itm.name,
		catalogueNumber: ci.catalogueNumber,
		serialNumber: itm.serialNumber,
		manufacturer: supplier.name,
		quantity: 1
	} as items `

		result.Parameters["euns"] = euns
	}
	return result
}

// db query to set Item printEUN by item uid
func SetItemPrintEUNQuery(eun string, printEUN bool) (result helpers.DatabaseQuery) {

	result.Query = `
	MATCH (itm:Item{eun: $eun}) 
	SET itm.printEUN = $printEUN
	RETURN count(itm) as count`

	result.ReturnAlias = "count"
	result.Parameters = make(map[string]interface{})
	result.Parameters["eun"] = eun
	result.Parameters["printEUN"] = printEUN

	return result
}

func GetOrderUidByOrderNumberQuery(orderNumber string) (result helpers.DatabaseQuery) {
	result.Query = `
	MATCH (o:Order{orderNumber: $orderNumber}) 
	RETURN o.uid as uid limit 1`

	result.ReturnAlias = "uid"
	result.Parameters = make(map[string]interface{})
	result.Parameters["orderNumber"] = orderNumber

	return result
}

// db query to get all orders for a given catalogue item by catalogue item uid
func GetOrdersForCatalogueItemQuery(catalogueItemUID string, facilityCode string) (result helpers.DatabaseQuery) {

	result.Query = `
	MATCH(f:Facility{code: $facilityCode})
	WITH f
	MATCH (f)<-[:BELONGS_TO_FACILITY]-(o{deleted: false})-[:HAS_ORDER_LINE]->(itm)-[:IS_BASED_ON]->(ci{uid: $catalogueItemUID}) 
	OPTIONAL MATCH (o)-[:HAS_SUPPLIER]->(s)  
	OPTIONAL MATCH (o)-[:HAS_ORDER_STATUS]->(os)
	OPTIONAL MATCH (o)-[:HAS_REQUESTOR]->(req)
	OPTIONAL MATCH (o)-[:HAS_PROCUREMENT_RESPONSIBLE]->(proc)
	RETURN DISTINCT {  
		uid: o.uid,
		name: o.name,
		orderNumber: o.orderNumber,
		requestNumber: o.requestNumber,
		contractNumber: o.contractNumber,
		orderDate: o.orderDate,
		supplier: s.name,
		orderStatus: case when os is not null then { uid: os.uid,name: os.name, code: os.code} else null end ,
		deliveryStatus: o.deliveryStatus,
		requestor: req.lastName + " " + req.firstName,
		procurementResponsible: proc.lastName + " " + proc.firstName,
		notes: o.notes,
		lastUpdateTime: o.lastUpdateTime,
		lastUpdateBy: o.lastUpdateBy
	} AS orders`

	result.ReturnAlias = "orders"
	result.Parameters = make(map[string]interface{})
	result.Parameters["catalogueItemUID"] = catalogueItemUID
	result.Parameters["facilityCode"] = facilityCode

	return result
}

// db query to get min and max order line price for a given facitlity
func GetMinAndMaxOrderLinePriceQuery(facilityCode string) (result helpers.DatabaseQuery) {

	result.Query = `
	MATCH (f:Facility{code: $facilityCode})<-[:BELONGS_TO_FACILITY]-(o{deleted: false})-[ol:HAS_ORDER_LINE]->(physicalItem)<-[:CONTAINS_ITEM]-(sys{deleted:false})
	where ol.price is not null and ol.currency is not NULL and ol.currency <> ""
	with ol.price as price
	return {min: toInteger(min(price)), max: toInteger(max(price))} as result`

	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["facilityCode"] = facilityCode

	return result
}
