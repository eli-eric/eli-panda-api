package systemsService

import (
	"panda/apigateway/helpers"
)

// employeeNameCypher renders an Employee's display name: fullName, falling back to "last first".
func employeeNameCypher(alias string) string {
	return `coalesce(` + alias + `.fullName, trim(coalesce(` + alias + `.lastName, "") + " " + coalesce(` + alias + `.firstName, "")))`
}

// GetSystemDetailQuery returns the full system detail consumed by GET/PATCH /v1/system/{uid}.
// Pattern comprehensions keep every optional relation at one row, so no DISTINCT/aggregation
// cartesian products can inflate counts.
func GetSystemDetailQuery(uid string, facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = `
	MATCH (sys:System{uid: $uid, deleted: false})-[:BELONGS_TO_FACILITY]->(:Facility{code: $facilityCode})
	OPTIONAL MATCH (parent:System)-[:HAS_SUBSYSTEM]->(sys)
	OPTIONAL MATCH (sys)-[:CONTAINS_ITEM]->(itm:Item)
	OPTIONAL MATCH (itm)-[:IS_BASED_ON]->(ci:CatalogueItem)
	CALL {
		WITH sys
		OPTIONAL MATCH fullPath = (root:System{deleted: false})-[:HAS_SUBSYSTEM*1..50]->(sys)
		WHERE NOT (root)<-[:HAS_SUBSYSTEM]-()
		RETURN CASE WHEN fullPath IS NULL THEN [] ELSE
			[n IN nodes(fullPath)[..-1] | {uid: n.uid, name: n.name, systemLevel: n.systemLevel}] END AS parentPath
		LIMIT 1
	}
	RETURN {
		uid: sys.uid,
		name: sys.name,
		description: sys.description,
		systemCode: sys.systemCode,
		systemAlias: sys.systemAlias,
		systemLevel: sys.systemLevel,
		lastUpdateTime: sys.lastUpdateTime,
		parentUid: parent.uid,
		parentPath: parentPath,
		systemType: head([(sys)-[:HAS_SYSTEM_TYPE]->(x) | {uid: x.uid, name: x.name, code: x.code}]),
		location: head([(sys)-[:HAS_LOCATION]->(x) | {uid: x.uid, name: x.name, code: x.code}]),
		zone: head([(sys)-[:HAS_ZONE]->(x) | {uid: x.uid, name: x.name, code: x.code}]),
		importance: head([(sys)-[:HAS_IMPORTANCE]->(x) | {uid: x.uid, name: x.name, code: x.code}]),
		attribute: head([(sys)-[:HAS_SYSTEM_ATTRIBUTE]->(x) | {uid: x.uid, name: x.name}]),
		responsible: head([(sys)-[:HAS_RESPONSIBLE]->(x) | {uid: x.uid, name: ` + employeeNameCypher("x") + `}]),
		responsibleTeam: head([(sys)-[:HAS_RESPONSIBLE_TEAM]->(x) | {uid: x.uid, name: x.name, code: x.code}]),
		owner: head([(sys)-[:HAS_OWNER]->(x) | {uid: x.uid, name: ` + employeeNameCypher("x") + `}]),
		operators: apoc.coll.sortMaps([(sys)-[:HAS_OPERATOR]->(x) | {uid: x.uid, name: ` + employeeNameCypher("x") + `}], "^name"),
		maintainedBy: apoc.coll.sortMaps([(sys)-[:IS_MAINTAINED_BY]->(x) | {uid: x.uid, name: ` + employeeNameCypher("x") + `}], "^name"),
		miniImageUrl: CASE WHEN sys.miniImageUrl IS NULL THEN null ELSE split(sys.miniImageUrl, ";") END,
		minimalSpareParstCount: sys.minimalSpareParstCount,
		sparePartsCoverageSum: sys.sparePartsCoverageSum,
		sp_coverage: sys.sp_coverage,
		physicalItem: CASE WHEN itm IS NULL THEN null ELSE {
			uid: itm.uid,
			name: coalesce(itm.name, ci.name, ""),
			eun: itm.eun,
			serialNumber: itm.serialNumber,
			notes: itm.notes,
			price: toFloat(itm.price),
			currency: itm.currency,
			printEUN: itm.printEUN,
			itemUsage: head([(itm)-[:HAS_ITEM_USAGE]->(x) | {uid: x.uid, name: x.name, code: x.code}]),
			conditionStatus: head([(itm)-[:HAS_CONDITION_STATUS]->(x) | {uid: x.uid, name: x.name, code: x.code}]),
			catalogueItem: CASE WHEN ci IS NULL THEN null ELSE {
				uid: ci.uid,
				name: ci.name,
				catalogueNumber: ci.catalogueNumber,
				category: head([(ci)-[:BELONGS_TO_CATEGORY]->(x) | {uid: x.uid, name: x.name}]),
				supplier: head([(ci)-[:HAS_SUPPLIER]->(x) | {uid: x.uid, name: x.name}])
			} END
		} END,
		statistics: {
			subsystemsCount: size([(sys)-[:HAS_SUBSYSTEM*1..50]->(x:System{deleted: false}) | x]),
			directSubsystemsCount: size([(sys)-[:HAS_SUBSYSTEM]->(x:System{deleted: false}) | x]),
			sparePartsCount: size([(sys)<-[:IS_SPARE_FOR]-(x:System{deleted: false}) | x]),
			sparesForCount: size([(sys)-[:IS_SPARE_FOR]->(x:System{deleted: false}) | x]),
			minimalSpareParstCount: sys.minimalSpareParstCount,
			sparePartsCoverageSum: sys.sparePartsCoverageSum,
			sp_coverage: sys.sp_coverage
		}
	} AS system
	LIMIT 1`

	result.ReturnAlias = "system"
	result.Parameters = map[string]interface{}{
		"uid":          uid,
		"facilityCode": facilityCode,
	}
	return result
}

