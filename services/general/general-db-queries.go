package general

import (
	"panda/apigateway/helpers"
	"strings"
)

func GetGraphNodesByUidQuery(uid string) (result helpers.DatabaseQuery) {
	result.Query = `match(n{uid:$uid})-[r]-(o) WHERE o.uid IS NOT NULL RETURN DISTINCT {
	 uid: o.uid, 
	 name: o.name, 
	 label: labels(o)[0], 
	 properties: apoc.map.removeKeys(properties(o), ['passwordHash','passwordToChange', 'isEnabled', 'deleted', 'username', 'printEUN', 'image']) } as nodes
	 union all
     match(n{uid:$uid}) RETURN DISTINCT {
	 uid: n.uid, 
	 name: n.name, 
	 label: labels(n)[0], 
	 properties: apoc.map.removeKeys(properties(n), ['passwordHash','passwordToChange', 'isEnabled', 'deleted', 'username', 'printEUN', 'image'] )} as nodes`
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid
	result.ReturnAlias = "nodes"

	return result
}

func GetGraphLinksByUidQuery(uid string) (result helpers.DatabaseQuery) {
	result.Query = `
	match(n{uid:$uid})-[r]->(o) WHERE o.uid IS NOT NULL RETURN DISTINCT { source: n.uid, target: o.uid, relationship: type(r) } as links
	union all
	match(n{uid:$uid})<-[r]-(o) WHERE o.uid IS NOT NULL RETURN DISTINCT { source: o.uid, target: n.uid, relationship: type(r) } as links`
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid
	result.ReturnAlias = "links"

	return result
}

func GetUUIDQuery() (result helpers.DatabaseQuery) {
	result.Query = `RETURN apoc.create.uuid() as uuid`
	result.Parameters = make(map[string]interface{})
	result.ReturnAlias = "uuid"

	return result
}

func GetGlobalSearchQuery(searchText string, skip int, limit int) (result helpers.DatabaseQuery) {
	result.Query = `
		// Search in Systems
		MATCH (n:System)
		WHERE (toLower(n.name) CONTAINS $searchText 
			OR toLower(n.systemCode) CONTAINS $searchText
			OR toLower(n.description) CONTAINS $searchText)
			AND (n.deleted IS NULL OR n.deleted <> true)
		WITH n, 'System' as nodeType
		RETURN {
			uid: n.uid, 
			name: n.name, 
			description: COALESCE(n.description, ''), 
			nodeType: nodeType
		} as result
		
		UNION ALL
		
		// Search in Orders
		MATCH (n:Order)
		WHERE (toLower(n.name) CONTAINS $searchText 
			OR toLower(n.orderNumber) CONTAINS $searchText
			OR toLower(n.contractNumber) CONTAINS $searchText
			OR toLower(n.requestNumber) CONTAINS $searchText)
			AND (n.deleted IS NULL OR n.deleted <> true)
		WITH n, 'Order' as nodeType
		RETURN {
			uid: n.uid, 
			name: n.name, 
			description: COALESCE(n.notes, ''), 
			nodeType: nodeType
		} as result
		
		UNION ALL
		
		// Search in CatalogueItems
		MATCH (n:CatalogueItem)
		WHERE (toLower(n.name) CONTAINS $searchText 
			OR toLower(n.catalogueNumber) CONTAINS $searchText
			OR toLower(n.description) CONTAINS $searchText)
			AND (n.deleted IS NULL OR n.deleted <> true)
		WITH n, 'CatalogueItem' as nodeType
		RETURN {
			uid: n.uid, 
			name: n.name, 
			description: COALESCE(n.description, ''), 
			nodeType: nodeType
		} as result
		
		// Add Items that have relations to Systems or Orders
		UNION ALL
		
		MATCH (n:Item)<-[:CONTAINS_ITEM]-(s:System)
		WHERE (toLower(n.name) CONTAINS $searchText 
			OR toLower(n.eun) CONTAINS $searchText
			OR toLower(n.serialNumber) CONTAINS $searchText)
			AND (n.deleted IS NULL OR n.deleted <> true)
		WITH DISTINCT s, 'System' as nodeType
		RETURN {
			uid: s.uid, 
			name: s.name, 
			description: COALESCE(s.description, ''), 
			nodeType: nodeType
		} as result
		
		UNION ALL
		
		MATCH (n:Item)<-[:HAS_ORDER_LINE]-(o:Order)
		WHERE (toLower(n.name) CONTAINS $searchText 
			OR toLower(n.eun) CONTAINS $searchText
			OR toLower(n.serialNumber) CONTAINS $searchText)
			AND (n.deleted IS NULL OR n.deleted <> true)
		WITH DISTINCT o, 'Order' as nodeType
		RETURN {
			uid: o.uid, 
			name: o.name, 
			description: COALESCE(o.notes, ''), 
			nodeType: nodeType
		} as result
		
		ORDER BY result.name
		SKIP $skip
		LIMIT $limit
	`

	result.Parameters = make(map[string]interface{})
	result.Parameters["searchText"] = strings.ToLower(searchText)
	result.Parameters["skip"] = skip
	result.Parameters["limit"] = limit

	result.ReturnAlias = "result"

	return result
}

