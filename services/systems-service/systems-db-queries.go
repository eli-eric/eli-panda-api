package systemsService

import (
	"panda/apigateway/helpers"
	"strings"
)

func GetSystemTypesCodebookQuery() (result helpers.DatabaseQuery) {
	result.Query = `MATCH (n:SystemTypeGroup)-[:CONTAINS_SYSTEM_TYPE]->(st) with n, st order by st.name 
					return {uid:st.uid, name: n.name+ " > "+ st.name} as result order by result.name`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	return result
}

func GetSystemImportancesCodebookQuery() (result helpers.DatabaseQuery) {
	result.Query = `MATCH(r:SystemImportance) RETURN {uid: r.uid,name:r.name} as result ORDER BY result.code`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	return result
}

func GetSystemCriticalityCodebookQuery() (result helpers.DatabaseQuery) {
	result.Query = `MATCH(r:SystemCriticality) RETURN {uid: r.uid,name:r.name} as result ORDER BY result.code`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	return result
}

func GetItemUsagesCodebookQuery() (result helpers.DatabaseQuery) {
	result.Query = `MATCH(r:ItemUsage) RETURN {uid: r.uid,name:r.name} as result ORDER BY result.code`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	return result
}

func GetItemConditionsCodebookQuery() (result helpers.DatabaseQuery) {
	result.Query = `MATCH(r:ItemCondition) RETURN {uid: r.uid,name:r.name} as result ORDER BY result.code`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	return result
}

func GetLocationsBySearchTextQuery(searchText string, limit int) (result helpers.DatabaseQuery) {
	searchText = strings.ToLower(searchText)
	result.Query = `
	MATCH (n:Location) where n.code is not null and not (n)-[:HAS_SUBLOCATION]->()
	with n 
	where (toLower(n.code) contains $searchText or toLower(n.name) contains $searchText) 
	optional match (parent)-[:HAS_SUBLOCATION*1..50]->(n) 
	with n, collect(parent.name) as parentNames
	return {uid: n.code, name: n.code + " - " +  n.name + " - " + apoc.text.join(reverse(parentNames), " > ")} as result
	order by result.name 
	limit $limit`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["searchText"] = searchText
	result.Parameters["limit"] = limit
	return result
}

func GetZonesCodebookQuery() (result helpers.DatabaseQuery) {
	result.Query = `MATCH(r:Zone) WHERE NOT ()-[:HAS_SUBZONE]->(r) RETURN {uid: r.uid,name:r.name} as zones ORDER BY zones.name`
	result.ReturnAlias = "zones"
	result.Parameters = make(map[string]interface{})
	return result
}

func GetSubZonesByParentUidCodebookQuery(parentUID string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH(parent:Zone{uid:$parentUID})-[:HAS_SUBZONE]->(r) RETURN {uid: r.uid,name:r.name} as zones ORDER BY zones.name`
	result.ReturnAlias = "zones"
	result.Parameters = make(map[string]interface{})
	result.Parameters["parentUID"] = parentUID
	return result
}