// spareSystemRefCypher projects a spare-related system for GET /v1/system/{uid}/spare-parts.
func spareSystemRefCypher(alias string) string {
	return `{
			uid: ` + alias + `.uid,
			name: ` + alias + `.name,
			systemCode: ` + alias + `.systemCode,
			systemLevel: ` + alias + `.systemLevel,
			parentPath: coalesce(head([p = (root:System{deleted: false})-[:HAS_SUBSYSTEM*1..50]->(` + alias + `) WHERE NOT (root)<-[:HAS_SUBSYSTEM]-() |
				[n IN nodes(p)[..-1] | {uid: n.uid, name: n.name}]]), []),
			location: head([(` + alias + `)-[:HAS_LOCATION]->(x) | {uid: x.uid, name: x.name, code: x.code}]),
			physicalItem: head([(` + alias + `)-[:CONTAINS_ITEM]->(itm:Item) | {
				uid: itm.uid,
				eun: itm.eun,
				itemUsage: head([(itm)-[:HAS_ITEM_USAGE]->(x) | {uid: x.uid, name: x.name, code: x.code}]),
				catalogueItem: head([(itm)-[:IS_BASED_ON]->(ci) | {uid: ci.uid, catalogueNumber: ci.catalogueNumber}])
			}])
		}`
}

func GetSystemSparePartsQuery(uid string, facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = `
	MATCH (sys:System{uid: $uid, deleted: false})-[:BELONGS_TO_FACILITY]->(:Facility{code: $facilityCode})
	RETURN {
		minimalSpareParstCount: sys.minimalSpareParstCount,
		sparePartsCoverageSum: sys.sparePartsCoverageSum,
		sp_coverage: sys.sp_coverage,
		spareParts: apoc.coll.sortMaps([(sys)<-[r:IS_SPARE_FOR]-(other:System{deleted: false}) |
			{coverage: r.coverage, sortName: other.name, system: ` + spareSystemRefCypher("other") + `}], "^sortName"),
		spareFor: apoc.coll.sortMaps([(sys)-[r:IS_SPARE_FOR]->(other:System{deleted: false}) |
			{coverage: r.coverage, sortName: other.name, system: ` + spareSystemRefCypher("other") + `}], "^sortName")
	} AS result`

	result.ReturnAlias = "result"
	result.Parameters = map[string]interface{}{
		"uid":          uid,
		"facilityCode": facilityCode,
	}
	return result
}

// rawPropertyCypher projects a HAS_CATALOGUE_PROPERTY edge (alias rel) to its property node (prop).
const rawPropertyCypher = `{
			uid: prop.uid,
			name: prop.name,
			type: head([(prop)-[:IS_PROPERTY_TYPE]->(t) | t.name]),
			unit: head([(prop)-[:HAS_UNIT]->(u) | u.name]),
			value: CASE WHEN rel.value IS NULL THEN null ELSE toString(rel.value) END,
			groups: [(prop)<-[:CONTAINS_PROPERTY]-(g) | {uid: g.uid, name: g.name}]
		}`

// GetPhysicalItemRawPropertiesQuery reads the catalogue property values of the item's catalogue
// item plus the property values of every service item applied to it. Override resolution
// happens in Go (resolveCatalogueProperties).
func GetPhysicalItemRawPropertiesQuery(itemUID, facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = `
	MATCH (itm:Item{uid: $itemUID})<-[:CONTAINS_ITEM]-(:System)-[:BELONGS_TO_FACILITY]->(:Facility{code: $facilityCode})
	OPTIONAL MATCH (itm)-[:IS_BASED_ON]->(ci:CatalogueItem)
	RETURN {
		catalogue: CASE WHEN ci IS NULL THEN [] ELSE [(ci)-[rel:HAS_CATALOGUE_PROPERTY]->(prop) | ` + rawPropertyCypher + `] END,
		serviceItems: [(itm)-[sb:IS_SERVICED_BY]->(si:ServiceItem) | {
			created: sb.created,
			properties: [(si)-[rel:HAS_CATALOGUE_PROPERTY]->(prop) | ` + rawPropertyCypher + `]
		}]
	} AS result`

	result.ReturnAlias = "result"
	result.Parameters = map[string]interface{}{"itemUID": itemUID, "facilityCode": facilityCode}
	return result
}

func GetPhysicalItemOrdersQuery(itemUID, facilityCode string) (result helpers.DatabaseQuery) {
	result.Query = `
	MATCH (itm:Item{uid: $itemUID})<-[:CONTAINS_ITEM]-(:System)-[:BELONGS_TO_FACILITY]->(:Facility{code: $facilityCode})
	RETURN {
		order: head([(o:Order)-[ol:HAS_ORDER_LINE]->(itm) | {uid: o.uid, name: o.name, orderDate: o.orderDate, isDelivered: ol.isDelivered}]),
		serviceItems: apoc.coll.sortMaps([(itm)-[sb:IS_SERVICED_BY]->(si:ServiceItem) | {
			uid: si.uid,
			name: si.name,
			isDelivered: coalesce(si.isDelivered, false),
			created: sb.created,
			sortCreated: coalesce(toString(sb.created), ""),
			order: head([(o:Order)-[:HAS_SERVICE_LINE]->(si) | {uid: o.uid, name: o.name, orderDate: o.orderDate}])
		}], "sortCreated")
	} AS result`

	result.ReturnAlias = "result"
	result.Parameters = map[string]interface{}{"itemUID": itemUID, "facilityCode": facilityCode}
	return result
}
