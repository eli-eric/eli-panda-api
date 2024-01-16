package systemsService

import (
	"fmt"
	"panda/apigateway/helpers"
	"panda/apigateway/services/systems-service/models"
	"strings"

	"github.com/google/uuid"
)

func GetSystemTypesCodebookQuery(facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = fmt.Sprintf(`MATCH (n:SystemTypeGroup)-[:CONTAINS_SYSTEM_TYPE]->(st) with n, st order by st.name 
	return {uid:st.uid, name: n.name+ " > "+ st.name + case when st.mask%v is null then "" else  " (" + st.mask%v  + ")" end } as result order by result.name`, facilityCode, facilityCode)
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["facilityCode"] = facilityCode
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
	return {uid: n.uid, name: n.code + " - " +  n.name + " - " + apoc.text.join(reverse(parentNames), " > ")} as result
	order by result.name 
	limit $limit`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["searchText"] = searchText
	result.Parameters["limit"] = limit
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func GetZonesCodebookQuery(facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH(f:Facility{code:$facilityCode}) WITH f
	MATCH(z:Zone)-[:HAS_SUBZONE]->(sz)-[:BELONGS_TO_FACILITY]->(f) return {uid:sz.uid, name: z.code+"-"+sz.code + " - " + sz.name + " ("+  z.name + ")"} as zone order by z.code, sz.code
		UNION
		MATCH(f:Facility{code:$facilityCode}) WITH f
		WITH f
		MATCH(z:Zone)-[:BELONGS_TO_FACILITY]->(f) where not ()-[:HAS_SUBZONE]->(z)  return {uid:z.uid, name:z.code + " - " +z.name } as zone order by z.code`
	result.ReturnAlias = "zone"
	result.Parameters = make(map[string]interface{})
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

func CreateNewSystemQuery(newSystem *models.System, facilityCode string, userUID string) (result helpers.DatabaseQuery) {
	result.Parameters = make(map[string]interface{})
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["uid"] = uuid.NewString()
	result.Parameters["name"] = newSystem.Name
	result.Parameters["description"] = newSystem.Description
	result.Parameters["systemCode"] = newSystem.SystemCode
	result.Parameters["systemAlias"] = newSystem.SystemAlias
	result.Parameters["lastUpdateBy"] = userUID

	result.Query = `
	MATCH(f:Facility{code: $facilityCode}) WITH f	
	MATCH(u:User{uid: $lastUpdateBy}) WITH u, f
	CREATE(s:System{uid: $uid, deleted: false, lastUpdateTime: datetime(), lastUpdatedBy: u.lastName + " " + u.firstName})-[:BELONGS_TO_FACILITY]->(f)
	SET 
	s.name = $name, 
	s.description = $description,
	s.systemCode = $systemCode,
	s.systemAlias = $systemAlias
	WITH s, u
	CREATE(s)-[:WAS_UPDATED_BY{ at: datetime(), action: "INSERT" }]->(u)	
	WITH s
	`

	if newSystem.Zone != nil && newSystem.Zone.UID != "" {
		result.Query += `WITH s MATCH(z:Zone{uid:$zoneUID}) MERGE(s)-[:HAS_ZONE]->(z) `
		result.Parameters["zoneUID"] = newSystem.Zone.UID
	}

	if newSystem.Location != nil && newSystem.Location.UID != "" {
		result.Query += `WITH s MATCH(l:Location{code:$locationUID})-[:BELONGS_TO_FACILITY]->(f{code:$facilityCode}) MERGE(s)-[:HAS_LOCATION]->(l) `
		result.Parameters["locationUID"] = newSystem.Location.UID
	}

	if newSystem.SystemType != nil && newSystem.SystemType.UID != "" {
		result.Query += `WITH s MATCH(st:SystemType{uid:$systemTypeUID}) MERGE(s)-[:HAS_SYSTEM_TYPE]->(st) `
		result.Parameters["systemTypeUID"] = newSystem.SystemType.UID
	}

	if newSystem.Owner != nil && newSystem.Owner.UID != "" {
		result.Query += `WITH s MATCH(owner:Employee{uid:$ownerUID}) MERGE(s)-[:HAS_OWNER]->(owner) `
		result.Parameters["ownerUID"] = newSystem.Owner.UID
	}

	if newSystem.Responsible != nil && newSystem.Responsible.UID != "" {
		result.Query += `WITH s MATCH(responsible:Employee{uid:$responsibleUID}) MERGE(s)-[:HAS_RESPONSIBLE]->(responsible) `
		result.Parameters["responsibleUID"] = newSystem.Responsible.UID
	}

	if newSystem.Importance != nil && newSystem.Importance.UID != "" {
		result.Query += `WITH s MATCH(imp:SystemImportance{uid:$importanceUID}) MERGE(s)-[:HAS_IMPORTANCE]->(imp) `
		result.Parameters["importanceUID"] = newSystem.Importance.UID
	}

	if newSystem.ParentUID != nil && *newSystem.ParentUID != "" {
		result.Query += `WITH s MATCH(parent:System{uid:$parentUID}) MERGE(parent)-[:HAS_SUBSYSTEM]->(s) `
		result.Parameters["parentUID"] = *newSystem.ParentUID
	}

	if newSystem.PhysicalItem != nil && newSystem.PhysicalItem.UID != "" {
		//unassign from previous system
		result.Query += `WITH s MATCH(prevSystem)-[rpiold:CONTAINS_ITEM]->(pi:Item{uid:$physicalItemUID})-[:IS_BASED_ON]->(ci) DELETE rpiold `

		result.Query += `WITH s, pi, ci MERGE(s)-[:CONTAINS_ITEM]->(pi) `
		result.Query += `
		WITH s, pi, ci
		SET pi.lastUpdateTime = datetime(), 
		pi.lastUpdatedBy = s.lastUpdatedBy,
		pi.serialNUmber = $serialNumber,
		pi.price = $price,
		pi.currency = $currency,
		ci.name = $catalogueName,
		ci.description = $catalogueDescription,
		ci.catalogueNumber = $catalogueNumber `

		result.Parameters["physicalItemUID"] = newSystem.PhysicalItem.UID
		result.Parameters["serialNumber"] = newSystem.PhysicalItem.SerialNumber
		result.Parameters["price"] = newSystem.PhysicalItem.Price
		result.Parameters["currency"] = newSystem.PhysicalItem.Currency
		result.Parameters["catalogueDescription"] = newSystem.PhysicalItem.CatalogueItem.Description
		result.Parameters["catalogueName"] = newSystem.PhysicalItem.CatalogueItem.Name
		result.Parameters["catalogueNumber"] = newSystem.PhysicalItem.CatalogueItem.CatalogueNumber

		if newSystem.PhysicalItem.ItemUsage != nil && newSystem.PhysicalItem.ItemUsage.UID != "" {
			result.Query += `WITH s, pi, ci OPTIONAL MATCH(pi)-[rpiUsage:HAS_ITEM_USAGE]->() DELETE rpiUsage `
			result.Query += `WITH s, pi, ci MATCH(piUsage:ItemUsage{uid:$itemUsageUID}) MERGE(pi)-[:HAS_ITEM_USAGE]->(piUsage) `
			result.Parameters["itemUsageUID"] = newSystem.PhysicalItem.ItemUsage.UID
		} else {
			result.Query += `WITH s, pi, ci OPTIONAL MATCH(pi)-[rpiUsage:HAS_ITEM_USAGE]->() DELETE rpiUsage `
		}

	}

	result.Query += `RETURN s.uid as result`

	result.ReturnAlias = "result"

	return result
}

func UpdateSystemQuery(newSystem *models.System, oldSystem *models.System, facilityCode string, userUID string) (result helpers.DatabaseQuery) {
	result.Parameters = make(map[string]interface{})
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["uid"] = oldSystem.UID

	result.Query = `MATCH(s:System{uid:$uid, deleted: false})-[:BELONGS_TO_FACILITY]->(f:Facility{code:$facilityCode}) `

	if newSystem.ParentUID != nil && *newSystem.ParentUID != "" {
		result.Query += `WITH s OPTIONAL MATCH(parent)-[oldParent:HAS_SUBSYSTEM]->(s) DELETE oldParent `
		result.Query += `WITH s MATCH(parent:System{uid:$parentUID}) MERGE(parent)-[:HAS_SUBSYSTEM]->(s) `
		result.Parameters["parentUID"] = *newSystem.ParentUID
	}

	if newSystem.Location != nil && newSystem.Location.UID != "" {
		result.Query += `WITH s OPTIONAL MATCH(s)-[rl:HAS_LOCATION]->() DELETE rl `
		result.Query += `WITH s MATCH(l:Location{code:$locationUID})-[:BELONGS_TO_FACILITY]->(f{code:$facilityCode}) MERGE(s)-[:HAS_LOCATION]->(l) `
		result.Parameters["locationUID"] = newSystem.Location.UID
	} else {
		result.Query += `WITH s OPTIONAL MATCH(s)-[rl:HAS_LOCATION]->() DELETE rl `
	}

	helpers.AutoResolveObjectToUpdateQuery(&result, *newSystem, *oldSystem, "s")

	if newSystem.PhysicalItem != nil && newSystem.PhysicalItem.UID != "" {
		//unassign from previous system if its another system
		if oldSystem.PhysicalItem != nil && oldSystem.PhysicalItem.UID != newSystem.PhysicalItem.UID {
			result.Query += `WITH s MATCH(prevSystem)-[rpiold:CONTAINS_ITEM]->(pi:Item{uid:$physicalItemUID})-[:IS_BASED_ON]->(ci) DELETE rpiold `
			result.Query += `WITH s, pi, ci MERGE(s)-[:CONTAINS_ITEM]->(pi) `
		} else if oldSystem.PhysicalItem != nil && oldSystem.PhysicalItem.UID == newSystem.PhysicalItem.UID {
			result.Query += `WITH s MATCH(s)-[:CONTAINS_ITEM]->(pi:Item{uid:$physicalItemUID})-[:IS_BASED_ON]->(ci) `
		}

		result.Query += `
		WITH s, pi, ci
		SET pi.lastUpdateTime = datetime(), 
		pi.lastUpdatedBy = s.lastUpdatedBy,
		pi.serialNUmber = $serialNumber,
		pi.price = $price,
		pi.currency = $currency,
		ci.name = $catalogueName,
		ci.description = $catalogueDescription,
		ci.catalogueNumber = $catalogueNumber `

		result.Parameters["physicalItemUID"] = newSystem.PhysicalItem.UID
		result.Parameters["serialNumber"] = newSystem.PhysicalItem.SerialNumber
		result.Parameters["price"] = newSystem.PhysicalItem.Price
		result.Parameters["currency"] = newSystem.PhysicalItem.Currency
		result.Parameters["catalogueDescription"] = newSystem.PhysicalItem.CatalogueItem.Description
		result.Parameters["catalogueName"] = newSystem.PhysicalItem.CatalogueItem.Name
		result.Parameters["catalogueNumber"] = newSystem.PhysicalItem.CatalogueItem.CatalogueNumber

		if newSystem.PhysicalItem.ItemUsage != nil && newSystem.PhysicalItem.ItemUsage.UID != "" {
			result.Query += `WITH s, pi, ci OPTIONAL MATCH(pi)-[rpiUsage:HAS_ITEM_USAGE]->() DELETE rpiUsage `
			result.Query += `WITH s, pi, ci MATCH(piUsage:ItemUsage{uid:$itemUsageUID}) MERGE(pi)-[:HAS_ITEM_USAGE]->(piUsage) `
			result.Parameters["itemUsageUID"] = newSystem.PhysicalItem.ItemUsage.UID
		} else {
			result.Query += `WITH s, pi, ci OPTIONAL MATCH(pi)-[rpiUsage:HAS_ITEM_USAGE]->() DELETE rpiUsage `
		}

	}

	result.Query += `RETURN s.uid as result`

	result.ReturnAlias = "result"

	return result
}

func DeleteSystemByUidQuery(uid string) (result helpers.DatabaseQuery) {

	result.Query = `MATCH(r:System) WHERE r.uid = $uid WITH r
	OPTIONAL MATCH (r)-[:HAS_SUBSYSTEM*1..50]->(child)
	WITH r, child, r.uid as uid
	SET r.deleted=true, child.deleted=true`

	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid

	return result
}

func GetSystemsForAutocomplete(search string, limit int, facilityCode string, onlyTechnologicalUnits bool) (result helpers.DatabaseQuery) {

	if onlyTechnologicalUnits {
		result.Query = `
	MATCH (n:System{systemLevel: 'TECHNOLOGY_UNIT', deleted: false})-[:BELONGS_TO_FACILITY]->(f)
	WHERE f.code = $facilityCode AND NOT (n)-[:HAS_SUBSYSTEM]->(:System{systemLevel: 'TECHNOLOGY_UNIT', deleted: false})
	WITH n
	OPTIONAL MATCH (parent{deleted: false})-[:HAS_SUBSYSTEM*1..50]->(n{systemLevel: 'TECHNOLOGY_UNIT', deleted: false})
	WITH n, collect(parent.name) AS parentNames
	WITH {uid: n.uid, name: n.name + " < " + apoc.text.join((parentNames), " < ")} AS result
	WHERE toLower(result.name) CONTAINS $searchText
	RETURN result
	ORDER BY result.name
	LIMIT $limit`

	} else {
		result.Query = `
	MATCH (n:System{deleted: false})-[:BELONGS_TO_FACILITY]->(f)
	WHERE f.code = $facilityCode AND NOT (n)-[:HAS_SUBSYSTEM]->()
	WITH n	
	OPTIONAL MATCH (parent{deleted: false})-[:HAS_SUBSYSTEM*1..50]->(n)
	WITH n, collect(parent.name) AS parentNames
	WITH {uid: n.uid, name: n.name + " < " + apoc.text.join((parentNames), " < ")} AS result
	WHERE toLower(result.name) CONTAINS $searchText
	RETURN result
	ORDER BY result.name
	LIMIT $limit`

	}

	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})

	result.Parameters["searchText"] = strings.ToLower(search)
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["limit"] = limit
	return result
}

