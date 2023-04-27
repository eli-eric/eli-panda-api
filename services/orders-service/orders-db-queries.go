package ordersService

import "panda/apigateway/helpers"

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
	MATCH (o:Order)-[r:UPDATED_BY]->(u:User)
	WITH distinct o, r, u
	OPTIONAL MATCH (o)-[:HAS_SUPPLIER]->(s)  
	OPTIONAL MATCH (o)-[:HAS_ORDER_STATUS]->(os)

	WITH distinct
	o.uid as uid, 
	o.name as name, 
	o.orderDate as orderDate,
	o.orderNumber as orderNumber, 
	o.contractNumber as contractNumber,
	o.requestNumber as requestNumber,
	o.notes as notes,	
	r.updated AS lastUpdateTime, 
	u.username AS userName, 
	s.name AS supplier, 
	os.name AS orderStatus

	ORDER BY lastUpdateTime DESC

	WITH distinct
	uid, 
	name, 
	orderDate,
	orderNumber, 
	contractNumber,
	requestNumber,
	notes,	
	lastUpdateTime, 
	userName, 
	supplier, 
	orderStatus,
	COLLECT({lastUpdateTime:lastUpdateTime, userName:userName})[0] AS lastUpdate 
	
	` + GetOrdersOrderByClauses(sorting) + `

	WHERE  
	orderStatus contains $search or
  	supplier CONTAINS $search or 
  	notes CONTAINS $search or 
  	name CONTAINS $search OR
  	orderNumber CONTAINS $search OR
  	requestNumber CONTAINS $search OR
  	contractNumber CONTAINS $search 

	WITH distinct uid, 
	name, 
	orderDate,
	orderNumber, 
	contractNumber,
	requestNumber,
	notes,	
	lastUpdateTime, 
	userName, 
	supplier, 
	orderStatus,
	lastUpdate

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
		lastUpdateDate: collect(lastUpdate.lastUpdateTime)[0],
		lastUpdatedBy: lastUpdate.userName
		} as orders
`
	result.ReturnAlias = "orders"
	result.Parameters = make(map[string]interface{})
	result.Parameters["search"] = searchString
	result.Parameters["limit"] = pagination.PageSize
	result.Parameters["skip"] = (pagination.Page - 1) * pagination.PageSize

	return result
}

////prev version of query - no duplicates
// `MATCH (o:Order)-[r:UPDATED_BY]->(u:User)
// WITH o, r, u
// OPTIONAL MATCH (o)-[:HAS_SUPPLIER]->(s)
// OPTIONAL MATCH (o)-[:HAS_ORDER_STATUS]->(os)

// WITH o, r.updated AS lastUpdateTime, u.username AS userName, s, os
// ORDER BY lastUpdateTime DESC

// WITH o, o.name as name, o.orderDate as orderDate, COLLECT({lastUpdateTime:lastUpdateTime, userName:userName})[0] AS lastUpdate, s, os

// WITH o,name, orderDate, lastUpdate, s, os ORDER BY orderDate DESC

// WHERE
// os.name contains $search or
//   s.name CONTAINS $search or
//   o.notes CONTAINS $search or
//   name CONTAINS $search OR
//   o.orderNumber CONTAINS $search OR
//   o.requestNumber CONTAINS $search OR
//   o.contractNumber CONTAINS $search

// WITH o, name, orderDate, lastUpdate, s, os
// SKIP 0
// LIMIT 50

// RETURN {
// name: name,
// orderNumber: o.orderNumber,
// requestNumber: o.requestNumber,
// contractNumber: o.contractNumber ,
// orderDate: orderDate,
// supplier: s.name,
// orderStatus: os.name,
// lastUpdateDate:lastUpdate.lastUpdateTime,
// lastUpdatedBy: lastUpdate.userName
// }`

func GetOrdersOrderByClauses(sorting *[]helpers.Sorting) string {

	if sorting == nil || len(*sorting) == 0 {
		return ` WITH distinct uid, 
		name, 
		orderDate,
		orderNumber, 
		contractNumber,
		requestNumber,
		notes,	
		lastUpdateTime, 
		userName, 
		supplier, 
		orderStatus,
		lastUpdate ORDER BY orderDate DESC `
	}

	var result string = ` WITH distinct uid, 
	name, 
	orderDate,
	orderNumber, 
	contractNumber,
	requestNumber,
	notes,	
	lastUpdateTime, 
	userName, 
	supplier, 
	orderStatus,
	lastUpdate  ORDER BY `

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
	MATCH (o:Order)-[:UPDATED_BY]->(u:User)
	WITH distinct o, u
	OPTIONAL MATCH (o)-[:HAS_SUPPLIER]->(s)  
	OPTIONAL MATCH (o)-[:HAS_ORDER_STATUS]->(os)

    WITH o, os, s
	WHERE  
	os.name contains $search or
  	s.name CONTAINS $search or 
  	o.notes CONTAINS $search or 
  	o.name CONTAINS $search OR
  	o.orderNumber CONTAINS $search OR
  	o.requestNumber CONTAINS $search OR
  	o.contractNumber CONTAINS $search 
		
    return count(o) as count
`
	result.ReturnAlias = "count"
	result.Parameters = make(map[string]interface{})
	result.Parameters["search"] = searchString
	return result
}
