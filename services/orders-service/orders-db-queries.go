package ordersService

import (
	"fmt"
	"panda/apigateway/helpers"
	"panda/apigateway/services/orders-service/models"
	"strings"

	"github.com/google/uuid"
)

func GetOrderStatusesCodebookQuery() (result helpers.DatabaseQuery) {
	result.Query = `MATCH(r:OrderStatus) RETURN {uid: r.uid,name:r.name} as orderStatuses ORDER BY orderStatuses.sortOrder ASC`
	result.ReturnAlias = "orderStatuses"
	result.Parameters = make(map[string]interface{})
	return result
}

func GetSuppliersAutoCompleteQuery(searchString string, limit int) (result helpers.DatabaseQuery) {
	result.Query = `MATCH(s:Supplier) WHERE toLower(s.name) CONTAINS toLower($searchString) RETURN {uid: s.uid,name:s.name + case when s.CIN is not null then " (" + s.CIN + ")" else "" end } as suppliers ORDER BY suppliers.name ASC LIMIT $limit`
	result.ReturnAlias = "suppliers"
	result.Parameters = make(map[string]interface{})
	result.Parameters["searchString"] = searchString
	result.Parameters["limit"] = limit
	return result
}

func GetOrdersBySearchTextFullTextQuery(searchString string, facilityCode string, pagination *helpers.Pagination, sorting *[]helpers.Sorting) (result helpers.DatabaseQuery) {

	if searchString == "" {
		result.Query = "MATCH(f:Facility{code: $facilityCode}) WITH f MATCH(o:Order{deleted:false})-[:BELONGS_TO_FACILITY]->(f) WITH o "
	} else {
		result.Query = `
		CALL db.index.fulltext.queryNodes('searchIndexOrders', $search) YIELD node AS o WHERE o:Order AND o.deleted = false WITH o
		MATCH(f:Facility{code: $facilityCode}) WITH f, o
		MATCH(o)-[:BELONGS_TO_FACILITY]->(f)
		WITH o `
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
	orderStatus: os.name,
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
	result.Parameters = make(map[string]interface{})
	result.Parameters["search"] = strings.ToLower(searchString)
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

func GetOrdersBySearchTextFullTextCountQuery(searchString string, facilityCode string) (result helpers.DatabaseQuery) {

	if searchString == "" {
		result.Query = "MATCH(f:Facility{code: $facilityCode}) WITH f MATCH(o:Order{deleted:false})-[:BELONGS_TO_FACILITY]->(f) WITH o "
	} else {
		result.Query = `
		CALL db.index.fulltext.queryNodes('searchIndexOrders', $search) YIELD node AS o WHERE o:Order AND o.deleted = false WITH o
		MATCH(f:Facility{code: $facilityCode}) WITH f, o
		MATCH(o)-[:BELONGS_TO_FACILITY]->(f)
		WITH o `
	}

	result.Query += `	
	OPTIONAL MATCH (o)-[:HAS_SUPPLIER]->(s)  
	OPTIONAL MATCH (o)-[:HAS_ORDER_STATUS]->(os)
	OPTIONAL MATCH (o)-[:HAS_REQUESTOR]->(req)
	OPTIONAL MATCH (o)-[:HAS_PROCUREMENT_RESPONSIBLE]->(proc)
		
    return count(o) as count
`
	result.ReturnAlias = "count"
	result.Parameters = make(map[string]interface{})
	result.Parameters["search"] = strings.ToLower(searchString)
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
	WITH o, s,os, ol, itm, ci, req, proc
	OPTIONAL MATCH (parentSystem)-[:HAS_SUBSYSTEM]->(sys)-[:CONTAINS_ITEM]->(itm)
	WITH o, s, os, req, proc, CASE WHEN itm IS NOT NULL THEN collect({ uid: itm.uid,  
		price: ol.price,
		currency: ol.currency, 
		name: itm.name, 
		catalogueNumber: ci.catalogueNumber, 
		catalogueUid: ci.uid, 
		system: CASE WHEN parentSystem IS NOT NULL THEN {uid: parentSystem.uid,name: parentSystem.name} ELSE NULL END }) ELSE NULL END as orderLines
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
	orderLines:  orderLines 
} AS order 
	`
	result.ReturnAlias = "order"
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func InsertNewOrderQuery(newOrder *models.OrderDetail, facilityCode string, userUID string) (result helpers.DatabaseQuery) {
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

	if newOrder.OrderLines != nil && len(newOrder.OrderLines) > 0 {
		result.Query += `WITH o MATCH(ccg:CatalogueCategory{uid: $catalogueCategoryGeneralUID}) WITH o, ccg `

		result.Parameters["catalogueCategoryGeneralUID"] = CATALOGUE_CATEGORY_GENERAL_UID

		for idxLine, orderLine := range newOrder.OrderLines {

			// the item is everytime new so we create a new one and the edge HAS_ORDER_LINE will have the price and lastUpdateTime
			result.Query += fmt.Sprintf(`
			MERGE (o)-[:HAS_ORDER_LINE{price: $price%[1]v, currency: $currency%[1]v, lastUpdateTime: datetime() }]->(itm%[1]v:Item{uid: $itemUID%[1]v, name: $itemName%[1]v, serialNumber: $serialNumber%[1]v, lastUpdateTime: datetime() }) 
			WITH o,ccg, itm%[1]v `, idxLine)

			result.Parameters[fmt.Sprintf("price%v", idxLine)] = orderLine.Price
			result.Parameters[fmt.Sprintf("currency%v", idxLine)] = orderLine.Currency
			result.Parameters[fmt.Sprintf("itemUID%v", idxLine)] = uuid.New().String()
			result.Parameters[fmt.Sprintf("itemName%v", idxLine)] = orderLine.Name
			result.Parameters[fmt.Sprintf("serialNumber%v", idxLine)] = orderLine.SerialNumber

			// assign system to the item only  if system(techn. unit) is set
			if orderLine.System != nil {
				result.Query += fmt.Sprintf(`MATCH(parentSystem%[1]v:System{uid: $systemUID%[1]v})  MERGE(parentSystem%[1]v)-[:HAS_SUBSYSTEM]->(sys%[1]v:System{ uid: $newSystemUID%[1]v, name: $itemName%[1]v  })-[:CONTAINS_ITEM]->(itm%[1]v) WITH o, ccg, itm%[1]v `, idxLine)

				result.Parameters[fmt.Sprintf("systemUID%v", idxLine)] = orderLine.System.UID
				result.Parameters[fmt.Sprintf("newSystemUID%v", idxLine)] = uuid.New().String()
			}

			newCatalogueItemUIDs := make(map[string]string, 0)
			// if catalogue item is not set, create new catalogue item
			if orderLine.CatalogueUID == "" {
				if newCatalogueItemUIDs[orderLine.CatalogueNumber] == "" {
					newCatalogueItemUIDs[orderLine.CatalogueNumber] = uuid.New().String()
				}
				result.Query += fmt.Sprintf(`MERGE (ci%[1]v:CatalogueItem{ name: $itemName%[1]v, catalogueNumber: $catalogueNumber%[1]v }) WITH o, itm%[1]v, ci%[1]v, ccg `, idxLine)
				result.Query += fmt.Sprintf(`SET ci%[1]v.uid = $catalogueItemUID%[1]v, ci%[1]v.lastUpdateTime = datetime() WITH o, itm%[1]v, ci%[1]v, ccg `, idxLine)
				result.Query += fmt.Sprintf(`MERGE (itm%[1]v)-[:IS_BASED_ON]->(ci%[1]v) WITH o, itm%[1]v, ci%[1]v, ccg `, idxLine)
				result.Query += fmt.Sprintf(`MERGE (ci%[1]v)-[:BELONGS_TO_CATEGORY]->(ccg) WITH o,ccg, itm%[1]v, ci%[1]v `, idxLine)

				result.Parameters[fmt.Sprintf("catalogueItemUID%v", idxLine)] = newCatalogueItemUIDs[orderLine.CatalogueNumber]
				result.Parameters[fmt.Sprintf("catalogueNumber%v", idxLine)] = orderLine.CatalogueNumber

			} else {
				result.Query += fmt.Sprintf(`MATCH(ci%[1]v:CatalogueItem{uid: $catalogueItemUID%[1]v}) WITH o,ccg, itm%[1]v, ci%[1]v `, idxLine)

				result.Parameters[fmt.Sprintf("catalogueItemUID%v", idxLine)] = orderLine.CatalogueUID
			}

			result.Query += fmt.Sprintf(`MERGE (itm%[1]v)-[:IS_BASED_ON]->(ci%[1]v) `, idxLine)

		}
	}

	result.Query += `
	RETURN DISTINCT o.uid as uid
	`
	result.ReturnAlias = "uid"

	result.Parameters["uid"] = uuid.New().String()
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["name"] = newOrder.Name
	result.Parameters["orderNumber"] = newOrder.OrderNumber
	result.Parameters["requestNumber"] = newOrder.RequestNumber
	result.Parameters["contractNumber"] = newOrder.ContractNumber
	result.Parameters["notes"] = newOrder.Notes
	result.Parameters["orderDate"] = newOrder.OrderDate.Local()
	result.Parameters["lastUpdateBy"] = userUID

	return result
}

// update order query
func UpdateOrderQuery(newOrder *models.OrderDetail, oldOrder *models.OrderDetail, facilityCode string, userUID string) (result helpers.DatabaseQuery) {
	result.Parameters = make(map[string]interface{}, 0)

	result.Query = `
	MATCH (o:Order{uid: $uid})-[:BELONGS_TO_FACILITY]->(f:Facility{code: $facilityCode}) 
	`

	helpers.AutoResolveObjectToUpdateQuery(&result, *newOrder, *oldOrder, "o")

	//compare new and old order lines and delete the ones that are not in the new order
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

func UpdateOrderLineDeliveryQuery(itemUID string, isDelivered bool, userUID string) (result helpers.DatabaseQuery) {
	result.Query = `
	MATCH(u:User{uid: $userUID})
	WITH u
	MATCH(o)-[ol:HAS_ORDER_LINE]->(itm:Item{uid: $itemUid})
	SET ol.isDelivered = $isDelivered, ol.deliveredTime = datetime(), ol.lastUpdateTime = datetime(), ol.lastUpdateBy = u.username, o.lastUpdateTime = datetime(), o.lastUpdateBy = u.username
	WITH o, u
	CREATE(o)-[:WAS_UPDATED_BY{at: datetime(), action: "UPDATE" }]->(u)

     return $itemUID as uid`

	result.ReturnAlias = "uid"
	result.Parameters = make(map[string]interface{})
	result.Parameters["itemUid"] = itemUID
	result.Parameters["userUID"] = userUID
	result.Parameters["isDelivered"] = isDelivered

	return result
}

func GetItemsForEunPrintQuery(euns []string) (result helpers.DatabaseQuery) {

	result.ReturnAlias = "items"
	result.Parameters = make(map[string]interface{})

	if len(euns) == 0 {
		result.Query = `
		MATCH (o:Order)-[:HAS_ORDER_LINE]->(itm:Item)-[:IS_BASED_ON]->(ci:CatalogueItem) WHERE itm.printEUN = true
	WITH o, itm, ci
	OPTIONAL MATCH (ci)-[:HAS_MANUFACTURER]->(man)
	WITH o, itm, ci, man
	RETURN DISTINCT {
		eun: itm.eun,
		name: itm.name,
		catalogueNumber: ci.catalogueNumber,
		serialNumber: itm.serialNumber,
		manufacturer: man.manufacturer,
		quantity: 1
	} as items `

	} else {
		result.Query = `	
	MATCH (o:Order)-[:HAS_ORDER_LINE]->(itm:Item)-[:IS_BASED_ON]->(ci:CatalogueItem) WHERE itm.eun IN $euns
	WITH o, itm, ci
	OPTIONAL MATCH (ci)-[:HAS_MANUFACTURER]->(man)
	WITH o, itm, ci, man
	RETURN DISTINCT {
		eun: itm.eun,
		name: itm.name,
		catalogueNumber: ci.catalogueNumber,
		serialNumber: itm.serialNumber,
		manufacturer: man.manufacturer,
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

const CATALOGUE_CATEGORY_GENERAL_UID string = "97598f04-948f-4da5-95b6-b2a44e0076db"
