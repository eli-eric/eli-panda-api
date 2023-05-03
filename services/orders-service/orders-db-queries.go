package ordersService

import (
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
	result.Query = `MATCH(s:Supplier) WHERE toLower(s.name) CONTAINS toLower($searchString) RETURN {uid: s.uid,name:s.name} as suppliers ORDER BY suppliers.name ASC LIMIT $limit`
	result.ReturnAlias = "suppliers"
	result.Parameters = make(map[string]interface{})
	result.Parameters["searchString"] = searchString
	result.Parameters["limit"] = limit
	return result
}

func GetOrdersBySearchTextFullTextQuery(searchString string, facilityCode string, pagination *helpers.Pagination, sorting *[]helpers.Sorting) (result helpers.DatabaseQuery) {

	if searchString == "" {
		result.Query = "MATCH(f:Facility{code: $facilityCode}) WITH f MATCH(o:Order)-[:BELONGS_TO_FACILITY]->(f) WITH o "
	} else {
		result.Query = `
		CALL db.index.fulltext.queryNodes('searchIndexOrders', $search) YIELD node AS o WHERE o:Order 
		MATCH(f:Facility{code: $facilityCode}) WITH f, o
		MATCH(o)-[:BELONGS_TO_FACILITY]->(f)
		WITH o `
	}

	result.Query += `	
	OPTIONAL MATCH (o)-[:HAS_SUPPLIER]->(s)  
	OPTIONAL MATCH (o)-[:HAS_ORDER_STATUS]->(os)
	RETURN DISTINCT {  
	uid: o.uid,
	name: o.name,
	orderNumber: o.orderNumber,
	requestNumber: o.requestNumber,
	contractNumber: o.contractNumber,
	orderDate: o.orderDate,
	supplier: s.name,
	orderStatus: os.name,
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
		return `ORDER BY orders.orderDate DESC `
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
		result.Query = "MATCH(f:Facility{code: $facilityCode}) WITH f MATCH(o:Order)-[:BELONGS_TO_FACILITY]->(f) WITH o "
	} else {
		result.Query = `
		CALL db.index.fulltext.queryNodes('searchIndexOrders', $search) YIELD node AS o WHERE o:Order 
		MATCH(f:Facility{code: $facilityCode}) WITH f, o
		MATCH(o)-[:BELONGS_TO_FACILITY]->(f)
		WITH o `
	}

	result.Query += `	
	OPTIONAL MATCH (o)-[:HAS_SUPPLIER]->(s)  
	OPTIONAL MATCH (o)-[:HAS_ORDER_STATUS]->(os)
		
    return count(o) as count
`
	result.ReturnAlias = "count"
	result.Parameters = make(map[string]interface{})
	result.Parameters["search"] = strings.ToLower(searchString)
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func GetOrderWithOrderLinesByUidQuery(uid string) (result helpers.DatabaseQuery) {
	result.Query = `
	MATCH(o:Order {uid: $uid})
	WITH o
	OPTIONAL MATCH (o)-[:HAS_SUPPLIER]->(s)  
	OPTIONAL MATCH (o)-[:HAS_ORDER_STATUS]->(os)
	OPTIONAL MATCH (o)-[ol:HAS_ORDER_LINE]->(itm)-[:IS_BASED_ON]->(ci)
	WITH o, s,os, ol, itm, ci 
	OPTIONAL MATCH (sys)-[:CONTAINS_ITEM]->(itm)
	WITH o, s, os, CASE WHEN itm IS NOT NULL THEN collect({ uid: itm.uid,  
		priceEur: ol.priceEur, 
		name: itm.name, 
		catalogueNumber: ci.catalogueNumber, 
		catalogueUid: ci.uid, 
		system: CASE WHEN sys IS NOT NULL THEN {uid: sys.uid,name: sys.name} ELSE NULL END }) ELSE NULL END as orderLines
	RETURN DISTINCT {  
	uid: o.uid,
	name: o.name,
	orderNumber: o.orderNumber,
	requestNumber: o.requestNumber,
	contractNumber: o.contractNumber,
	notes: o.notes,
	supplier: CASE WHEN s IS NOT NULL THEN {uid: s.uid,name: s.name} ELSE NULL END,
	orderStatus: CASE WHEN s IS NOT NULL THEN {uid: os.uid,name: os.name} ELSE NULL END,
	orderDate: o.orderDate,
	orderLines:  orderLines 
} AS order 
	`
	result.ReturnAlias = "order"
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid
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
		lastUpdateBy: u.username
	})-[:BELONGS_TO_FACILITY]->(f)
	with o,u
	CREATE(o)-[:WAS_CHANGED_BY{ updateTime: datetime() }]->(u)	
	`
	if newOrder.Supplier != nil {
		result.Query += `WITH o MATCH(s:Supplier{uid: $supplierUID}) MERGE (o)-[:HAS_SUPPLIER]->(s) `
		result.Parameters["supplierUID"] = newOrder.Supplier.UID
	}

	if newOrder.OrderStatus != nil {
		result.Query += `WITH o MATCH(os:OrderStatus{uid: $orderStatusUID}) MERGE (o)-[:HAS_ORDER_STATUS]->(os) `
		result.Parameters["orderStatusUID"] = newOrder.OrderStatus.UID
	}

	// if newOrder.OrderLines != nil && len(newOrder.OrderLines) > 0 {
	// 	result.Query += `WITH o `
	// 	for idxLine, orderLine := range newOrder.OrderLines {

	// 		result.Query += fmt.Sprintf(`MERGE (o)-[:HAS_ORDER_LINE{priceEur: $priceEur%[1]v}]->(itm:Item{uid: $itemUID%[1]v}) `, idxLine)
	// 		// if catalogue item is not set, create new catalogue item
	// 		if orderLine.CatalogueUID == "" {

	// 		}
	// 		result.Query += `MATCH(ci:CatalogueItem{uid: $catalogueItemUID` + orderLine.UID + `}) `

	// 		result.Query += `MERGE (itm)-[:IS_BASED_ON]->(ci) `
	// 		result.Parameters["catalogueItemUID"+orderLine.UID] = orderLine.CatalogueUID
	// 		result.Parameters[fmt.Sprintf("priceEur%v", idxLine)] = orderLine.PriceEUR
	// 		result.Parameters[fmt.Sprintf("itemUID%v", idxLine)] = orderLine.UID
	// 	}
	// }

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
	result.Parameters["orderDate"] = newOrder.OrderDate.UTC()
	result.Parameters["lastUpdateBy"] = userUID

	return result
}