func GetSystemsBySearchTextFullTextQuery(searchString string, facilityCode string, pagination *helpers.Pagination, sorting *[]helpers.Sorting) (result helpers.DatabaseQuery) {

	result.Parameters = make(map[string]interface{})

	if searchString == "" {
		result.Query = "MATCH(f:Facility{code: $facilityCode}) WITH f MATCH(sys:System{deleted:false})-[:BELONGS_TO_FACILITY]->(f) WHERE NOT ()-[:HAS_SUBSYSTEM]->(sys) WITH sys "
	} else {
		result.Query = `
		CALL {
			CALL db.index.fulltext.queryNodes('searchIndexSystems', $fulltextSearch) YIELD node AS sys WHERE sys:System AND sys.deleted = false return sys
			UNION
			MATCH (sys{deleted:false})-[:CONTAINS_ITEM]->(physicalItem)-[:IS_BASED_ON]->(catalogueItem) 
			WHERE toLower(physicalItem.eun) STARTS WITH $search OR toLower(catalogueItem.catalogueNumber) STARTS WITH $search OR toLower(catalogueItem.name) STARTS WITH $search
			RETURN sys
		}
		WITH sys
		MATCH(f:Facility{code: $facilityCode}) WITH f, sys
		MATCH(sys)-[:BELONGS_TO_FACILITY]->(f)
		WITH sys `
	}

	result.Query += `	
	OPTIONAL MATCH (sys)-[:HAS_LOCATION]->(loc)  
	OPTIONAL MATCH (sys)-[:HAS_ZONE]->(zone)  
	OPTIONAL MATCH (sys)-[:HAS_SYSTEM_TYPE]->(st)	
	OPTIONAL MATCH (sys)-[:HAS_OWNER]->(own)
	OPTIONAL MATCH (sys)-[:HAS_RESPONSIBLE]->(responsilbe)
	OPTIONAL MATCH (sys)-[:HAS_IMPORTANCE]->(imp)
	OPTIONAL MATCH (sys)-[:CONTAINS_ITEM]->(physicalItem)-[:IS_BASED_ON]->(catalogueItem)-[:BELONGS_TO_CATEGORY]->(ciCategory)	
	OPTIONAL MATCH (physicalItem)-[:HAS_ITEM_USAGE]->(itemUsage)
	OPTIONAL MATCH (parents)-[:HAS_SUBSYSTEM*1..50]->(sys)
	OPTIONAL MATCH (sys)-[:HAS_SUBSYSTEM*1..50]->(subsys)
	RETURN DISTINCT {  
	uid: sys.uid,
	description: sys.description,
	name: sys.name,
	parentPath: case when parents is not null then reverse(collect(distinct {uid: parents.uid, name: parents.name})) else null end,
	hasSubsystems: case when subsys is not null then true else false end,
	systemCode: sys.systemCode,
	systemAlias: sys.systemAlias,
	systemLevel: sys.systemLevel,
	isTechnologicalUnit: sys.isTechnologicalUnit,
	location: case when loc is not null then {uid: loc.code, name: loc.name} else null end,
	zone: case when zone is not null then {uid: zone.uid, name: zone.name} else null end,
	systemType: case when st is not null then {uid: st.uid, name: st.name} else null end,
	owner: case when own is not null then {uid: own.uid, name: own.lastName + " " + own.firstName} else null end,
	responsible: case when responsilbe is not null then {uid: responsilbe.uid, name: responsilbe.lastName + " " + responsilbe.firstName} else null end,
	importance: case when imp is not null then {uid: imp.uid, name: imp.name} else null end,	
	lastUpdateTime: sys.lastUpdateTime,
	lastUpdateBy: sys.lastUpdateBy,
	physicalItem: case when physicalItem is not null then {
		uid: physicalItem.uid, 
		eun: physicalItem.eun, 
		serialNumber: physicalItem.serialNumber,
		price: physicalItem.price,
		currency: physicalItem.currency,
		itemUsage: case when itemUsage is not null then {uid: itemUsage.uid, name: itemUsage.name} else null end,
		catalogueItem: case when catalogueItem is not null then {
			uid: catalogueItem.uid,
			name: catalogueItem.name,
			catalogueNumber: catalogueItem.catalogueNumber,
			category: case when ciCategory is not null then {uid: ciCategory.uid, name: ciCategory.name} else null end
			} else null end	
	} else null end,
	statistics: {subsystemsCount: count(subsys)}
	} AS systems

` + GetSystemsOrderByClauses(sorting) + `

	SKIP $skip
	LIMIT $limit

`
	result.ReturnAlias = "systems"

	result.Parameters["search"] = strings.ToLower(searchString)
	result.Parameters["fulltextSearch"] = strings.ToLower(helpers.GetFullTextSearchString(searchString))
	result.Parameters["limit"] = pagination.PageSize
	result.Parameters["skip"] = (pagination.Page - 1) * pagination.PageSize
	result.Parameters["facilityCode"] = facilityCode

	return result
}

