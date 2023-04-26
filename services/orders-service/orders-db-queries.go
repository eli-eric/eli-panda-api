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
func GetOrdersBySearchTextQuery(searchString string, pagination *helpers.Pagination, sorting *helpers.Sorting) (result helpers.DatabaseQuery) {
	result.Query = `
	MATCH (o:Order)-[r:UPDATED_BY]->(u:User)
	WITH o, r, u
	OPTIONAL MATCH (o)-[:HAS_SUPPLIER]->(s)  
	OPTIONAL MATCH (o)-[:HAS_ORDER_STATUS]->(os)

	WITH o, r, u, s, os
	ORDER BY o.orderDate desc

	WITH o, r.updated AS lastUpdateTime, u.username AS userName, s, os
	ORDER BY lastUpdateTime DESC
	WITH o, COLLECT({lastUpdateTime:lastUpdateTime, userName:userName})[0] AS lastUpdate, s, os

	WHERE
  
	os.name contains $search or
  	s.name CONTAINS $search or 
  	o.notes CONTAINS $search or 
  	o.name CONTAINS $search OR
  	o.orderNumber CONTAINS $search OR
  	o.requestNumber CONTAINS $search OR
  	o.contractNumber CONTAINS $search 

	WITH o, lastUpdate, s, os
	SKIP $skip
	LIMIT $limit

	RETURN 
  o.name AS orderName,
  o.orderNumber AS orderNumber,
  o.requestNumber AS requestNumber,
  o.contractNumber AS contractNumber,
  o.orderDate AS orderDate,
  s.name AS supplierName,
  os.name AS orderStateName,
  lastUpdate.lastUpdateTime as lastUpdateDate,
  lastUpdate.userName as lastUpdatedBy
`
	result.ReturnAlias = "orders"
	result.Parameters = make(map[string]interface{})
	result.Parameters["search"] = searchString
	result.Parameters["limit"] = pagination.PageSize
	result.Parameters["skip"] = (pagination.Page - 1) * pagination.PageSize
	return result
}
