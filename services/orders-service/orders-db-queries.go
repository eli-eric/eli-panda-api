package ordersService

import (
	"panda/apigateway/helpers"
	"strings"
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

// get orders by search text with pagination and sorting
func GetOrdersBySearchTextQuery(searchString string, pagination *helpers.Pagination, sorting *[]helpers.Sorting) (result helpers.DatabaseQuery) {
	result.Query = `
	MATCH (o:Order)
	WITH  o
	OPTIONAL MATCH (o)-[:HAS_SUPPLIER]->(s)  
	OPTIONAL MATCH (o)-[:HAS_ORDER_STATUS]->(os)

	WITH
	o.uid as uid, 
	o.name as name, 
	o.orderDate as orderDate,
	o.orderNumber as orderNumber, 
	o.contractNumber as contractNumber,
	o.requestNumber as requestNumber,
	o.notes as notes,	
	o.lastUpdateTime AS lastUpdateTime, 
	o.lastUpdateBy AS lastUpdateBy, 
	s.name AS supplier, 
	os.name AS orderStatus 
	
	` + GetOrdersOrderByClauses(sorting) + `

	WHERE  
	toLower(orderStatus) contains $search or
	toLower(supplier) CONTAINS $search or 
	toLower(notes) CONTAINS $search or 
	toLower(name) CONTAINS $search OR
	toLower(orderNumber) CONTAINS $search OR
	toLower(requestNumber) CONTAINS $search OR
	toLower(contractNumber) CONTAINS $search 

	WITH uid, 
	name, 
	orderDate,
	orderNumber, 
	contractNumber,
	requestNumber,
	notes,	
	lastUpdateTime, 
	lastUpdateBy,	 
	supplier, 
	orderStatus

	SKIP $skip
	LIMIT $limit

	RETURN DISTINCT {  
		uid: uid,
		name: name,
		orderNumber: orderNumber,
		requestNumber: requestNumber,
		contractNumber: contractNumber ,
		orderDate: orderDate,
		supplier: supplier,
		orderStatus: orderStatus,
		notes: notes,
		lastUpdateTime: lastUpdateTime,
		lastUpdateBy: lastUpdateBy
		} as orders
`
	result.ReturnAlias = "orders"
	result.Parameters = make(map[string]interface{})
	result.Parameters["search"] = strings.ToLower(searchString)
	result.Parameters["limit"] = pagination.PageSize
	result.Parameters["skip"] = (pagination.Page - 1) * pagination.PageSize

	return result
}

func GetOrdersBySearchTextFullTextQuery(searchString string, pagination *helpers.Pagination, sorting *[]helpers.Sorting) (result helpers.DatabaseQuery) {

	if searchString == "" {
		result.Query = "MATCH(o:Order) WITH o "
	} else {
		result.Query = "CALL db.index.fulltext.queryNodes('searchIndexOrders', $search) YIELD node AS o WHERE o:Order "
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

` + GetOrdersOrderByClausesV2(sorting) + `

	SKIP $skip
	LIMIT $limit

`
	result.ReturnAlias = "orders"
	result.Parameters = make(map[string]interface{})
	result.Parameters["search"] = strings.ToLower(searchString)
	result.Parameters["limit"] = pagination.PageSize
	result.Parameters["skip"] = (pagination.Page - 1) * pagination.PageSize

	return result
}

func GetOrdersOrderByClauses(sorting *[]helpers.Sorting) string {

	if sorting == nil || len(*sorting) == 0 {
		return ` WITH uid, 
		name, 
		orderDate,
		orderNumber, 
		contractNumber,
		requestNumber,
		notes,	
		lastUpdateTime, 
		lastUpdateBy, 
		supplier, 
		orderStatus
		ORDER BY orderDate DESC `
	}

	var result string = ` WITH uid, 
	name, 
	orderDate,
	orderNumber, 
	contractNumber,
	requestNumber,
	notes,	
	lastUpdateTime, 
	lastUpdateBy,
	supplier, 
	orderStatus ORDER BY `

	for i, sort := range *sorting {
		if i > 0 {
			result += ", "
		}
		result += sort.ID
		if sort.DESC {
			result += " DESC "
		}
	}

	return result
}

func GetOrdersOrderByClausesV2(sorting *[]helpers.Sorting) string {

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

func GetOrdersBySearchTextCountQuery(searchString string) (result helpers.DatabaseQuery) {
	result.Query = `
	MATCH (o:Order)
	WITH o
	OPTIONAL MATCH (o)-[:HAS_SUPPLIER]->(s)  
	OPTIONAL MATCH (o)-[:HAS_ORDER_STATUS]->(os)

    WITH
	o.uid as uid, 
	o.name as name, 
	o.orderDate as orderDate,
	o.orderNumber as orderNumber, 
	o.contractNumber as contractNumber,
	o.requestNumber as requestNumber,
	o.notes as notes,	
	o.lastUpdateTime AS lastUpdateTime, 
	o.lastUpdateBy AS lastUpdateBy, 
	s.name AS supplier, 
	os.name AS orderStatus 

	WHERE  
	toLower(orderStatus) contains $search or
	toLower(supplier) CONTAINS $search or 
	toLower(notes) CONTAINS $search or 
	toLower(name) CONTAINS $search OR
	toLower(orderNumber) CONTAINS $search OR
	toLower(requestNumber) CONTAINS $search OR
	toLower(contractNumber) CONTAINS $search 
		
    return count(uid) as count
`
	result.ReturnAlias = "count"
	result.Parameters = make(map[string]interface{})
	result.Parameters["search"] = strings.ToLower(searchString)
	return result
}

func GetOrdersBySearchTextFullTextCountQuery(searchString string) (result helpers.DatabaseQuery) {

	if searchString == "" {
		result.Query = "MATCH(o:Order) WITH o "
	} else {
		result.Query = "CALL db.index.fulltext.queryNodes('searchIndexOrders', $search) YIELD node AS o WHERE o:Order "
	}

	result.Query += `	
	OPTIONAL MATCH (o)-[:HAS_SUPPLIER]->(s)  
	OPTIONAL MATCH (o)-[:HAS_ORDER_STATUS]->(os)
		
    return count(o) as count
`
	result.ReturnAlias = "count"
	result.Parameters = make(map[string]interface{})
	result.Parameters["search"] = strings.ToLower(searchString)
	return result
}