func GetSystemsOrderByClauses(sorting *[]helpers.Sorting) string {

	if sorting == nil || len(*sorting) == 0 {
		return `ORDER BY systems.lastUpdateTime DESC `
	}

	var result string = ` ORDER BY `

	for i, sort := range *sorting {
		if i > 0 {
			result += ", "
		}
		result += "systems." + sort.ID
		if sort.DESC {
			result += " DESC "
		}
	}

	return result
}

func GetSystemsBySearchTextFullTextCountQuery(searchString string, facilityCode string) (result helpers.DatabaseQuery) {

	result.Parameters = make(map[string]interface{})

	if searchString == "" {
		result.Query = "MATCH(f:Facility{code: $facilityCode}) WITH f MATCH(sys:System{deleted:false})-[:BELONGS_TO_FACILITY]->(f) WHERE NOT ()-[:HAS_SUBSYSTEM]->(sys) WITH sys "
	} else {
		result.Query = `
		CALL {
			CALL db.index.fulltext.queryNodes('searchIndexSystems', $fulltextSearch) YIELD node AS sys WHERE sys:System AND sys.deleted = false return sys
			UNION
			MATCH (sys{deleted:false})-[:CONTAINS_ITEM]->(physicalItem)-[:IS_BASED_ON]->(catalogueItem) 
			WHERE toLower(physicalItem.eun) STARTS WITH $search OR toLower(catalogueItem.catalogueNumber) STARTS WITH $search OR toLower(catalogueItem.name) STARTS WITH $search
			RETURN sys
		}
		WITH sys
		MATCH(f:Facility{code: $facilityCode}) WITH f, sys
		MATCH(sys)-[:BELONGS_TO_FACILITY]->(f)
		WITH sys `
	}

	result.Query += `	
	OPTIONAL MATCH (sys)-[:HAS_LOCATION]->(loc)  
	OPTIONAL MATCH (sys)-[:HAS_ZONE]->(zone)  
	OPTIONAL MATCH (sys)-[:HAS_SYSTEM_TYPE]->(st)	
	OPTIONAL MATCH (sys)-[:HAS_OWNER]->(own)
	OPTIONAL MATCH (sys)-[:HAS_RESPONSIBLE]->(responsilbe)
	OPTIONAL MATCH (sys)-[:HAS_IMPORTANCE]->(imp)
		
    return count(sys) as count
`
	result.ReturnAlias = "count"

	result.Parameters["search"] = strings.ToLower(searchString)
	result.Parameters["fulltextSearch"] = strings.ToLower(helpers.GetFullTextSearchString(searchString))
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func GetSubSystemsQuery(parentUID string, facilityCode string) (result helpers.DatabaseQuery) {

	result.Parameters = make(map[string]interface{})

	result.Query = `
	MATCH(f:Facility{code: $facilityCode}) 
	WITH f
	MATCH(parent:System{uid: $parentUID})-[:BELONGS_TO_FACILITY]->(f) WITH parent
	MATCH(parent)-[:HAS_SUBSYSTEM]->(sys{ deleted: false }) WITH sys `

	result.Query += `	
	OPTIONAL MATCH (sys)-[:HAS_LOCATION]->(loc)  
	OPTIONAL MATCH (sys)-[:HAS_ZONE]->(zone)  
	OPTIONAL MATCH (sys)-[:HAS_SYSTEM_TYPE]->(st)	
	OPTIONAL MATCH (sys)-[:HAS_OWNER]->(own)
	OPTIONAL MATCH (sys)-[:HAS_RESPONSIBLE]->(responsilbe)
	OPTIONAL MATCH (sys)-[:HAS_IMPORTANCE]->(imp)
	OPTIONAL MATCH (sys)-[:CONTAINS_ITEM]->(physicalItem)-[:IS_BASED_ON]->(catalogueItem)-[:BELONGS_TO_CATEGORY]->(ciCategory)	
	OPTIONAL MATCH (physicalItem)-[:HAS_ITEM_USAGE]->(itemUsage)
	OPTIONAL MATCH (parents)-[:HAS_SUBSYSTEM*1..50]->(sys)
	OPTIONAL MATCH (sys)-[:HAS_SUBSYSTEM*1..50]->(subsys)
	RETURN DISTINCT {  
		uid: sys.uid,
	description: sys.description,
	name: sys.name,
	parentPath: case when parents is not null then reverse(collect(distinct {uid: parents.uid, name: parents.name})) else null end,
	hasSubsystems: case when subsys is not null then true else false end,
	systemCode: sys.systemCode,
	systemAlias: sys.systemAlias,
	systemLevel: sys.systemLevel,
	isTechnologicalUnit: sys.isTechnologicalUnit,
	location: case when loc is not null then {uid: loc.code, name: loc.name} else null end,
	zone: case when zone is not null then {uid: zone.uid, name: zone.name} else null end,
	systemType: case when st is not null then {uid: st.uid, name: st.name} else null end,
	owner: case when own is not null then {uid: own.uid, name: own.lastName + " " + own.firstName} else null end,
	responsible: case when responsilbe is not null then {uid: responsilbe.uid, name: responsilbe.lastName + " " + responsilbe.firstName} else null end,
	importance: case when imp is not null then {uid: imp.uid, name: imp.name} else null end,	
	lastUpdateTime: sys.lastUpdateTime,
	lastUpdateBy: sys.lastUpdateBy,
	physicalItem: case when physicalItem is not null then {
		uid: physicalItem.uid, 
		eun: physicalItem.eun, 
		serialNumber: physicalItem.serialNumber,
		price: physicalItem.price,
		currency: physicalItem.currency,
		itemUsage: case when itemUsage is not null then {uid: itemUsage.uid, name: itemUsage.name} else null end,
		catalogueItem: case when catalogueItem is not null then {
			uid: catalogueItem.uid,
			name: catalogueItem.name,
			catalogueNumber: catalogueItem.catalogueNumber,
			category: case when ciCategory is not null then {uid: ciCategory.uid, name: ciCategory.name} else null end
		} else null end	
		} else null end,
		statistics: {subsystemsCount: count(subsys)}
		} AS systems
	ORDER BY systems.name
	LIMIT 1000
`
	result.ReturnAlias = "systems"

	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["parentUID"] = parentUID

	return result
}

func SystemDetailQuery(uid string, facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH(sys:System{uid: $uid, deleted: false})-[:BELONGS_TO_FACILITY]->(f) WHERE f.code = $facilityCode
	WITH sys
	OPTIONAL MATCH (sys)-[:HAS_LOCATION]->(loc)  
	OPTIONAL MATCH (sys)-[:HAS_ZONE]->(zone)  
	OPTIONAL MATCH (sys)-[:HAS_SYSTEM_TYPE]->(st)	
	OPTIONAL MATCH (sys)-[:HAS_OWNER]->(own)
	OPTIONAL MATCH (sys)-[:HAS_RESPONSIBLE]->(responsilbe)
	OPTIONAL MATCH (sys)-[:HAS_IMPORTANCE]->(imp)
	OPTIONAL MATCH (sys)-[:CONTAINS_ITEM]->(physicalItem)-[:IS_BASED_ON]->(catalogueItem)-[:BELONGS_TO_CATEGORY]->(ciCategory)	
	OPTIONAL MATCH (physicalItem)-[:HAS_ITEM_USAGE]->(itemUsage)
	OPTIONAL MATCH (parents)-[:HAS_SUBSYSTEM*1..50]->(sys)
	OPTIONAL MATCH (sys)-[:HAS_SUBSYSTEM*1..50]->(subsys)
	RETURN DISTINCT {  
	uid: sys.uid,
	description: sys.description,
	name: sys.name,
	parentPath: case when parents is not null then reverse(collect(distinct {uid: parents.uid, name: parents.name})) else null end,
	systemCode: sys.systemCode,
	systemAlias: sys.systemAlias,
	systemLevel: sys.systemLevel,
	isTechnologicalUnit: sys.isTechnologicalUnit,
	location: case when loc is not null then {uid: loc.code, name: loc.name} else null end,
	zone: case when zone is not null then {uid: zone.uid, name: zone.name} else null end,
	systemType: case when st is not null then {uid: st.uid, name: st.name} else null end,
	owner: case when own is not null then {uid: own.uid, name: own.lastName + " " + own.firstName} else null end,
	responsible: case when responsilbe is not null then {uid: responsilbe.uid, name: responsilbe.lastName + " " + responsilbe.firstName} else null end,
	importance: case when imp is not null then {uid: imp.uid, name: imp.name} else null end,	
	lastUpdateTime: sys.lastUpdateTime,
	lastUpdateBy: sys.lastUpdateBy,
	physicalItem: case when physicalItem is not null then {
		uid: physicalItem.uid, 
		eun: physicalItem.eun, 
		serialNumber: physicalItem.serialNumber,
		price: physicalItem.price,
		currency: physicalItem.currency,
		itemUsage: case when itemUsage is not null then {uid: itemUsage.uid, name: itemUsage.name} else null end,
		catalogueItem: case when catalogueItem is not null then {
			uid: catalogueItem.uid,
			name: catalogueItem.name,
			catalogueNumber: catalogueItem.catalogueNumber,
			category: case when ciCategory is not null then {uid: ciCategory.uid, name: ciCategory.name} else null end
		} else null end	
	} else null end,
	statistics: {subsystemsCount: count(subsys)}
} AS system`
	result.ReturnAlias = "system"
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func GetSystemRelationshipsQuery(uid string) (result helpers.DatabaseQuery) {
	result.Query = `
	MATCH(sys:System{uid: $uid, deleted: false})
	MATCH (parents)-[rParent:HAS_SUBSYSTEM]->(sys)	
	return distinct {
		direction: "to",
		foreignSystemName: parents.name,
		relationUid: id(rParent),
		relationTypeCode: "HAS_SUBSYSTEM"
		} as relationships
	UNION
	MATCH(sys:System{uid: $uid, deleted: false})
	MATCH (sys)-[rSubsys:HAS_SUBSYSTEM]->(subsys)	
	return distinct {
		direction: "from",
		foreignSystemName: subsys.name,
		relationUid: id(rSubsys),
		relationTypeCode: "HAS_SUBSYSTEM"
		} as relationships
	UNION
	MATCH(sys:System{uid: $uid, deleted: false})
	MATCH (parents)-[rParent:IS_SPARE_FOR]->(sys)	
	return distinct {
		direction: "to",
		foreignSystemName: parents.name,
		relationUid: id(rParent),
		relationTypeCode: "IS_SPARE_FOR"
		} as relationships
	UNION
	MATCH(sys:System{uid: $uid, deleted: false})
	MATCH (sys)-[rSubsys:IS_SPARE_FOR]->(subsys)	
	return distinct {
		direction: "from",
		foreignSystemName: subsys.name,
		relationUid: id(rSubsys),
		relationTypeCode: "IS_SPARE_FOR"
		} as relationships;`

	result.ReturnAlias = "relationships"
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid

	return result
}

// create new system relationship query
func CreateNewSystemRelationshipQuery(newRelationship *models.SystemRelationshipRequest, facilityCode string, userUID string) (result helpers.DatabaseQuery) {

	result.Parameters = make(map[string]interface{})
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["uid"] = uuid.NewString()
	result.Parameters["fromSystemUID"] = newRelationship.SystemFromUID
	result.Parameters["toSystemUID"] = newRelationship.SystemToUID
	result.Parameters["relationshipTypeCode"] = newRelationship.RelationTypeCode
	result.Parameters["lastUpdateBy"] = userUID

	result.Query = `
	MATCH(f:Facility{code: $facilityCode}) WITH f	
	MATCH(u:User{uid: $lastUpdateBy}) WITH u, f
	MATCH(fromSystem:System{uid: $fromSystemUID, deleted: false})-[:BELONGS_TO_FACILITY]->(f)
	MATCH(toSystem:System{uid: $toSystemUID, deleted: false})-[:BELONGS_TO_FACILITY]->(f)`

	if newRelationship.RelationTypeCode == "IS_SPARE_FOR" {
		result.Query += `CREATE(fromSystem)-[newRel:IS_SPARE_FOR]->(toSystem) `
	} else {
		result.Query += `REALTIONSHIP NOT DEFINED`
	}

	result.Query += `
	WITH fromSystem, toSystem, u, newRel
	CREATE(fromSystem)-[:WAS_UPDATED_BY{ at: datetime(), action: "UPDATE" }]->(u)	
	WITH fromSystem, toSystem, newRel
	CREATE(toSystem)-[:WAS_UPDATED_BY{ at: datetime(), action: "UPDATE" }]->(u)	
	WITH fromSystem, toSystem, newRel
	`

	result.Query += `RETURN id(newRel) as result`

	result.ReturnAlias = "result"

	return result
}

func DeleteSystemRelationshipQuery(uid int64, userUID string) (result helpers.DatabaseQuery) {
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid
	result.Parameters["lastUpdateBy"] = userUID

	result.Query = `
	MATCH (u:User{uid: $lastUpdateBy}) WITH u
	MATCH ()-[r]-() WHERE id(r) = $uid DELETE r
	WITH u
	CREATE(u)-[:WAS_UPDATED_BY{ at: datetime(), action: "DELETE" }]->(u)	
	return true as result`

	result.ReturnAlias = "result"

	return result
}
