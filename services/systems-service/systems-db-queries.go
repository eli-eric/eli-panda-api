package systemsService

import (
	"fmt"
	"panda/apigateway/helpers"
	codebookModels "panda/apigateway/services/codebook-service/models"
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
	where $searchText = '' or ((toLower(n.code) contains $searchText or toLower(n.name) contains $searchText))
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

func GetZonesCodebookQuery(facilityCode string, searchString string) (result helpers.DatabaseQuery) {
	result.Query = `
	MATCH(f:Facility{code:$facilityCode}) WITH f	
	MATCH(z:Zone)-[:BELONGS_TO_FACILITY]->(f) where not ()-[:HAS_SUBZONE]->(z) AND ($searchString = '' OR toLower(z.code) contains $searchString) return {uid:z.uid, name:z.code + " - " +z.name } as zone order by z.code
	UNION
	MATCH(f:Facility{code:$facilityCode}) WITH f
	MATCH(z:Zone)-[:HAS_SUBZONE]->(sz)-[:BELONGS_TO_FACILITY]->(f) where ($searchString = '' OR toLower(sz.code) contains $searchString) return {uid:sz.uid, name: z.code+"-"+sz.code + " - " + sz.name + " ("+  z.name + ")"} as zone order by z.code, sz.code`
	result.ReturnAlias = "zone"
	result.Parameters = make(map[string]interface{})
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["searchString"] = strings.ToLower(searchString)

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

func DeleteSystemByUidQuery(uid, userUid string) (result helpers.DatabaseQuery) {

	result.Query = `MATCH(r:System) WHERE r.uid = $uid WITH r
	OPTIONAL MATCH (r)-[:HAS_SUBSYSTEM*1..50]->(child)
	WITH r, child, r.uid as uid
	SET r.deleted=true, child.deleted=true
	WITH r, child
	MATCH(u:User{uid: $userUID})
	CREATE(r)-[:WAS_UPDATED_BY{ at: datetime(), action: "DELETE" }]->(u)
	CREATE(child)-[:WAS_UPDATED_BY{ at: datetime(), action: "DELETE" }]->(u)`

	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid
	result.Parameters["userUID"] = userUid

	return result
}

func GetPhysicalItemsBySystemUidRecursiveQuery(systemUid string) (result helpers.DatabaseQuery) {

	result.Query = `
	MATCH(s:System) WHERE s.uid = $systemUid WITH s
	OPTIONAL MATCH (s)-[:HAS_SUBSYSTEM*1..50]->(subs)
	WITH collect(s.uid) as ss, collect(subs.uid) as subss
	WITH ss + subss as sUids
	OPTIONAL MATCH(systems:System)-[:CONTAINS_ITEM]->(itm) WHERE systems.uid in sUids
	RETURN CASE WHEN systems IS NULL THEN NULL ELSE { systemUid: systems.uid, itemUid: itm.uid, systemName: systems.name, itemName: itm.name } END as result;
	`

	result.Parameters = make(map[string]interface{})
	result.Parameters["systemUid"] = systemUid
	result.ReturnAlias = "result"

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

func GetSystemsSearchFilterQueryOnly(searchString string, facilityCode string, filering *[]helpers.ColumnFilter) (result helpers.DatabaseQuery) {

	result.Parameters = make(map[string]interface{})
	result.Parameters["search"] = strings.ToLower(searchString)
	result.Parameters["fulltextSearch"] = strings.ToLower(helpers.GetFullTextSearchString(searchString))
	result.Parameters["facilityCode"] = facilityCode
	catalogueCategoryFilter := helpers.GetFilterValueCodebook(filering, "category")

	if searchString == "" && (filering == nil || len(*filering) == 0) {
		result.Query = "MATCH(f:Facility{code: $facilityCode}) WITH f MATCH(sys:System{deleted:false})-[:BELONGS_TO_FACILITY]->(f) WHERE NOT ()-[:HAS_SUBSYSTEM]->(sys) WITH sys "
	} else if filering != nil && len(*filering) > 0 {
		// explicitlly set search string to empty string if we have filters
		searchString = ""

		if catalogueCategoryFilter != nil {
			result.Query = `
		MATCH(cat:CatalogueCategory{uid:$filterCatalogueCategory})
		OPTIONAL MATCH(cat)-[:HAS_SUBCATEGORY*1..20]->(subs)
		WITH collect(subs.uid) + cat.uid as catUids
		MATCH (f{code: $facilityCode})<-[:BELONGS_TO_FACILITY]-(sys:System)-[:CONTAINS_ITEM]->(physicalItem)-[:IS_BASED_ON]->(catalogueItem)-[:BELONGS_TO_CATEGORY]->(ciCategory)
		WHERE ciCategory.uid in catUids
		WITH sys, catalogueItem, ciCategory, physicalItem
		`
			result.Parameters["filterCatalogueCategory"] = (*catalogueCategoryFilter).UID
		} else {
			result.Query = "MATCH(f:Facility{code: $facilityCode}) WITH f MATCH(sys:System{deleted:false})-[:BELONGS_TO_FACILITY]->(f) WITH sys "
			result.Query += `
			OPTIONAL MATCH (sys)-[:CONTAINS_ITEM]->(physicalItem)-[:IS_BASED_ON]->(catalogueItem)-[:BELONGS_TO_CATEGORY]->(ciCategory) WITH sys, physicalItem, catalogueItem, ciCategory
			`
		}
	} else {

		// check if the search string is uuid
		_, err := uuid.Parse(searchString)

		if err == nil {
			result.Query = `
		MATCH(f:Facility{code: $facilityCode}) WITH f
		MATCH(sys:System{uid: $search})-[:BELONGS_TO_FACILITY]->(f) WITH sys
		`
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
	}

	//apply filters

	//parentSystem
	if filterVal := helpers.GetFilterValueCodebook(filering, "parentSystem"); filterVal != nil {
		result.Query += ` MATCH(p{uid: $filterParentSystemUID})-[:HAS_SUBSYSTEM*1..50]->(sys) WITH sys, physicalItem, catalogueItem, ciCategory `
		result.Parameters["filterParentSystemUID"] = (*filterVal).UID
	}

	//system name
	if filterVal := helpers.GetFilterValueString(filering, "name"); filterVal != nil {
		result.Query += ` WHERE toLower(sys.name) CONTAINS $filterName `
		result.Parameters["filterName"] = strings.ToLower(*filterVal)
	}

	//system description
	if filterVal := helpers.GetFilterValueString(filering, "description"); filterVal != nil {
		result.Query += ` WITH sys, physicalItem, catalogueItem, ciCategory WHERE toLower(sys.description) CONTAINS $filterDescription `
		result.Parameters["filterDescription"] = strings.ToLower(*filterVal)
	}

	//system level
	if filterVal := helpers.GetFilterValueListString(filering, "systemLevel"); filterVal != nil {
		result.Query += ` WITH sys, physicalItem, catalogueItem, ciCategory WHERE sys.systemLevel IN $filterSystemLevel `
		result.Parameters["filterSystemLevel"] = filterVal
	}

	//system code
	if filterVal := helpers.GetFilterValueString(filering, "systemCode"); filterVal != nil {
		result.Query += ` WITH sys, physicalItem, catalogueItem, ciCategory WHERE toLower(sys.systemCode) CONTAINS $filterSystemCode `
		result.Parameters["filterSystemCode"] = strings.ToLower(*filterVal)
	}

	// system sp_coverage
	if filterVal := helpers.GetFilterValueRangeFloat64(filering, "sparePartsCoverage"); filterVal != nil {
		result.Query += ` WITH sys, physicalItem, catalogueItem, ciCategory WHERE sys.sp_coverage IS NOT NULL AND (($filterSPCoverageFrom IS NULL OR toFloat(sys.sp_coverage) >= $filterSPCoverageFrom) AND ($filterSPCoverageTo IS NULL OR toFloat(sys.sp_coverage) <= $filterSPCoverageTo)) `
		result.Parameters["filterSPCoverageFrom"] = filterVal.Min
		result.Parameters["filterSPCoverageTo"] = filterVal.Max
	}

	//system type
	if filterVal := helpers.GetFilterValueCodebook(filering, "systemType"); filterVal != nil {
		result.Query += ` MATCH (sys)-[:HAS_SYSTEM_TYPE]->(st) WHERE st.uid = $filterSystemType `
		result.Parameters["filterSystemType"] = (*filterVal).UID
	} else {
		result.Query += `OPTIONAL MATCH (sys)-[:HAS_SYSTEM_TYPE]->(st) `
	}

	//system location
	if filterVal := helpers.GetFilterValueCodebook(filering, "location"); filterVal != nil {
		result.Query += ` MATCH (sys)-[:HAS_LOCATION]->(loc) WHERE loc.uid = $filterLocation `
		result.Parameters["filterLocation"] = (*filterVal).UID
	} else {
		result.Query += ` OPTIONAL MATCH (sys)-[:HAS_LOCATION]->(loc) `
	}

	//system zone
	if filterVal := helpers.GetFilterValueCodebook(filering, "zone"); filterVal != nil {
		result.Query += ` MATCH (sys)-[:HAS_ZONE]->(zone) WHERE zone.uid = $filterZone `
		result.Parameters["filterZone"] = (*filterVal).UID
	} else {
		result.Query += ` OPTIONAL MATCH (sys)-[:HAS_ZONE]->(zone) `
	}

	//system responsible
	if filterVal := helpers.GetFilterValueCodebook(filering, "responsible"); filterVal != nil {
		result.Query += ` MATCH (sys)-[:HAS_RESPONSIBLE]->(responsible) WHERE responsible.uid = $filterResponsible `
		result.Parameters["filterResponsible"] = (*filterVal).UID
	} else {
		result.Query += ` OPTIONAL MATCH (sys)-[:HAS_RESPONSIBLE]->(responsible) `
	}

	//system importance
	if filterVal := helpers.GetFilterValueCodebook(filering, "importance"); filterVal != nil {
		result.Query += ` MATCH (sys)-[:HAS_IMPORTANCE]->(imp) WHERE imp.uid = $filterImportance `
		result.Parameters["filterImportance"] = (*filterVal).UID
	} else {
		result.Query += ` OPTIONAL MATCH (sys)-[:HAS_IMPORTANCE]->(imp) `
	}

	//system attribute - TODO: has to be on SystemType
	// if filterVal := helpers.GetFilterValueCodebook(filering, "systemAttribute"); filterVal != nil {
	// 	result.Query += ` MATCH (sys)-[:HAS_SYSTEM_ATTRIBUTE]->(sysAttr) WHERE sysAttr.uid = $filterSystemAttribute `
	// 	result.Parameters["filterSystemAttribute"] = (*filterVal).UID
	// } else {
	// 	result.Query += ` OPTIONAL MATCH (sys)-[:HAS_SYSTEM_ATTRIBUTE]->(sysAttr) `
	// }

	//physical item filters
	//we have to get all physical items filter values first and then apply them
	itemUsageFilter := helpers.GetFilterValueListString(filering, "itemUsage")
	eunFilter := helpers.GetFilterValueString(filering, "eun")
	serialNumberFilter := helpers.GetFilterValueString(filering, "serialNumber")
	catalogueNumberFilter := helpers.GetFilterValueString(filering, "catalogueNumber")
	catalogueNameFilter := helpers.GetFilterValueString(filering, "catalogueName")
	supplierFilter := helpers.GetFilterValueCodebook(filering, "supplier")
	priceFilter := helpers.GetFilterValueRangeFloat64(filering, "price")
	orderFilter := helpers.GetFilterValueString(filering, "order")
	orderNameFilter := helpers.GetFilterValueString(filering, "orderName")
	orderNumberFilter := helpers.GetFilterValueString(filering, "orderNumber")
	orderRequestNumberFilter := helpers.GetFilterValueString(filering, "orderRequestNumber")
	orderContractNumberFilter := helpers.GetFilterValueString(filering, "orderContractNumber")

	if orderContractNumberFilter != nil || orderRequestNumberFilter != nil || orderNumberFilter != nil || orderNameFilter != nil || itemUsageFilter != nil || eunFilter != nil || serialNumberFilter != nil || catalogueNumberFilter != nil || catalogueNameFilter != nil || supplierFilter != nil || catalogueCategoryFilter != nil || orderFilter != nil { // || priceFilter != nil {

		if orderFilter != nil {
			result.Query += ` MATCH (physicalItem)<-[:HAS_ORDER_LINE]-(order) WHERE order.orderNumber CONTAINS $filterOrder OR order.requestNumber CONTAINS $filterOrder OR order.contractNumber CONTAINS $filterOrder `
			result.Parameters["filterOrder"] = *orderFilter
		}

		if orderNameFilter != nil {
			result.Query += ` MATCH (physicalItem)<-[:HAS_ORDER_LINE]-(order) WHERE order.name CONTAINS $filterOrderName `
			result.Parameters["filterOrderName"] = strings.ToLower(*orderNameFilter)
		}

		if orderNumberFilter != nil {
			result.Query += ` MATCH (physicalItem)<-[:HAS_ORDER_LINE]-(order) WHERE order.orderNumber CONTAINS $filterOrderNumber `
			result.Parameters["filterOrderNumber"] = *orderNumberFilter
		}

		if orderRequestNumberFilter != nil {
			result.Query += ` MATCH (physicalItem)<-[:HAS_ORDER_LINE]-(order) WHERE order.requestNumber CONTAINS $filterOrderRequestNumber `
			result.Parameters["filterOrderRequestNumber"] = *orderRequestNumberFilter
		}

		if orderContractNumberFilter != nil {
			result.Query += ` MATCH (physicalItem)<-[:HAS_ORDER_LINE]-(order) WHERE order.contractNumber CONTAINS $filterOrderContractNumber `
			result.Parameters["filterOrderContractNumber"] = *orderContractNumberFilter
		}

		if orderFilter == nil && orderNameFilter == nil && orderNumberFilter == nil && orderRequestNumberFilter == nil && orderContractNumberFilter == nil {
			result.Query += ` OPTIONAL MATCH (physicalItem)<-[:HAS_ORDER_LINE]-(order) `
		}

		if itemUsageFilter != nil {
			result.Query += ` MATCH (physicalItem)-[:HAS_ITEM_USAGE]->(itemUsage) WHERE itemUsage.uid IN $filterItemUsage `
			result.Parameters["filterItemUsage"] = itemUsageFilter
		} else {
			result.Query += ` OPTIONAL MATCH (physicalItem)-[:HAS_ITEM_USAGE]->(itemUsage) `
		}

		if supplierFilter != nil {
			result.Query += ` MATCH (catalogueItem)-[:HAS_SUPPLIER]->(supplier) WHERE supplier.uid = $filterSupplier `
			result.Parameters["filterSupplier"] = (*supplierFilter).UID
		} else {
			result.Query += ` OPTIONAL MATCH (catalogueItem)-[:HAS_SUPPLIER]->(supplier) `
		}

		if priceFilter != nil {
			result.Query += ` WITH sys, physicalItem, catalogueItem, ciCategory, itemUsage, imp, responsible, loc, zone, st, supplier, order `
			result.Query += ` MATCH (physicalItem)<-[ol:HAS_ORDER_LINE]-() WHERE ($filterPriceFrom IS NULL OR ol.price >= $filterPriceFrom) AND ($filterPriceTo IS NULL OR ol.price <= $filterPriceTo) `
			result.Parameters["filterPriceFrom"] = priceFilter.Min
			result.Parameters["filterPriceTo"] = priceFilter.Max
		} else {
			result.Query += ` WITH sys, physicalItem, catalogueItem, ciCategory, itemUsage, imp, responsible, loc, zone, st, supplier, order `
			result.Query += ` OPTIONAL MATCH (physicalItem)<-[ol:HAS_ORDER_LINE]-() `
		}

		if eunFilter != nil {
			result.Query += ` WITH sys, physicalItem, catalogueItem, ciCategory, itemUsage, imp, responsible, loc, zone, st, supplier, ol, order `
			result.Query += ` WHERE toLower(physicalItem.eun) CONTAINS $filterEUN `
			result.Parameters["filterEUN"] = strings.ToLower(*eunFilter)
		}

		if serialNumberFilter != nil {
			result.Query += ` WITH sys, physicalItem, catalogueItem, ciCategory, itemUsage, imp, responsible, loc, zone, st, supplier, ol, order `
			result.Query += ` WHERE toLower(physicalItem.serialNumber) CONTAINS $filterSerialNumber `
			result.Parameters["filterSerialNumber"] = strings.ToLower(*serialNumberFilter)
		}

		if catalogueNumberFilter != nil {
			result.Query += ` WITH sys, physicalItem, catalogueItem, ciCategory, itemUsage, imp, responsible, loc, zone, st, supplier , ol, order`
			result.Query += ` WHERE toLower(catalogueItem.catalogueNumber) CONTAINS $filterCatalogueNumber `
			result.Parameters["filterCatalogueNumber"] = strings.ToLower(*catalogueNumberFilter)
		}

		if catalogueNameFilter != nil {
			result.Query += ` WITH sys, physicalItem, catalogueItem, ciCategory, itemUsage, imp, responsible, loc, zone, st, supplier, ol, order `
			result.Query += ` WHERE toLower(catalogueItem.name) CONTAINS $filterCatalogueName `
			result.Parameters["filterCatalogueName"] = strings.ToLower(*catalogueNameFilter)
		}

		for i, filter := range *filering {
			if filter.Type == "" {
				continue
			}

			if filter.Type == "text" {
				if filterPropvalue := helpers.GetFilterValueString(filering, filter.Id); filterPropvalue != nil {
					result.Query += ` WITH sys, physicalItem, catalogueItem, ciCategory, itemUsage, imp, responsible, loc, zone, st, supplier, ol, order `
					result.Query += fmt.Sprintf(` MATCH(prop{uid: $propUID%v})<-[pv]-(%v) WHERE toLower(pv.value) contains $propFilterVal%v `, i, itemTypeByPropType[filter.PropType], i)
					result.Parameters[fmt.Sprintf("propUID%v", i)] = filter.Id
					result.Parameters[fmt.Sprintf("propFilterVal%v", i)] = strings.ToLower(*filterPropvalue)
				}
			} else if filter.Type == "list" {
				if filterPropvalue := helpers.GetFilterValueListString(filering, filter.Id); filterPropvalue != nil {
					result.Query += ` WITH sys, physicalItem, catalogueItem, ciCategory, itemUsage, imp, responsible, loc, zone, st, supplier, ol, order `
					result.Query += fmt.Sprintf(` MATCH(prop{uid: $propUID%v})<-[pv]-(%v) WHERE pv.value IN $propFilterVal%v `, i, itemTypeByPropType[filter.PropType], i)
					result.Parameters[fmt.Sprintf("propUID%v", i)] = filter.Id
					result.Parameters[fmt.Sprintf("propFilterVal%v", i)] = filterPropvalue
				}
			} else if filter.Type == "number" {
				if filterPropvalue := helpers.GetFilterValueRangeFloat64(filering, filter.Id); filterPropvalue != nil {
					result.Query += ` WITH sys, physicalItem, catalogueItem, ciCategory, itemUsage, imp, responsible, loc, zone, st, supplier, ol, order `
					result.Query += fmt.Sprintf(` MATCH(prop{uid: $propUID%v})<-[pv]-(%v) WHERE ($propFilterValFrom%v IS NULL OR toFloat(pv.value) >= $propFilterValFrom%v) AND ($propFilterValTo%v IS NULL OR toFloat(pv.value) <= $propFilterValTo%v) `, i, itemTypeByPropType[filter.PropType], i, i, i, i)
					result.Parameters[fmt.Sprintf("propUID%v", i)] = filter.Id

					result.Parameters[fmt.Sprintf("propFilterValFrom%v", i)] = filterPropvalue.Min
					result.Parameters[fmt.Sprintf("propFilterValTo%v", i)] = filterPropvalue.Max
				}
			} else if filter.Type == "range" {
				if filterPropvalue := helpers.GetFilterValueRangeFloat64(filering, filter.Id); filterPropvalue != nil {
					result.Query += ` WITH sys, physicalItem, catalogueItem, ciCategory, itemUsage, imp, responsible, loc, zone, st, supplier, ol, order `
					result.Query += fmt.Sprintf(` MATCH(prop{uid: $propUID%v})<-[pv]-(%v)
					WITH sys, physicalItem, catalogueItem, ciCategory, itemUsage, imp, responsible, loc, zone, st, supplier, ol, order, apoc.convert.fromJsonMap(pv.value) as jsonValue					
					WHERE (toFloat(jsonValue.min) <= $propFilterValFrom%[1]v - $propFilterValTo%[1]v) 
					AND (toFloat(jsonValue.max) >= $propFilterValFrom%[1]v + $propFilterValTo%[1]v) 
					`, i, itemTypeByPropType[filter.PropType])

					result.Parameters[fmt.Sprintf("propUID%v", i)] = filter.Id
					result.Parameters[fmt.Sprintf("propFilterValFrom%v", i)] = filterPropvalue.Min
					if filterPropvalue.Max == nil {
						filterPropvalue.Max = &plusMinusZero
					}
					result.Parameters[fmt.Sprintf("propFilterValTo%v", i)] = filterPropvalue.Max
				}
			}
		}

	} else {
		if catalogueCategoryFilter != nil {
			result.Query += ` 
			OPTIONAL MATCH (physicalItem)<-[ol:HAS_ORDER_LINE]-(order)
			OPTIONAL MATCH (physicalItem)-[:HAS_ITEM_USAGE]->(itemUsage) 
			OPTIONAL MATCH (catalogueItem)-[:HAS_SUPPLIER]->(supplier) `
		} else {
			result.Query += ` 
			OPTIONAL MATCH (sys)-[:CONTAINS_ITEM]->(physicalItem)-[:IS_BASED_ON]->(catalogueItem)-[:BELONGS_TO_CATEGORY]->(ciCategory) 		
			OPTIONAL MATCH (physicalItem)<-[ol:HAS_ORDER_LINE]-(order)
			OPTIONAL MATCH (physicalItem)-[:HAS_ITEM_USAGE]->(itemUsage) 
			OPTIONAL MATCH (catalogueItem)-[:HAS_SUPPLIER]->(supplier) `
		}
	}

	return result
}

var plusMinusZero float64 = 0

// this map is used to map filter prop type to type of item that has this prop related to in the query/func "GetSystemsSearchFilterQueryOnly"
var itemTypeByPropType = map[string]string{
	"CATALOGUE_ITEM": "catalogueItem",
	"PHYSICAL_ITEM":  "physicalItem",
}

func GetSystemsOrderByClauses(sorting *[]helpers.Sorting) string {

	if sorting == nil || len(*sorting) == 0 {
		return `ORDER BY systems.hasSubsystems desc, systems.systemLevelOrder, systems.lastUpdateTime DESC `
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

func GetSystemsBySearchTextFullTextQuery(searchString string, facilityCode string, pagination *helpers.Pagination, sorting *[]helpers.Sorting, filering *[]helpers.ColumnFilter) (result helpers.DatabaseQuery) {

	result = GetSystemsSearchFilterQueryOnly(searchString, facilityCode, filering)

	result.Query += `
	OPTIONAL MATCH (parents{deleted: false})-[:HAS_SUBSYSTEM*1..50]->(sys)
	OPTIONAL MATCH (sys)-[:HAS_SUBSYSTEM*1..50]->(subsys{deleted: false})
	OPTIONAL MATCH (sys)-[:IS_SPARE_FOR]->(spareOUT)
    OPTIONAL MATCH (sys)<-[:IS_SPARE_FOR]-(spareIN)
	RETURN DISTINCT {  
	uid: sys.uid,
	description: sys.description,
	name: sys.name,
	parentPath: case when parents is not null then reverse(collect(distinct {uid: parents.uid, name: parents.name})) else null end,
	hasSubsystems: case when subsys is not null then true else false end,
	sparesIn: count(distinct spareIN),
	sparesOut: count(distinct spareOUT),
	systemCode: sys.systemCode,
	systemAlias: sys.systemAlias,
	systemLevel: sys.systemLevel,
	minimalSpareParstCount: sys.minimalSpareParstCount,
	sparePartsCoverageSum: sys.sparePartsCoverageSum,
	sp_coverage: sys.sp_coverage,
	miniImageUrl: split(sys.miniImageUrl, ";"),	
	systemLevelOrder: case sys.systemLevel WHEN 'TECHNOLOGY_UNIT' THEN 1 WHEN 'KEY_SYSTEMS' THEN 2 ELSE 3 END,
	isTechnologicalUnit: sys.isTechnologicalUnit,
	location: case when loc is not null then {uid: loc.uid, name: loc.name, code: loc.code} else null end,
	zone: case when zone is not null then {uid: zone.uid, name: zone.name, code: zone.code} else null end,
	systemType: case when st is not null then {uid: st.uid, name: st.name} else null end,	
	responsible: case when responsible is not null then {uid: responsible.uid, name: responsible.lastName + " " + responsible.firstName} else null end,
	importance: case when imp is not null then {uid: imp.uid, name: imp.name} else null end,	
	lastUpdateTime: sys.lastUpdateTime,
	lastUpdateBy: sys.lastUpdateBy,	
	physicalItem: case when physicalItem is not null then {
		uid: physicalItem.uid, 
		eun: physicalItem.eun, 
		serialNumber: physicalItem.serialNumber,
		orderNumber: case when order is not null then order.orderNumber else null end,
	    orderUid: case when order is not null then order.uid else null end,
		price: case when ol is not null then apoc.number.format(ol.price, '#,##0') else null end,
		currency: ol.currency,
		itemUsage: case when itemUsage is not null then {uid: itemUsage.uid, name: itemUsage.name} else null end,
		catalogueItem: case when catalogueItem is not null then {
			uid: catalogueItem.uid,
			name: catalogueItem.name,
			catalogueNumber: catalogueItem.catalogueNumber,
			category: case when ciCategory is not null then {uid: ciCategory.uid, name: ciCategory.name} else null end,
			supplier: case when supplier is not null then {uid: supplier.uid, name: supplier.name} else null end
			} else null end	
	} else null end,
	statistics: {
			subsystemsCount: count(DISTINCT subsys),
			minimalSpareParstCount: sys.minimalSpareParstCount,
			sparePartsCoverageSum: sys.sparePartsCoverageSum,
			sp_coverage: sys.sp_coverage
	}
	} AS systems

` + GetSystemsOrderByClauses(sorting) + `

	SKIP $skip
	LIMIT $limit

`
	result.ReturnAlias = "systems"

	result.Parameters["limit"] = pagination.PageSize
	result.Parameters["skip"] = (pagination.Page - 1) * pagination.PageSize

	return result
}

func GetSystemsBySearchTextFullTextCountQuery(searchString string, facilityCode string, filering *[]helpers.ColumnFilter) (result helpers.DatabaseQuery) {

	result = GetSystemsSearchFilterQueryOnly(searchString, facilityCode, filering)

	result.Query += ` RETURN COUNT(sys) as count `
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
	OPTIONAL MATCH (sys)-[:HAS_RESPONSIBLE]->(responsilbe)
	OPTIONAL MATCH (sys)-[:HAS_IMPORTANCE]->(imp)
	OPTIONAL MATCH (sys)-[:CONTAINS_ITEM]->(physicalItem)-[:IS_BASED_ON]->(catalogueItem)-[:BELONGS_TO_CATEGORY]->(ciCategory)	
	OPTIONAL MATCH (catalogueItem)-[:HAS_SUPPLIER]->(supplier)
	OPTIONAL MATCH (physicalItem)-[:HAS_ITEM_USAGE]->(itemUsage)
	OPTIONAL MATCH (parents{deleted: false})-[:HAS_SUBSYSTEM*1..50]->(sys)
	OPTIONAL MATCH (sys)-[:HAS_SUBSYSTEM*1..50]->(subsys{deleted: false})
	OPTIONAL MATCH (sys)-[:IS_SPARE_FOR]->(spareOUT)
    OPTIONAL MATCH (sys)<-[:IS_SPARE_FOR]-(spareIN)
	OPTIONAL MATCH (physicalItem)<-[ol:HAS_ORDER_LINE]-(order)
	
	RETURN DISTINCT {  
		uid: sys.uid,
	description: sys.description,
	name: sys.name,
	parentPath: case when parents is not null then reverse(collect(distinct {uid: parents.uid, name: parents.name})) else null end,
	hasSubsystems: case when subsys is not null then true else false end,
	sparesIn: count(distinct spareIN),
	sparesOut: count(distinct spareOUT),
	systemCode: sys.systemCode,
	systemAlias: sys.systemAlias,
	systemLevel: sys.systemLevel,
	minimalSpareParstCount: sys.minimalSpareParstCount,
	sparePartsCoverageSum: sys.sparePartsCoverageSum,
	sp_coverage: sys.sp_coverage,
	miniImageUrl: split(sys.miniImageUrl, ";"),
	systemLevelOrder: case sys.systemLevel WHEN 'TECHNOLOGY_UNIT' THEN 1 WHEN 'KEY_SYSTEMS' THEN 2 ELSE 3 END,
	isTechnologicalUnit: sys.isTechnologicalUnit,
	location: case when loc is not null then {uid: loc.uid, name: loc.name, code: loc.code} else null end,
	zone: case when zone is not null then {uid: zone.uid, name: zone.name, code: zone.code} else null end,
	systemType: case when st is not null then {uid: st.uid, name: st.name} else null end,
	responsible: case when responsilbe is not null then {uid: responsilbe.uid, name: responsilbe.lastName + " " + responsilbe.firstName} else null end,
	importance: case when imp is not null then {uid: imp.uid, name: imp.name} else null end,	
	lastUpdateTime: sys.lastUpdateTime,
	lastUpdateBy: sys.lastUpdateBy,	
	physicalItem: case when physicalItem is not null then {
		uid: physicalItem.uid, 
		eun: physicalItem.eun, 
		serialNumber: physicalItem.serialNumber,
		orderNumber: case when order is not null then order.orderNumber else null end,
	    orderUid: case when order is not null then order.uid else null end,
		price: case when ol is not null then apoc.number.format(ol.price, '#,##0') else null end,
		currency: physicalItem.currency,
		itemUsage: case when itemUsage is not null then {uid: itemUsage.uid, name: itemUsage.name} else null end,
		catalogueItem: case when catalogueItem is not null then {
			uid: catalogueItem.uid,
			name: catalogueItem.name,
			catalogueNumber: catalogueItem.catalogueNumber,
			category: case when ciCategory is not null then {uid: ciCategory.uid, name: ciCategory.name} else null end,
			supplier: case when supplier is not null then {uid: supplier.uid, name: supplier.name} else null end
		} else null end	
		} else null end,
		statistics: {
			subsystemsCount: count(DISTINCT subsys),
			minimalSpareParstCount: sys.minimalSpareParstCount,
			sparePartsCoverageSum: sys.sparePartsCoverageSum,
			sp_coverage: sys.sp_coverage
	}
		} AS systems
	ORDER BY systems.hasSubsystems desc, systems.systemLevelOrder, systems.name
	LIMIT 1000
`
	result.ReturnAlias = "systems"

	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["parentUID"] = parentUID

	return result
}

func GetSystemsByUidsQuery(uids []string) (result helpers.DatabaseQuery) {

	result.Parameters = make(map[string]interface{})

	result.Query = `	
	MATCH(sys) WHERE sys.uid IN $uids `

	result.Query += `	
	OPTIONAL MATCH (sys)-[:HAS_LOCATION]->(loc)  
	OPTIONAL MATCH (sys)-[:HAS_ZONE]->(zone)  
	OPTIONAL MATCH (sys)-[:HAS_SYSTEM_TYPE]->(st)	
	OPTIONAL MATCH (sys)-[:HAS_RESPONSIBLE]->(responsilbe)
	OPTIONAL MATCH (sys)-[:HAS_IMPORTANCE]->(imp)
	OPTIONAL MATCH (sys)-[:CONTAINS_ITEM]->(physicalItem)-[:IS_BASED_ON]->(catalogueItem)-[:BELONGS_TO_CATEGORY]->(ciCategory)	
	OPTIONAL MATCH (catalogueItem)-[:HAS_SUPPLIER]->(supplier)
	OPTIONAL MATCH (physicalItem)-[:HAS_ITEM_USAGE]->(itemUsage)
	OPTIONAL MATCH (physicalItem)<-[:HAS_ORDER_LINE]-(order)
	OPTIONAL MATCH (parents{deleted: false})-[:HAS_SUBSYSTEM*1..50]->(sys)
	OPTIONAL MATCH (sys)-[:HAS_SUBSYSTEM*1..50]->(subsys{deleted: false})
	OPTIONAL MATCH (sys)-[:IS_SPARE_FOR]->(spareOUT)
    OPTIONAL MATCH (sys)<-[:IS_SPARE_FOR]-(spareIN)
	RETURN DISTINCT {  
		uid: sys.uid,
	description: sys.description,
	name: sys.name,
	parentPath: case when parents is not null then reverse(collect(distinct {uid: parents.uid, name: parents.name})) else null end,
	hasSubsystems: case when subsys is not null then true else false end,
	sparesIn: count(distinct spareIN),
	sparesOut: count(distinct spareOUT),
	systemCode: sys.systemCode,
	systemAlias: sys.systemAlias,
	systemLevel: sys.systemLevel,	
	miniImageUrl: split(sys.miniImageUrl, ";"),
	systemLevelOrder: case sys.systemLevel WHEN 'TECHNOLOGY_UNIT' THEN 1 WHEN 'KEY_SYSTEMS' THEN 2 ELSE 3 END,
	isTechnologicalUnit: sys.isTechnologicalUnit,
	location: case when loc is not null then {uid: loc.uid, name: loc.name, code: loc.code} else null end,
	zone: case when zone is not null then {uid: zone.uid, name: zone.name, code: zone.code} else null end,
	systemType: case when st is not null then {uid: st.uid, name: st.name} else null end,
	responsible: case when responsilbe is not null then {uid: responsilbe.uid, name: responsilbe.lastName + " " + responsilbe.firstName} else null end,
	importance: case when imp is not null then {uid: imp.uid, name: imp.name} else null end,	
	lastUpdateTime: sys.lastUpdateTime,
	lastUpdateBy: sys.lastUpdateBy,	
	physicalItem: case when physicalItem is not null then {
		uid: physicalItem.uid, 
		eun: physicalItem.eun, 
		serialNumber: physicalItem.serialNumber,
		orderNumber: case when order is not null then order.orderNumber else null end,
		price: physicalItem.price,
		currency: physicalItem.currency,
		itemUsage: case when itemUsage is not null then {uid: itemUsage.uid, name: itemUsage.name} else null end,
		catalogueItem: case when catalogueItem is not null then {
			uid: catalogueItem.uid,
			name: catalogueItem.name,
			catalogueNumber: catalogueItem.catalogueNumber,
			category: case when ciCategory is not null then {uid: ciCategory.uid, name: ciCategory.name} else null end,
			supplier: case when supplier is not null then {uid: supplier.uid, name: supplier.name} else null end
		} else null end	
		} else null end,
	statistics: {
			subsystemsCount: count(DISTINCT subsys),
			minimalSpareParstCount: sys.minimalSpareParstCount,
			sparePartsCoverageSum: sys.sparePartsCoverageSum,
			sp_coverage: sys.sp_coverage
	}
		} AS systems
	ORDER BY systems.hasSubsystems desc, systems.systemLevelOrder, systems.name
	LIMIT 1000
`
	result.ReturnAlias = "systems"

	result.Parameters["uids"] = uids

	return result
}

func SystemDetailQuery(uid string, facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH(sys:System{uid: $uid, deleted: false})-[:BELONGS_TO_FACILITY]->(f) WHERE f.code = $facilityCode
	WITH sys
	OPTIONAL MATCH (sys)-[:HAS_LOCATION]->(loc)  
	OPTIONAL MATCH (sys)-[:HAS_ZONE]->(zone)  
	OPTIONAL MATCH (sys)-[:HAS_SYSTEM_TYPE]->(st)	
	OPTIONAL MATCH (sys)-[:HAS_RESPONSIBLE]->(responsilbe)
	OPTIONAL MATCH (sys)-[:HAS_IMPORTANCE]->(imp)
	OPTIONAL MATCH (sys)-[:CONTAINS_ITEM]->(physicalItem)-[:IS_BASED_ON]->(catalogueItem)-[:BELONGS_TO_CATEGORY]->(ciCategory)	
	OPTIONAL MATCH (physicalItem)-[:HAS_ITEM_USAGE]->(itemUsage)
	OPTIONAL MATCH (parents{deleted: false})-[:HAS_SUBSYSTEM*1..50]->(sys)
	OPTIONAL MATCH (sys)-[:HAS_SUBSYSTEM*1..50]->(subsys{deleted: false})	
	RETURN DISTINCT {  
	uid: sys.uid,
	description: sys.description,
	name: sys.name,
	parentPath: case when parents is not null then reverse(collect(distinct {uid: parents.uid, name: parents.name})) else null end,
	systemCode: sys.systemCode,
	systemAlias: sys.systemAlias,
	systemLevel: sys.systemLevel,
	minimalSpareParstCount: sys.minimalSpareParstCount,
	sparePartsCoverageSum: sys.sparePartsCoverageSum,
	sp_coverage: sys.sp_coverage,
	miniImageUrl: split(sys.miniImageUrl, ";"),
	isTechnologicalUnit: sys.isTechnologicalUnit,
	location: case when loc is not null then {uid: loc.uid, name: loc.name, code: loc.code} else null end,
	zone: case when zone is not null then {uid: zone.uid, name: zone.name} else null end,
	systemType: case when st is not null then {uid: st.uid, name: st.name} else null end,
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
	statistics: {
			subsystemsCount: count(DISTINCT subsys),
			minimalSpareParstCount: sys.minimalSpareParstCount,
			sparePartsCoverageSum: sys.sparePartsCoverageSum,
			sp_coverage: sys.sp_coverage
	}
} AS system`
	result.ReturnAlias = "system"
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func GetSystemByEunQuery(eun string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH (ciCategory)<-[:BELONGS_TO_CATEGORY]-(catalogueItem:CatalogueItem)<-[:IS_BASED_ON]-(physicalItem:Item)<-[:CONTAINS_ITEM]-(sys:System) WHERE toLower(physicalItem.eun) = toLower($eun)
	WITH sys, catalogueItem, physicalItem, ciCategory
	OPTIONAL MATCH (sys)-[:HAS_LOCATION]->(loc)  
	OPTIONAL MATCH (sys)-[:HAS_ZONE]->(zone)  
	OPTIONAL MATCH (sys)-[:HAS_SYSTEM_TYPE]->(st)	
	OPTIONAL MATCH (sys)-[:HAS_RESPONSIBLE]->(responsilbe)
	OPTIONAL MATCH (sys)-[:HAS_IMPORTANCE]->(imp)		
	OPTIONAL MATCH (physicalItem)-[:HAS_ITEM_USAGE]->(itemUsage)
	OPTIONAL MATCH (parents{deleted: false})-[:HAS_SUBSYSTEM*1..50]->(sys)
	OPTIONAL MATCH (sys)-[:HAS_SUBSYSTEM*1..50]->(subsys{deleted: false})
	RETURN DISTINCT {  
	uid: sys.uid,
	description: sys.description,
	name: sys.name,
	parentPath: case when parents is not null then reverse(collect(distinct {uid: parents.uid, name: parents.name})) else null end,
	systemCode: sys.systemCode,	
	systemLevel: sys.systemLevel,
	miniImageUrl: split(sys.miniImageUrl, ";"),	
	isTechnologicalUnit: sys.isTechnologicalUnit,
	location: case when loc is not null then {uid: loc.uid, name: loc.name, code: loc.code} else null end,
	zone: case when zone is not null then {uid: zone.uid, name: zone.name} else null end,
	systemType: case when st is not null then {uid: st.uid, name: st.name} else null end,
	responsible: case when responsilbe is not null then {uid: responsilbe.uid, name: responsilbe.lastName + " " + responsilbe.firstName} else null end,
	importance: case when imp is not null then {uid: imp.uid, name: imp.name} else null end,	
	subsystems: case when subsys is not null then collect({uid: subsys.uid, name: subsys.name}) else null end,
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
	statistics: {
			subsystemsCount: count(DISTINCT subsys),
			minimalSpareParstCount: sys.minimalSpareParstCount,
			sparePartsCoverageSum: sys.sparePartsCoverageSum,
			sp_coverage: sys.sp_coverage
	}} AS system;`
	result.ReturnAlias = "system"

	result.Parameters = make(map[string]interface{})
	result.Parameters["eun"] = strings.Trim(eun, " ")
	return result
}

func GetSystemRelationshipsQuery(uid string) (result helpers.DatabaseQuery) {
	result.Query = `
	MATCH(sys:System{uid: $uid, deleted: false})
	MATCH (parents{deleted: false})-[rParent:HAS_SUBSYSTEM]->(sys)	
	return distinct {
		direction: "to",
		foreignSystemName: parents.name,
		relationUid: id(rParent),
		relationTypeCode: "HAS_SUBSYSTEM"
		} as relationships
	UNION
	MATCH(sys:System{uid: $uid, deleted: false})
	MATCH (sys)-[rSubsys:HAS_SUBSYSTEM]->(subsys{deleted: false})	
	return distinct {
		direction: "from",
		foreignSystemName: subsys.name,
		relationUid: id(rSubsys),
		relationTypeCode: "HAS_SUBSYSTEM"
		} as relationships
	UNION
	MATCH(sys:System{uid: $uid, deleted: false})
	MATCH (parents{deleted: false})-[rParent:IS_SPARE_FOR]->(sys)	
	return distinct {
		direction: "to",
		foreignSystemName: parents.name,
		relationUid: id(rParent),
		relationTypeCode: "IS_SPARE_FOR"
		} as relationships
	UNION
	MATCH(sys:System{uid: $uid, deleted: false})
	MATCH (sys)-[rSubsys:IS_SPARE_FOR]->(subsys{deleted: false})	
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

func GetSystemTypeMask(systemTypeUID, facilityCode string) (result helpers.DatabaseQuery) {

	result.Query = `MATCH (st:SystemType{uid: $systemTypeUID}) `

	switch facilityCode {
	case "B":
		result.Query += ` RETURN st.maskB as mask `
	case "A":
		result.Query += ` RETURN st.maskA as mask `
	case "N":
		result.Query += ` RETURN st.maskN as mask `
	default:
		result.Query += ` RETURN "" `
	}

	result.Parameters = make(map[string]interface{})
	result.Parameters["systemTypeUID"] = systemTypeUID

	result.ReturnAlias = "mask"

	return result
}

func GetSystemTypeCode(systemTypeUID string) (result helpers.DatabaseQuery) {

	result.Query = `MATCH (st:SystemType{uid: $systemTypeUID}) RETURN st.code as code`

	result.Parameters = make(map[string]interface{})
	result.Parameters["systemTypeUID"] = systemTypeUID

	result.ReturnAlias = "code"

	return result
}

func GetZoneCode(zoneUID string) (result helpers.DatabaseQuery) {

	result.Query = `
	MATCH (z:Zone{uid: $zoneUID}) 
	WITH z
	OPTIONAL MATCH (pz)-[:HAS_SUBZONE]->(z)
	WITH CASE WHEN pz IS NOT NULL THEN pz.code ELSE z.code END as code
	RETURN code as code `

	result.Parameters = make(map[string]interface{})
	result.Parameters["zoneUID"] = zoneUID

	result.ReturnAlias = "code"

	return result
}

func GetSubZoneCode(zoneUID string) (result helpers.DatabaseQuery) {

	result.Query = `
	OPTIONAL MATCH (z:Zone{uid: $zoneUID}) 
	WHERE ()-[:HAS_SUBZONE]->(z)	
	WITH CASE WHEN z IS NOT NULL THEN z.code ELSE "" END as code	
	RETURN code as code `

	result.Parameters = make(map[string]interface{})
	result.Parameters["zoneUID"] = zoneUID

	result.ReturnAlias = "code"

	return result
}

func GetLocationCode(locationUID string) (result helpers.DatabaseQuery) {

	result.Query = `MATCH (l:Location{uid: $locationUID}) RETURN l.code as code`

	result.Parameters = make(map[string]interface{})
	result.Parameters["locationUID"] = locationUID

	result.ReturnAlias = "code"

	return result
}

func GetZoneName(zoneUID string) (result helpers.DatabaseQuery) {

	result.Query = `MATCH (z:Zone{uid: $zoneUID}) RETURN z.name as name`

	result.Parameters = make(map[string]interface{})
	result.Parameters["zoneUID"] = zoneUID

	result.ReturnAlias = "name"

	return result
}

func GetLocationName(locationUID string) (result helpers.DatabaseQuery) {

	result.Query = `MATCH (l:Location{uid: $locationUID}) RETURN l.name as name`

	result.Parameters = make(map[string]interface{})
	result.Parameters["locationUID"] = locationUID

	result.ReturnAlias = "name"

	return result
}

func GetNewSystemCode(systemCodePrefix string, serialNumberLength int, facilityCode string) (result helpers.DatabaseQuery) {

	result.Query = `
	OPTIONAL MATCH(st:System)-[:BELONGS_TO_FACILITY]->(f{code: $facilityCode}) WHERE st.systemCode STARTS WITH $systemCodePrefix
	WITH st ORDER BY st.systemCode DESC LIMIT 1
	WITH CASE WHEN st IS NOT NULL THEN toInteger(split(st.systemCode, $systemCodePrefix)[1]) + 1 ELSE 1 END as serialNumber
	RETURN $systemCodePrefix +  apoc.text.lpad(toString(serialNumber), $serialNumberLength, "0") as newCode;`

	result.Parameters = make(map[string]interface{})
	result.Parameters["systemCodePrefix"] = systemCodePrefix
	result.Parameters["serialNumberLength"] = serialNumberLength
	result.Parameters["facilityCode"] = facilityCode

	result.ReturnAlias = "newCode"

	return result
}

func GetPhysicalItemPropertiesQuery(physicalItemUID string) (result helpers.DatabaseQuery) {

	result.Query = `
	MATCH(itm:Item{uid: $physicalItemUID})-[:IS_BASED_ON]->(ci)-[:BELONGS_TO_CATEGORY]->(cat)
	WITH itm, cat
	OPTIONAL MATCH(parentCategories)-[:HAS_SUBCATEGORY*1..20]->(cat)
	WITH itm, collect(parentCategories.uid) + cat.uid as categoryUids
	OPTIONAL MATCH (cat)-[:CONTAINS_PHYSICAL_ITEM_PROPERTY]->(prop) WHERE cat.uid in categoryUids
	WITH itm, prop
	OPTIONAL MATCH (prop)-[:IS_PROPERTY_TYPE]->(ptype)
	OPTIONAL MATCH (prop)-[:HAS_UNIT]->(punit)
	OPTIONAL MATCH (itm)-[pv:HAS_PHYSICAL_ITEM_PROPERTY]->(prop)
	RETURN DISTINCT CASE WHEN prop IS NOT NULL THEN {
						value: CASE WHEN pv IS NOT NULL THEN pv.value ELSE null END,
						property: {
						uid: prop.uid,					
						name: prop.name,
						listOfValues: apoc.text.split(case when prop.listOfValues = "" then null else  prop.listOfValues END, ";"),
						defaultValue: prop.defaultValue,
						type: CASE WHEN ptype IS NOT NULL THEN {name: ptype.name, uid: ptype.uid} ELSE null END,
						unit: CASE WHEN punit IS NOT NULL THEN {name: punit.name, uid: punit.uid} ELSE null END					
						}
					} ELSE NULL END as properties ORDER BY properties.property.name;`

	result.Parameters = make(map[string]interface{})
	result.Parameters["physicalItemUID"] = physicalItemUID

	result.ReturnAlias = "properties"

	return result
}

func UpdatePhysicalItemDetailsQuery(physicalItemUID string, details []models.PhysicalItemDetail, userUID string) (result helpers.DatabaseQuery) {

	result.Parameters = make(map[string]interface{})
	result.Parameters["physicalItemUID"] = physicalItemUID
	result.Parameters["lastUpdateBy"] = userUID

	result.Query = `
	MATCH (itm:Item{uid: $physicalItemUID})
	`

	for i, detail := range details {

		propIdxUid := fmt.Sprintf("propUID%v", i)
		propValIdx := fmt.Sprintf("propVal%v", i)

		result.Parameters[propIdxUid] = detail.Property.UID
		result.Parameters[propValIdx] = detail.Value

		result.Query += `
		WITH itm
		MATCH (prop{uid: $` + propIdxUid + `})
		MERGE (itm)-[pv:HAS_PHYSICAL_ITEM_PROPERTY]->(prop)
		SET pv.value = $` + propValIdx
	}

	result.Query += `
	WITH itm
	CREATE(itm)-[:WAS_UPDATED_BY{ at: datetime(), action: "UPDATE" }]->(u{uid: $lastUpdateBy})
	RETURN true as result`

	result.ReturnAlias = "result"

	return result
}

func GetSystemHistoryQuery(systemUID string) (result helpers.DatabaseQuery) {

	result.Query = `
	CALL { 
		MATCH(s:System{uid: $systemUID})
		WITH s
		MATCH(s)-[upd:WAS_UPDATED_BY]->(usr)
		RETURN {uid: apoc.create.uuid(), changedAt: upd.at, changedBy: usr.lastName + " " + usr.firstName , historyType: "GENERAL", action: upd.action} as history
		UNION		
		MATCH(s:System{uid: $systemUID})
		WITH s
		MATCH(s)<-[upd:IS_ORIGINATED_FROM]-(physicalItem)<-[:CONTAINS_ITEM]-(s2)
		WITH s,upd,s2
		MATCH(usr:User{uid: upd.userUid})
		RETURN {uid: apoc.create.uuid(), changedAt: upd.at, changedBy: usr.lastName + " " + usr.firstName , historyType: "ITEM", detail: { systemUid: s2.uid, systemName: s2.name, direction: "IN" }} as history
		UNION
		MATCH(s:System{uid: $systemUID})
		WITH s
		MATCH(s)-[upd:WAS_MOVED_FROM]->(s2:System)
		WITH s,upd,s2
		MATCH(usr:User{uid: upd.userUid})
		RETURN {uid: apoc.create.uuid(), changedAt: upd.at, changedBy: usr.lastName + " " + usr.firstName , historyType: "MOVE" , detail: { systemUid: s2.uid, systemName: s2.name, direction: "OUT" }} as history
		UNION		
		MATCH(s:System{uid: $systemUID})
		MATCH(itm)-[upd:WAS_MOVED_FROM]->(s)
		WITH s,upd,itm
		MATCH(usr:User{uid: upd.userUid})
		RETURN {uid: apoc.create.uuid(), changedAt: upd.at, changedBy: usr.lastName + " " + usr.firstName , historyType: "ITEM_MOVE" , detail: { systemUid: upd.movedToUid, systemName: itm.name, direction: "IN" }} as history
		UNION		
		MATCH(s:System{uid: $systemUID})
		MATCH(itm)-[upd:WAS_MOVED_TO]->(s)
		WITH s,upd,itm
		MATCH(usr:User{uid: upd.userUid})
		RETURN {uid: apoc.create.uuid(), changedAt: upd.at, changedBy: usr.lastName + " " + usr.firstName , historyType: "ITEM_MOVE" , detail: { systemUid: upd.movedFromUid, systemName: itm.name, direction: "OUT" }} as history		
		UNION
		MATCH(s:System{uid: $systemUID})
		WITH s
		MATCH(s)<-[upd:WAS_MOVED_FROM]-(s2:System)
		WITH s,upd,s2
		MATCH(usr:User{uid: upd.userUid})
		RETURN {uid: apoc.create.uuid(), changedAt: upd.at, changedBy: usr.lastName + " " + usr.firstName , historyType: "MOVE", detail: { systemUid: s2.uid, systemName: s2.name, direction: "IN" }} as history
		}
		RETURN history
		ORDER BY history.changedAt DESC;
	`

	result.Parameters = make(map[string]interface{})
	result.Parameters["systemUID"] = systemUID

	result.ReturnAlias = "history"

	return result
}

func GetSystemTypeGroupsQuery(facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH(f:Facility{code: $facilityCode}) MATCH(n:SystemTypeGroup)-[:BELONGS_TO_FACILITY]->(f)
					RETURN { name: n.name, uid: n.uid } as result ORDER BY result.name`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func GetSystemTypesBySystemTypeGroupQuery(systemTypeGroupUid, facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = `
	MATCH(n:SystemTypeGroup{uid: $systemTypeGroupUid})-[:CONTAINS_SYSTEM_TYPE]->(st:SystemType)
	OPTIONAL MATCH (st)-[:HAS_SYSTEM_ATTRIBUTE]->(attr:SystemAttribute)
	RETURN 
	{ name: st.name, 
	  uid: st.uid, 
	  code: st.code,
	  mask: case when $facilityCode = "B" then st.maskB WHEN $facilityCode = "A" THEN st.maskA WHEN $facilityCode = "N" THEN st.maskN END ,
	  systemAttribute: case when attr is not null then { name: attr.name, uid: attr.uid } else null end
	  } as result ORDER BY result.name`

	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["systemTypeGroupUid"] = systemTypeGroupUid
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func DeleteSystemTypeGroupQuery(systemTypeGroupUid string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH(n:SystemTypeGroup{uid: $systemTypeGroupUid}) DETACH DELETE n RETURN true as result`
	result.Parameters = make(map[string]interface{})
	result.Parameters["systemTypeGroupUid"] = systemTypeGroupUid
	result.ReturnAlias = "result"
	return result
}

func GetSystemTypeGroupRelatedNodeLabelsCountQuery(systemTypeGroupUid string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH(grp:SystemTypeGroup{uid: $systemTypeGroupUid})
	WITH grp
	MATCH(grp)-[:CONTAINS_SYSTEM_TYPE]->(st:SystemType)<-[:HAS_SYSTEM_TYPE]-(n)
	RETURN { label: labels(n)[0], count: count(n) } as result;`
	result.Parameters = make(map[string]interface{})
	result.Parameters["systemTypeGroupUid"] = systemTypeGroupUid
	result.ReturnAlias = "result"
	return result
}

func DeleteSystemTypeQuery(systemTypeUid string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH(n:SystemType{uid: $systemTypeUid}) DETACH DELETE n RETURN true as result`
	result.Parameters = make(map[string]interface{})
	result.Parameters["systemTypeUid"] = systemTypeUid
	result.ReturnAlias = "result"
	return result
}

func GetSystemTypeRelatedNodeLabelsCountQuery(systemTypeUid string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH(st:SystemType{uid: $systemTypeUid})
	WITH st
	MATCH(st)<-[:HAS_SYSTEM_TYPE]-(n)
	RETURN { label: labels(n)[0], count: count(n) } as result;`
	result.Parameters = make(map[string]interface{})
	result.Parameters["systemTypeUid"] = systemTypeUid
	result.ReturnAlias = "result"
	return result
}

func CreateSystemTypeGroupQuery(systemTypeGroup *codebookModels.Codebook, facilityCode, userUID string) (result helpers.DatabaseQuery) {
	result.Parameters = make(map[string]interface{})
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["userUID"] = userUID
	result.Parameters["name"] = systemTypeGroup.Name
	systemTypeGroup.UID = uuid.NewString()
	result.Parameters["uid"] = systemTypeGroup.UID

	result.Query = `MATCH(f:Facility{code: $facilityCode}) MATCH(u:User{uid: $userUID})
	CREATE(n:SystemTypeGroup{uid: $uid, name: $name})-[:BELONGS_TO_FACILITY]->(f)
	CREATE(n)-[:WAS_UPDATED_BY{ at: datetime(), action: "CREATE" }]->(u)
	RETURN true as result`

	result.ReturnAlias = "result"

	return result
}

func UpdateSystemTypeGroupQuery(systemTypeGroup *codebookModels.Codebook, userUID string) (result helpers.DatabaseQuery) {
	result.Parameters = make(map[string]interface{})
	result.Parameters["userUID"] = userUID
	result.Parameters["name"] = systemTypeGroup.Name
	result.Parameters["uid"] = systemTypeGroup.UID

	result.Query = `
	MATCH(u:User{uid: $userUID})
	MATCH(n:SystemTypeGroup{uid: $uid})	
	SET n.name = $name
	CREATE(n)-[:WAS_UPDATED_BY{ at: datetime(), action: "UPDATE" }]->(u)
	RETURN true as result`

	result.ReturnAlias = "result"

	return result
}

func CreateSystemTypeQuery(systemType *models.SystemType, facilityCode, userUID, systemTypeGroupUID string) (result helpers.DatabaseQuery) {

	result.Parameters = make(map[string]interface{})
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["userUID"] = userUID
	result.Parameters["name"] = systemType.Name
	result.Parameters["code"] = systemType.Code
	systemType.UID = uuid.NewString()
	result.Parameters["uid"] = systemType.UID
	result.Parameters["mask"] = systemType.Mask
	result.Parameters["systemTypeGroupUID"] = systemTypeGroupUID

	result.Query = `
	MATCH(u:User{uid: $userUID})
	CREATE(st:SystemType{uid: $uid, name: $name, code: $code }) `

	if facilityCode == helpers.FACILITY_CODE_ALPS {
		result.Query += ` SET st.maskA = $mask `
	} else if facilityCode == helpers.FACILITY_CODE_BEAMLINES {
		result.Query += ` SET st.maskB = $mask `
	} else if facilityCode == helpers.FACILITY_CODE_NP {
		result.Query += ` SET st.maskN = $mask `
	}

	if systemType.SystemAttribute != nil {
		result.Query += ` WITH st, u		
		MATCH(attr:SystemAttribute{uid: $attrUID})
		CREATE(st)-[:HAS_SYSTEM_ATTRIBUTE]->(attr) `

		result.Parameters["attrUID"] = systemType.SystemAttribute.UID
	}

	result.Query += ` CREATE(st)-[:WAS_UPDATED_BY{ at: datetime(), action: "CREATE" }]->(u)
	WITH st
	MATCH(grp:SystemTypeGroup{uid: $systemTypeGroupUID})
	CREATE(grp)-[:CONTAINS_SYSTEM_TYPE]->(st)
	RETURN true as result`

	result.ReturnAlias = "result"

	return result
}

func UpdateSystemTypeQuery(systemType *models.SystemType, facilityCode, userUID string) (result helpers.DatabaseQuery) {

	result.Parameters = make(map[string]interface{})
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["userUID"] = userUID
	result.Parameters["name"] = systemType.Name
	result.Parameters["code"] = systemType.Code
	result.Parameters["uid"] = systemType.UID
	result.Parameters["mask"] = systemType.Mask

	result.Query = `
	MATCH(u:User{uid: $userUID})
	MATCH(st:SystemType{uid: $uid }) SET st.name = $name, st.code = $code, `

	if facilityCode == helpers.FACILITY_CODE_ALPS {
		result.Query += ` st.maskA = $mask `
	} else if facilityCode == helpers.FACILITY_CODE_BEAMLINES {
		result.Query += ` st.maskB = $mask `
	} else if facilityCode == helpers.FACILITY_CODE_NP {
		result.Query += ` st.maskN = $mask `
	}

	if systemType.SystemAttribute != nil {
		result.Query += ` WITH st, u
		OPTIONAL MATCH(st)-[r:HAS_SYSTEM_ATTRIBUTE]->(attr:SystemAttribute)
		DELETE r
		WITH st, u
		MATCH(attr:SystemAttribute{uid: $attrUID})
		CREATE(st)-[:HAS_SYSTEM_ATTRIBUTE]->(attr) `

		result.Parameters["attrUID"] = systemType.SystemAttribute.UID
	} else {
		result.Query += ` WITH st, u
		MATCH(st)-[r:HAS_SYSTEM_ATTRIBUTE]->(attr:SystemAttribute)
		DELETE r `
	}

	result.Query += `WITH st, u CREATE(st)-[:WAS_UPDATED_BY{ at: datetime(), action: "UPDATE" }]->(u)
	RETURN true as result`

	result.ReturnAlias = "result"

	return result
}

func GetSystemAttributesCodebookQuery(facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH(f:Facility{code: $facilityCode}) 
					MATCH(attr:SystemAttribute)-[:BELONGS_TO_FACILITY]->(f)
					RETURN { name: attr.name, uid: attr.uid } as result ORDER BY result.name`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func GetEunsQuery(facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = `match(itm:Item)<-[:CONTAINS_ITEM]-(s)-[:BELONGS_TO_FACILITY]->(f:Facility{code: $facilityCode})
	where itm.eun is not null and itm.eun <> ""
	return { eun: itm.eun } as result
	order by result.eun desc`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func SyncSystemLocationByEUNQuery(eun, locationUid, userUID string) (result helpers.DatabaseQuery) {
	result.Parameters = make(map[string]interface{})
	result.Parameters["locationUid"] = locationUid
	result.Parameters["eun"] = eun
	result.Parameters["userUID"] = userUID

	result.Query = `
	MATCH(itm:Item{eun: $eun})<-[:CONTAINS_ITEM]-(s:System)
	WITH s
	OPTIONAL MATCH (s)-[lr:HAS_LOCATION]->(loc)
	DELETE lr
	WITH s
	MATCH (l:Location{uid: $locationUid})
	CREATE (s)-[:HAS_LOCATION]->(l)
	WITH s
	MATCH(u:User{uid: $userUID})
	CREATE(s)-[:WAS_UPDATED_BY{ at: datetime(), action: "UPDATE" }]->(u)
	RETURN true as result`

	result.ReturnAlias = "result"

	return result
}

func GetAllLocationsFlatQuery(facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH(f:Facility{code: $facilityCode}) 
					MATCH(loc:Location)-[:BELONGS_TO_FACILITY]->(f)
					WHERE loc.code IS NOT NULL AND loc.code <> ""
					RETURN {
						name	: loc.name,
						uid		: loc.uid,
						code	: loc.code
					} as result ORDER BY result.code`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func GetAllSystemTypesQuery() (result helpers.DatabaseQuery) {
	result.Query = `MATCH(st:SystemType)
	RETURN {
	uid: st.uid,
	name: st.name,
	code: st.code
	} AS result
	ORDER BY result.code;`
	result.ReturnAlias = "result"
	return result
}

func GetAllZonesQuery(facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH(f:Facility{code: $facilityCode}) 
	MATCH(z:Zone)-[:BELONGS_TO_FACILITY]->(f)
	WHERE NOT ()-[:HAS_SUBZONE]->(z) 
	RETURN {
		name	: z.name,
		uid		: z.uid,
		code	: z.code
	} as result ORDER BY result.code;`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func CreateNewSystemForNewSystemCodeQuery(parentUID, systemCode, systemTypeUID, zoneUID, facilityCode, userUID string) (result helpers.DatabaseQuery) {
	result.Parameters = make(map[string]interface{})
	result.Parameters["parentUID"] = parentUID
	result.Parameters["systemTypeUID"] = systemTypeUID
	result.Parameters["zoneUID"] = zoneUID
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["userUID"] = userUID
	result.Parameters["systemLevel"] = "SUBSYSTEMS_AND_PARTS"
	result.Parameters["systemCode"] = systemCode

	result.Query = `
	MATCH(f:Facility{code: $facilityCode})
	MATCH(u:User{uid: $userUID})
	MATCH(parent:System{uid: $parentUID, deleted: false})-[:BELONGS_TO_FACILITY]->(f)
	MATCH(st:SystemType{uid: $systemTypeUID})
	MATCH(z:Zone{uid: $zoneUID})
	CREATE(s:System{uid: apoc.create.uuid(), name: st.name + " " + $systemCode, systemCode: $systemCode, systemLevel: $systemLevel, isTechnologicalUnit: false, lastUpdateTime: datetime(), lastUpdateBy: u.username, deleted: false})-[:BELONGS_TO_FACILITY]->(f)
	CREATE(s)-[:HAS_SYSTEM_TYPE]->(st)
	CREATE(s)-[:HAS_ZONE]->(z)
	CREATE(parent)-[:HAS_SUBSYSTEM]->(s)
	CREATE(s)-[:WAS_UPDATED_BY{ at: datetime(), action: "CREATE" }]->(u)
	RETURN s.uid as result`

	result.ReturnAlias = "result"

	return result
}

func RecalculateSparePartsCoverageCoeficientQuery() (result helpers.DatabaseQuery) {
	result.Query = `
	match(sp)-[spf:IS_SPARE_FOR]->(o)
    with sp, count(spf) as spfs, collect(spf) as spc
    FOREACH (x IN spc | SET x.coverage = 1.0 / spfs)
	RETURN true as result`
	result.ReturnAlias = "result"
	return result
}

func RecalculateSystemSparePartsCoverageQuery() (result helpers.DatabaseQuery) {
	result.Query = `
	match(s:System) where s.minimalSpareParstCount is not null and s.minimalSpareParstCount > 0
	optional match(sp)-[spf:IS_SPARE_FOR]->(s) 
	with s, coalesce(sum(spf.coverage),0) as coverage
	set s.sparePartsCoverageSum = coverage, s.sp_coverage = coverage/s.minimalSpareParstCount
	RETURN true as result`
	result.ReturnAlias = "result"
	return result
}

func RecalculateSystemSparePartsEmptyCoverageQuery() (result helpers.DatabaseQuery) {
	result.Query = `	
	match(s:System) where s.minimalSpareParstCount is null OR s.minimalSpareParstCount = 0
	set s.sp_coverage = NULL
	RETURN true as result`
	result.ReturnAlias = "result"
	return result
}

func MovePhysicalItemQuery(movementInfo *models.PhysicalItemMovement, userUID string) (result helpers.DatabaseQuery) {
	result.Parameters = make(map[string]interface{})
	result.Parameters["userUID"] = userUID
	result.Parameters["sourceSystemUid"] = movementInfo.SourceSystemUID
	if movementInfo.DestinationSystemUID == "" {
		movementInfo.DestinationSystemUID = uuid.NewString()
	}
	result.Parameters["destinationSystemUid"] = movementInfo.DestinationSystemUID
	result.Parameters["parentSystemUid"] = movementInfo.ParentSystemUID
	result.Parameters["systemName"] = movementInfo.SystemName

	result.Query = `
	MATCH(u:User{uid: $userUID})
	MATCH(source:System{uid: $sourceSystemUid})-[rsource:CONTAINS_ITEM]->(itm:Item)
	WITH source, rsource, itm, u
	DELETE rsource
	WITH source, itm, u `

	if movementInfo.DeleteSourceSystem {
		result.Query += `
		SET source.deleted = true
		CREATE(source)-[:WAS_UPDATED_BY{ at: datetime(), action: "DELETE" }]->(u)
		`
	}

	result.Query += `
	WITH source, itm, u
	MERGE(destination:System{uid: $destinationSystemUid})
	SET destination.name = $systemName,
	 destination.deleted = false,
	 destination.systemLevel = source.systemLevel,
	 destination.lastUpdateTime = datetime(), 
	 destination.lastUpdateBy = u.username
	CREATE(destination)-[:CONTAINS_ITEM]->(itm) `
	if movementInfo.ParentSystemUID != "" {
		result.Query += `WITH destination, itm
		 MATCH(parent:System{uid: $parentSystemUid}) 
		 CREATE(parent)-[:HAS_SUBSYSTEM]->(destination) `
	}

	if movementInfo.Location != nil {
		result.Parameters["locationUid"] = movementInfo.Location.UID
		// if there already is a location relationship, delete it
		result.Query += `
		WITH destination, itm
		OPTIONAL MATCH(destination)-[rLocation:HAS_LOCATION]->(loc)
		DELETE rLocation
		WITH destination, itm
		MATCH(l:Location{uid: $locationUid})
		CREATE(destination)-[:HAS_LOCATION]->(l) `
	}

	if movementInfo.Condition != nil {
		result.Parameters["conditionUid"] = movementInfo.Condition.UID
		// if there already is a condition relationship, delete it
		result.Query += `
		WITH destination, itm
		OPTIONAL MATCH(itm)-[rCondition:HAS_CONDITION_STATUS]->(cond)
		DELETE rCondition
		WITH destination, itm
		MATCH(cond:ItemCondition{uid: $conditionUid})
		CREATE(itm)-[:HAS_CONDITION_STATUS]->(cond) `
	}

	if movementInfo.ItemUsage != nil {
		result.Parameters["itemUsageUid"] = movementInfo.ItemUsage.UID
		// if there already is a item usage relationship, delete it
		result.Query += `
		WITH destination, itm
		OPTIONAL MATCH(itm)-[rItemUsage:HAS_ITEM_USAGE]->(iu)
		DELETE rItemUsage
		WITH destination, itm
		MATCH(iu:ItemUsage{uid: $itemUsageUid})
		CREATE(itm)-[:HAS_ITEM_USAGE]->(iu) `
	}

	result.Query += ` RETURN destination.uid as result`

	result.ReturnAlias = "result"

	return result
}

func HasSystemPhysicalItemQuery(systemUID string) (result helpers.DatabaseQuery) {
	result.Query = `OPTIONAL MATCH(s:System{uid: $systemUID})-[:CONTAINS_ITEM]->(itm:Item) RETURN itm IS NOT NULL as result`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["systemUID"] = systemUID
	return result
}

func CopySystemTypeQuery(sourceSystemUID, destinationSystemUID string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH (source:System {uid: $sourceSystemUID})
					OPTIONAL MATCH (source)-[:HAS_SYSTEM_TYPE]->(st)
					MATCH (destination:System {uid: $destinationSystemUID})
					OPTIONAL MATCH (destination)-[r:HAS_SYSTEM_TYPE]->(std)
					WITH source, st, destination, r
					WHERE st IS NOT NULL // Ensure source has HAS_SYSTEM_TYPE
					DELETE r
					CREATE (destination)-[:HAS_SYSTEM_TYPE]->(st)
					RETURN true as result`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["sourceSystemUID"] = sourceSystemUID
	result.Parameters["destinationSystemUID"] = destinationSystemUID
	return result
}

func CopySystemResponsibleQuery(sourceSystemUID, destinationSystemUID string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH (source:System {uid: $sourceSystemUID})
					OPTIONAL MATCH (source)-[:HAS_RESPONSIBLE]->(r)
					MATCH (destination:System {uid: $destinationSystemUID})
					OPTIONAL MATCH (destination)-[rd:HAS_RESPONSIBLE]->(responsible)
					WITH source, r, destination, rd
					WHERE r IS NOT NULL // Ensure source has HAS_RESPONSIBLE
					DELETE rd
					CREATE (destination)-[:HAS_RESPONSIBLE]->(r)
					RETURN true as result`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["sourceSystemUID"] = sourceSystemUID
	result.Parameters["destinationSystemUID"] = destinationSystemUID
	return result
}

func CopySystemResponsibleTeamQuery(sourceSystemUID, destinationSystemUID string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH (source:System {uid: $sourceSystemUID})
					OPTIONAL MATCH (source)-[:HAS_RESPONSIBLE_TEAM]->(rt)
					MATCH (destination:System {uid: $destinationSystemUID})
					OPTIONAL MATCH (destination)-[rtd:HAS_RESPONSIBLE_TEAM]->(responsibleTeam)
					WITH source, rt, destination, rtd
					WHERE rt IS NOT NULL // Ensure source has HAS_RESPONSIBLE_TEAM
					DELETE rtd
					CREATE (destination)-[:HAS_RESPONSIBLE_TEAM]->(rt)
					RETURN true as result`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["sourceSystemUID"] = sourceSystemUID
	result.Parameters["destinationSystemUID"] = destinationSystemUID
	return result
}

func RecordSystemUpdateHistoryQuery(systemUID, userUID, action string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH (s:System{uid: $systemUID})
					MATCH (u:User{uid: $userUID})
					CREATE (s)-[:WAS_UPDATED_BY{ at: datetime(), action: $action }]->(u)
					RETURN true as result`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["systemUID"] = systemUID
	result.Parameters["userUID"] = userUID
	result.Parameters["action"] = action
	return result
}

func RecordSystemMoveHistoryQuery(systemUID, userUID, action, targetSystemUID string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH (s:System{uid: $systemUID})
					MATCH (u:User{uid: $userUID})
					MATCH (t:System{uid: $targetSystemUID})
					CREATE (s)-[:WAS_MOVED_FROM{ at: datetime(), action: $action, userUid: $userUID }]->(t)
					RETURN true as result`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["systemUID"] = systemUID
	result.Parameters["userUID"] = userUID
	result.Parameters["action"] = action
	result.Parameters["targetSystemUID"] = targetSystemUID
	return result
}

func RecordItemMoveHistoryQuery(userUID, targetSystemUID, sourceSystemUID string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH (u:User{uid: $userUID})
					MATCH (s:System{uid: $sourceSystemUID})-[:CONTAINS_ITEM]->(i:Item)
					MATCH (t:System{uid: $targetSystemUID})
					CREATE (i)-[:WAS_MOVED_TO{ at: datetime(), action: "ITEM_MOVE", userUid: $userUID, movedFromUid: $sourceSystemUID }]->(t)
					CREATE (i)-[:WAS_MOVED_FROM{ at: datetime(), action: "ITEM_MOVE", userUid: $userUID, movedToUid: $targetSystemUID }]->(s)
					RETURN true as result`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["userUID"] = userUID
	result.Parameters["targetSystemUID"] = targetSystemUID
	result.Parameters["sourceSystemUID"] = sourceSystemUID
	return result
}

func SetMissingFacilityToSystems(facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH(s:System) WHERE NOT (s)-[:BELONGS_TO_FACILITY]->(:Facility)
					MATCH(f:Facility{code: $facilityCode})
					CREATE(s)-[:BELONGS_TO_FACILITY]->(f)
					RETURN true as result`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func ReplacePhysicalItemsQuery(movementInfo *models.PhysicalItemMovement, userUID, facilityCode string) (result helpers.DatabaseQuery) {

	result.Query = `
	MATCH(u:User{uid: $userUID})
	MATCH(f:Facility{code: $facilityCode})
	MATCH(sourceSystem:System{uid: $sourceSystemUID})-[rSource:CONTAINS_ITEM]->(sourceItem:Item)
	MATCH(destinationSystem:System{uid: $destinationSystemUID})-[rDest:CONTAINS_ITEM]->(destinationItem:Item)
	MATCH(parentSystem:System{uid: $parentSystemUID})
	WITH sourceSystem, sourceItem, destinationSystem, destinationItem, u, f, parentSystem, rSource, rDest
    CREATE(oldSystem:System{
	uid: $oldSystemUid,
	name: destinationSystem.name, 	
	systemLevel: sourceSystem.systemLevel, 
	lastUpdateTime: datetime(), 
	lastUpdateBy: u.username, 
	deleted: false})-[:BELONGS_TO_FACILITY]->(f)
	CREATE(parentSystem)-[:HAS_SUBSYSTEM]->(oldSystem)
	CREATE(oldSystem)-[:CONTAINS_ITEM]->(destinationItem)
	CREATE(destinationSystem)-[:CONTAINS_ITEM]->(sourceItem)
	DELETE rDest, rSource
	`
	if movementInfo.DeleteSourceSystem {
		result.Query += `
		SET sourceSystem.deleted = true
		CREATE(sourceSystem)-[:WAS_UPDATED_BY{ at: datetime(), action: "DELETE" }]->(u)
		`
	}

	result.Query += ` RETURN $oldSystemUid as result`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["sourceSystemUID"] = movementInfo.SourceSystemUID
	result.Parameters["destinationSystemUID"] = movementInfo.DestinationSystemUID
	result.Parameters["parentSystemUID"] = movementInfo.ParentSystemUID
	result.Parameters["oldSystemUid"] = uuid.NewString()
	result.Parameters["userUID"] = userUID
	result.Parameters["facilityCode"] = facilityCode
	return result
}

func MoveSystemsQuery(movementInfo *models.SystemsMovement, userUID string) (result helpers.DatabaseQuery) {

	result.Query = `
	MATCH(newParent:System{uid: $newParentUID})
	MATCH(oldParent:System)-[rOld:HAS_SUBSYSTEM]->(s:System) WHERE s.uid in $systemsToMoveUids
	CREATE(newParent)-[:HAS_SUBSYSTEM]->(s)
	DELETE rOld
	CREATE (s)-[:WAS_MOVED_FROM { at: datetime(), userUid: $userUID }]->(oldParent)
	RETURN DISTINCT newParent.uid as result`

	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["newParentUID"] = movementInfo.TargetParentSystemUid
	result.Parameters["systemsToMoveUids"] = movementInfo.SystemsToMoveUids
	result.Parameters["userUID"] = userUID
	return result
}

func CreateNewSystemFromJiraQuery(request *models.JiraSystemImportRequest, facilityCode string, userUID string) (result helpers.DatabaseQuery) {
	result.Parameters = make(map[string]interface{})
	result.Parameters["facilityCode"] = facilityCode
	result.Parameters["name"] = request.Name
	result.Parameters["systemCode"] = request.Code
	result.Parameters["description"] = request.Description
	result.Parameters["lastUpdateBy"] = userUID
	result.Parameters["zoneUID"] = request.ZoneUID
	result.Parameters["systemTypeUID"] = request.SystemTypeUID
	result.Parameters["parentSystemUID"] = request.ParentSystemUID
	result.Parameters["linkUrl"] = request.LinkUrl
	result.Parameters["linkName"] = request.LinkName

	result.Query = `
	MATCH(f:Facility{code: $facilityCode}) 
	MATCH(u:User{uid: $lastUpdateBy}) 
	MATCH(z:Zone{uid: $zoneUID})-[:BELONGS_TO_FACILITY]->(f)
	MATCH(st:SystemType{uid: $systemTypeUID})
	MATCH(parent:System{uid: $parentSystemUID, deleted: false})-[:BELONGS_TO_FACILITY]->(f)
	WITH DISTINCT f, u, z, st, parent
	
	MERGE(s:System{systemCode: $systemCode})
	ON CREATE SET 
		s.uid = apoc.create.uuid(),
		s.name = $name,
    s.description = $description,
		s.systemLevel = "SUBSYSTEMS_AND_PARTS",
		s.deleted = false,
		s.lastUpdateTime = datetime(),
		s.lastUpdatedBy = u.lastName + " " + u.firstName
	WITH DISTINCT s, f, z, st, parent, u
	
	MERGE(s)-[:BELONGS_TO_FACILITY]->(f)
	MERGE(s)-[:HAS_ZONE]->(z)
	MERGE(s)-[:HAS_SYSTEM_TYPE]->(st)
	MERGE(parent)-[:HAS_SUBSYSTEM]->(s)
	
	WITH DISTINCT s, u
	CREATE(s)-[:WAS_UPDATED_BY{ at: datetime(), action: "INSERT" }]->(u)
	
	// Create file link for JIRA if URL is provided
	WITH DISTINCT s
	FOREACH (ignoreMe IN CASE WHEN $linkUrl <> "" AND $linkName <> "" THEN [1] ELSE [] END |
		CREATE(fl:FileLink{ 
			uid: apoc.create.uuid(), 
			name: $linkName, 
			url: $linkUrl, 
			tags: "jira" })
		CREATE(s)-[:HAS_FILE_LINK{createdAt: datetime()}]->(fl)
	)
	
	RETURN s.uid as result`

	result.ReturnAlias = "result"

	return result
}

func checkSystemCodeExistsQuery(systemCode string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH (s:System{deleted: false}) WHERE toLower(s.systemCode) = $systemCode RETURN count(s) > 0 as result`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["systemCode"] = strings.ToLower(systemCode)
	return result
}
