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

	var result string = ` WITH distinct uid, 
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