func GetGlobalSearchCountQuery(searchText string) (result helpers.DatabaseQuery) {
	result.Query = `
		CALL {
			// Count Systems
			MATCH (n:System)
			WHERE (toLower(n.name) CONTAINS $searchText 
				OR toLower(n.systemCode) CONTAINS $searchText
				OR toLower(n.description) CONTAINS $searchText)
				AND (n.deleted IS NULL OR n.deleted <> true)
			RETURN count(n) as systemCount
		}
		CALL {
			// Count Orders
			MATCH (n:Order)
			WHERE (toLower(n.name) CONTAINS $searchText 
				OR toLower(n.orderNumber) CONTAINS $searchText
				OR toLower(n.contractNumber) CONTAINS $searchText
				OR toLower(n.requestNumber) CONTAINS $searchText)
				AND (n.deleted IS NULL OR n.deleted <> true)
			RETURN count(n) as orderCount
		}
		CALL {
			// Count CatalogueItems
			MATCH (n:CatalogueItem)
			WHERE (toLower(n.name) CONTAINS $searchText 
				OR toLower(n.catalogueNumber) CONTAINS $searchText
				OR toLower(n.description) CONTAINS $searchText)
				AND (n.deleted IS NULL OR n.deleted <> true)
			RETURN count(n) as catalogueCount
		}
		CALL {
			// Count distinct Systems from Items
			MATCH (n:Item)<-[:CONTAINS_ITEM]-(s:System)
			WHERE (toLower(n.name) CONTAINS $searchText 
				OR toLower(n.eun) CONTAINS $searchText
				OR toLower(n.serialNumber) CONTAINS $searchText)
				AND (n.deleted IS NULL OR n.deleted <> true)
			RETURN count(DISTINCT s) as itemSystemCount
		}
		CALL {
			// Count distinct Orders from Items
			MATCH (n:Item)<-[:HAS_ORDER_LINE]-(o:Order)
			WHERE (toLower(n.name) CONTAINS $searchText 
				OR toLower(n.eun) CONTAINS $searchText
				OR toLower(n.serialNumber) CONTAINS $searchText)
				AND (n.deleted IS NULL OR n.deleted <> true)
			RETURN count(DISTINCT o) as itemOrderCount
		}
		
		RETURN (systemCount + orderCount + catalogueCount + itemSystemCount + itemOrderCount) as totalCount
	`

	result.Parameters = make(map[string]interface{})
	result.Parameters["searchText"] = strings.ToLower(searchText)

	result.ReturnAlias = "totalCount"

	return result
}
