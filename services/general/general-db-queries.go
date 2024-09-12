package general

import (
	"panda/apigateway/helpers"
)

func GetGraphNodesByUidQuery(uid string) (result helpers.DatabaseQuery) {
	result.Query = `match(n{uid:$uid})-[r]-(o) RETURN DISTINCT {
	 uid: o.uid, 
	 name: o.name, 
	 label: labels(o)[0], 
	 properties: apoc.map.removeKeys(properties(o), ['passwordHash','passwordToChange', 'isEnabled', 'deleted', 'username', 'printEUN']) } as nodes`
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid
	result.ReturnAlias = "nodes"

	return result
}

func GetGraphLinksByUidQuery(uid string) (result helpers.DatabaseQuery) {
	result.Query = `
	match(n{uid:$uid})-[r]->(o) RETURN DISTINCT { source: n.uid, target: o.uid, relationship: type(r) } as links
	union all
	match(n{uid:$uid})<-[r]-(o) RETURN DISTINCT { source: o.uid, target: n.uid, relationship: type(r) } as links`
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid
	result.ReturnAlias = "links"

	return result
}
