package systemsService

import (
	"panda/apigateway/helpers"
)


func GetSystemTypesCodebookQuery() (result helpers.DatabaseQuery) {
	result.Query = `MATCH (n:SystemTypeGroup)-[:CONTAINS_SYSTEM_TYPE]->(st) with n, st order by st.name 
					return {uid:st.uid, name: n.name+ " > "+ st.name} as result order by result.name`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	return result
}
