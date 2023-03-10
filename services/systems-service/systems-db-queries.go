package systemsService

import (
	"panda/apigateway/helpers"
	"panda/apigateway/services/systems-service/models"
	"strings"

	"github.com/google/uuid"
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

func GetLocationsBySearchTextQuery(searchText string, limit int, facilityCode string) (result helpers.DatabaseQuery) {
	searchText = strings.ToLower(searchText)
	result.Query = `
	MATCH (n:Location)-[:BELONGS_TO_FACILITY]->(f) where f.code = $facilityCode and n.code is not null and not (n)-[:HAS_SUBLOCATION]->()
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
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func GetZonesCodebookQuery() (result helpers.DatabaseQuery) {
	result.Query = `MATCH(r:Zone) WHERE NOT ()-[:HAS_SUBZONE]->(r) RETURN {uid: r.uid,name: r.code + " - " + r.name} as zones ORDER BY zones.name`
	result.ReturnAlias = "zones"
	result.Parameters = make(map[string]interface{})
	return result
}

func GetSubZonesByParentUidCodebookQuery(parentUID string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH(parent:Zone{uid:$parentUID})-[:HAS_SUBZONE]->(r) RETURN {uid: r.uid,name:r.code + " - " + r.name} as zones ORDER BY zones.name`
	result.ReturnAlias = "zones"
	result.Parameters = make(map[string]interface{})
	result.Parameters["parentUID"] = parentUID
	return result
}

func GetSubSystemsQuery(parentUID string, facilityCode string) (result helpers.DatabaseQuery) {

	//we have to diff queries depend if it is or not a root parent
	if parentUID != "" {
		result.Query = `
		MATCH(r:System)-[:BELONGS_TO_FACILITY]->(f) WHERE f.code = $facilityCode WITH r			
		MATCH (parent)-[:HAS_SUBSYSTEM]->(r)
		where parent.uid = $parentUID
		return {uid: r.uid, name: r.name} as result;`
	} else {
		result.Query = `
		MATCH(r:System)-[:BELONGS_TO_FACILITY]->(f)			
		where not ()-[:HAS_SUBSYSTEM]->(r) and f.code = $facilityCode
		return {uid: r.uid, name: r.name} as result;`
	}

	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["parentUID"] = parentUID
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func SystemImageByUidQuery(uid string) (result helpers.DatabaseQuery) {
	result.Query = `match(r:System{uid: $uid})	
	return r.image as image`
	result.ReturnAlias = "image"
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid

	return result
}

func SystemDetailQuery(uid string, facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH(r:System{uid: $uid})-[:BELONGS_TO_FACILITY]->(f) WHERE f.code = $facilityCode
	WITH r,f
OPTIONAL MATCH(r)-[:HAS_LOCATION]->(l)
OPTIONAL MATCH(r)-[:HAS_ZONE]->(z)
OPTIONAL MATCH(r)-[:HAS_SYSTEM_TYPE]->(st)
OPTIONAL MATCH(r)-[:HAS_IMPORTANCE]->(imp)
OPTIONAL MATCH(r)-[:HAS_OWNER]->(own)
OPTIONAL MATCH(r)-[:HAS_CRITICALITY]->(cc)
OPTIONAL MATCH(r)-[:CONTAINS_ITEM]->(itm)
OPTIONAL MATCH(parent)-[:HAS_SUBSYSTEM*1..50]->(r)
WITH r,l, z, st,itm, imp, own,cc, case when parent is not null then collect({uid: parent.uid, name: parent.name}) else null end as parents
WITH r,l, z, st,itm, imp, own,cc, reverse(parents) as parents
RETURN {
    uid: r.uid, 
    name: r.name, 
    description: r.description,
    location: case when l is not null then {uid: l.code, name: l.name} else null end, 
    systemType: case when st is not null then {uid: st.uid, name: st.name} else null end,
    systemCode: r.systemCode,
    systemALias: r.systemAlias,
    importance: case when imp is not null then {uid: imp.uid, name: imp.name} else null end,
    owner: case when own is not null then {uid: own.uid, name: own.lastName + " " + own.firstName} else null end,
    zone: case when z is not null then {uid: z.uid, name: z.name} else null end,
    parentPath: parents,
	itemUID: itm.uid    
    } as result`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func CreateNewSystemQuery(newSystem *models.SystemForm, facilityCode string) (result helpers.DatabaseQuery) {
	result.Parameters = make(map[string]interface{})
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["uid"] = uuid.NewString()
	result.Parameters["name"] = newSystem.Name
	result.Parameters["description"] = newSystem.Description
	result.Parameters["systemCode"] = newSystem.SystemCode
	result.Parameters["systemAlias"] = newSystem.SystemAlias

	result.Query = `
	CREATE(s:System{uid: $uid})
	SET 
	s.name = $name, 
	s.description = $description,
	s.systemCode = $systemCode,
	s.systemAlias = $systemAlias
	WITH s
	MATCH(f:Facility{code: $facilityCode})
	CREATE(s)-[:BELONGS_TO_FACILITY]->(f)
	WITH s
	`

	if newSystem.ParentUID != nil {
		result.Query += `WITH s MATCH(parent:System{uid:$parentUID}) MERGE(parent)-[:HAS_SUBSYSTEM]->(s) `
		result.Parameters["parentUID"] = newSystem.ParentUID
	}

	if newSystem.ZoneUID != nil {
		result.Query += `WITH s MATCH(z:Zone{uid:$zoneUID}) MERGE(s)-[:HAS_ZONE]->(z) `
		result.Parameters["zoneUID"] = newSystem.ZoneUID
	}

	if newSystem.LocationUID != nil {
		result.Query += `WITH s MATCH(l:Location{code:$locationUID})-[:BELONGS_TO_FACILITY]->(f{code:$facilityCode}) MERGE(s)-[:HAS_LOCATION]->(l) `
		result.Parameters["locationUID"] = newSystem.LocationUID
	}

	if newSystem.SystemTypeUID != nil {
		result.Query += `WITH s MATCH(st:SystemType{uid:$systemTypeUID}) MERGE(s)-[:HAS_SYSTEM_TYPE]->(st) `
		result.Parameters["systemTypeUID"] = newSystem.SystemTypeUID
	}

	if newSystem.OwnerUID != nil {
		result.Query += `WITH s MATCH(owner:User{uid:$ownerUID}) MERGE(s)-[:HAS_OWNER]->(owner) `
		result.Parameters["ownerUID"] = newSystem.OwnerUID
	}

	if newSystem.ImportanceUID != nil {
		result.Query += `WITH s MATCH(imp:SystemImportance{uid:$importanceUID}) MERGE(s)-[:HAS_IMPORTANCE]->(imp) `
		result.Parameters["importanceUID"] = newSystem.ImportanceUID
	}

	if newSystem.Image != nil {
		result.Query += `WITH s SET s.image = $image `
		result.Parameters["image"] = newSystem.Image
	}

	result.Query += `RETURN s.uid as result`

	result.ReturnAlias = "result"

	return result
}
