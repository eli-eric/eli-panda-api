package catalogueService

import (
	"encoding/json"
	"fmt"
	"panda/apigateway/helpers"
	"panda/apigateway/services/catalogue-service/models"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func CatalogueItemsFiltersPrepareQuery(search string, categoryUid string, filering *[]helpers.ColumnFilter, skip, limit int) (result helpers.DatabaseQuery) {

	result.Parameters = make(map[string]interface{})
	result.Parameters["searchText"] = strings.TrimSpace(strings.ToLower(search))
	result.Parameters["skip"] = skip
	result.Parameters["limit"] = limit
	result.Parameters["categoryUid"] = categoryUid

	// category filter
	filterCategory := helpers.GetFilterValueCodebook(filering, "category")
	if filterCategory != nil {
		categoryUid = filterCategory.UID
		result.Parameters["categoryUid"] = categoryUid
	}

	// Start optimized query with CALL subquery
	result.Query = `CALL {
		WITH toLower($searchText) AS q`

	if categoryUid == "" {
		result.Query += `
		// Start from items & category for all categories
		MATCH (cat:CatalogueCategory)<-[:BELONGS_TO_CATEGORY]-(itm)`
	} else {
		result.Query += `
		// Start from specific category and its subcategories
		MATCH(category:CatalogueCategory)
		WHERE $categoryUid = '' OR category.uid = $categoryUid
		OPTIONAL MATCH(category)-[:HAS_SUBCATEGORY*1..15]->(subs)
		WITH collect(subs.uid) + category.uid AS categoryUids, q
		MATCH (itm:CatalogueItem)-[:BELONGS_TO_CATEGORY]->(cat)
		WHERE cat.uid IN categoryUids`
	}

	result.Query += `
		OPTIONAL MATCH (itm)-[:HAS_SUPPLIER]->(supp)
		OPTIONAL MATCH (itm)-[pv:HAS_CATALOGUE_PROPERTY]->(:CatalogueProperty)

		WITH DISTINCT itm, cat, q,
			 coalesce(toLower(itm.name), '') AS iname,
			 coalesce(toLower(itm.description), '') AS idesc,
			 coalesce(toLower(itm.catalogueNumber), '') AS icn,
			 coalesce(toLower(supp.name), '') AS sname,
			 coalesce(toLower(toString(pv.value)), '') AS pval

		WHERE q = '' OR iname CONTAINS q OR idesc CONTAINS q OR icn CONTAINS q OR sname CONTAINS q OR pval CONTAINS q `

	if filering != nil && len(*filering) > 0 {

		//catalogue name
		if filterVal := helpers.GetFilterValueString(filering, "name"); filterVal != nil {

			result.Query += ` WITH itm, cat, q, iname, idesc, icn, sname, pval  `
			result.Query += ` WHERE iname CONTAINS $filterName `
			result.Parameters["filterName"] = strings.ToLower(*filterVal)
		}

		//catalogue number
		if filterVal := helpers.GetFilterValueString(filering, "catalogueNumber"); filterVal != nil {

			result.Query += ` WITH itm, cat, q, iname, idesc, icn, sname, pval  `
			result.Query += ` WHERE icn CONTAINS $filterCatalogueNumber `
			result.Parameters["filterCatalogueNumber"] = strings.ToLower(*filterVal)
		}

		//supplier
		if filterVal := helpers.GetFilterValueString(filering, "supplier"); filterVal != nil {

			result.Query += ` WITH itm, cat, q, iname, idesc, icn, sname, pval  `
			result.Query += ` WHERE sname CONTAINS $filterSupplier `
			result.Parameters["filterSupplier"] = strings.ToLower(*filterVal)
		}

		//manufacturerUrl
		if filterVal := helpers.GetFilterValueString(filering, "manufacturerUrl"); filterVal != nil {

			result.Query += ` WITH itm, cat, q, iname, idesc, icn, sname, pval  `
			result.Query += ` WHERE toLower(itm.manufacturerUrl) CONTAINS $filterManufacturerUrl `
			result.Parameters["filterManufacturerUrl"] = strings.ToLower(*filterVal)
		}

		//description
		if filterVal := helpers.GetFilterValueString(filering, "description"); filterVal != nil {

			result.Query += ` WITH itm, cat, q, iname, idesc, icn, sname, pval  `
			result.Query += ` WHERE idesc CONTAINS $filterDescription `
			result.Parameters["filterDescription"] = strings.ToLower(*filterVal)
		}

		for i, filter := range *filering {
			if filter.Type == "" {
				continue
			}

			if filter.Type == "text" {
				if filterPropvalue := helpers.GetFilterValueString(filering, filter.Id); filterPropvalue != nil {
					result.Query += ` WITH itm, cat, q, iname, idesc, icn, sname, pval `
					result.Query += fmt.Sprintf(` MATCH(itm)-[pv%v:HAS_CATALOGUE_PROPERTY]->(prop%v{uid: $propUID%v}) WHERE toLower(pv%v.value) CONTAINS $propFilterVal%v `, i, i, i, i, i)

					result.Parameters[fmt.Sprintf("propUID%v", i)] = filter.Id
					result.Parameters[fmt.Sprintf("propFilterVal%v", i)] = strings.ToLower(*filterPropvalue)
				}
			} else if filter.Type == "list" {
				if filterPropvalue := helpers.GetFilterValueListString(filering, filter.Id); filterPropvalue != nil {
					result.Query += ` WITH itm, cat, q, iname, idesc, icn, sname, pval `
					result.Query += fmt.Sprintf(` MATCH(itm)-[pv%v:HAS_CATALOGUE_PROPERTY]->(prop%v{uid: $propUID%v}) WHERE pv%v.value IN $propFilterVal%v `, i, i, i, i, i)

					result.Parameters[fmt.Sprintf("propUID%v", i)] = filter.Id
					result.Parameters[fmt.Sprintf("propFilterVal%v", i)] = filterPropvalue
				}
			} else if filter.Type == "number" {
				if filterPropvalue := helpers.GetFilterValueRangeFloat64(filering, filter.Id); filterPropvalue != nil {
					result.Query += ` WITH itm, cat, q, iname, idesc, icn, sname, pval `
					result.Query += fmt.Sprintf(` MATCH(itm)-[pv%v:HAS_CATALOGUE_PROPERTY]->(prop%v{uid: $propUID%v}) WHERE ($propFilterValFrom%v IS NULL OR toFloat(pv%v.value) >= $propFilterValFrom%v) AND ($propFilterValTo%v IS NULL OR toFloat(pv%v.value) <= $propFilterValTo%v) `, i, i, i, i, i, i, i, i, i)

					result.Parameters[fmt.Sprintf("propUID%v", i)] = filter.Id
					result.Parameters[fmt.Sprintf("propFilterValFrom%v", i)] = filterPropvalue.Min
					result.Parameters[fmt.Sprintf("propFilterValTo%v", i)] = filterPropvalue.Max
				}
			} else if filter.Type == "range" {
				if filterPropvalue := helpers.GetFilterValueRangeFloat64(filering, filter.Id); filterPropvalue != nil {
					result.Query += ` WITH itm, cat, q, iname, idesc, icn, sname, pval `
					result.Query += fmt.Sprintf(` MATCH(itm)-[pv%v:HAS_CATALOGUE_PROPERTY]->(prop%v{uid: $propUID%v})
					WITH itm, cat, q, iname, idesc, icn, sname, pval, apoc.convert.fromJsonMap(pv%v.value) as jsonValue%v
					WHERE jsonValue%v IS NOT NULL AND ((toFloat(jsonValue%v.min) <= $propFilterValFrom%v - $propFilterValTo%v) 
					AND (toFloat(jsonValue%v.max) >= $propFilterValFrom%v + $propFilterValTo%v)) `, i, i, i, i, i, i, i, i, i, i, i, i)

					result.Parameters[fmt.Sprintf("propUID%v", i)] = filter.Id
					result.Parameters[fmt.Sprintf("propFilterValFrom%v", i)] = filterPropvalue.Min
					if filterPropvalue.Max == nil {
						filterPropvalue.Max = &plusMinusZero
					}
					result.Parameters[fmt.Sprintf("propFilterValTo%v", i)] = filterPropvalue.Max
				}
			}
		}

	}

	return result
}

var plusMinusZero float64 = 0

// Helper function to build filter conditions for catalogue items queries
func buildCatalogueItemsFilterConditions(filering *[]helpers.ColumnFilter) (filterConditions []string, parameters map[string]interface{}) {
	parameters = make(map[string]interface{})
	filterConditions = []string{}

	if filering == nil || len(*filering) == 0 {
		return filterConditions, parameters
	}

	//catalogue name
	if filterVal := helpers.GetFilterValueString(filering, "name"); filterVal != nil {
		filterConditions = append(filterConditions, "iname CONTAINS $filterName")
		parameters["filterName"] = strings.ToLower(*filterVal)
	}

	//catalogue number
	if filterVal := helpers.GetFilterValueString(filering, "catalogueNumber"); filterVal != nil {
		filterConditions = append(filterConditions, "icn CONTAINS $filterCatalogueNumber")
		parameters["filterCatalogueNumber"] = strings.ToLower(*filterVal)
	}

	//supplier
	if filterVal := helpers.GetFilterValueString(filering, "supplier"); filterVal != nil {
		filterConditions = append(filterConditions, "sname CONTAINS $filterSupplier")
		parameters["filterSupplier"] = strings.ToLower(*filterVal)
	}

	//description
	if filterVal := helpers.GetFilterValueString(filering, "description"); filterVal != nil {
		filterConditions = append(filterConditions, "idesc CONTAINS $filterDescription")
		parameters["filterDescription"] = strings.ToLower(*filterVal)
	}

	return filterConditions, parameters
}

// Get catalogue items with pagination and filters
func CatalogueItemsFiltersPaginationQuery(search string, categoryUid string, skip int, limit int, filering *[]helpers.ColumnFilter, sorting *[]helpers.Sorting) (result helpers.DatabaseQuery) {

	result = CatalogueItemsFiltersPrepareQuery(search, categoryUid, filering, skip, limit)

	// Add filters if any
	filterConditions, filterParams := buildCatalogueItemsFilterConditions(filering)
	if len(filterConditions) > 0 {
		// Add filter parameters to the result
		for key, value := range filterParams {
			result.Parameters[key] = value
		}
		result.Query += " AND (" + strings.Join(filterConditions, " AND ") + ")"
	}

	// Add latest update subquery and pagination
	result.Query += `

		// Latest update per item via ordered subquery
		CALL {
			WITH itm
			OPTIONAL MATCH (itm)-[u:WAS_UPDATED_BY]->(usr)
			RETURN u.at AS lastUpdateTime, usr.lastName + ' ' + usr.firstName AS lastUpdateUser
			ORDER BY lastUpdateTime DESC
			LIMIT 1
		}

		// Page on the reduced set
		RETURN itm, cat, lastUpdateTime, lastUpdateUser`

	// Add sorting
	if sorting != nil && len(*sorting) > 0 {
		result.Query += " ORDER BY "
		for ids, sort := range *sorting {
			sortName := sort.ID
			// handle special cases
			if sortName == "partNumber" {
				sortName = "itm.catalogueNumber"
			} else if sortName == "catalogueNumber" {
				sortName = "itm.catalogueNumber"
			} else if sortName == "categoryName" {
				sortName = "cat.name"
			} else if sortName == "name" {
				sortName = "itm.name"
			} else if sortName == "lastUpdateTime" {
				sortName = "lastUpdateTime"
			}

			if ids == len(*sorting)-1 {
				result.Query += fmt.Sprintf(" %s %s ", sortName, helpers.GetSortingDirectionString(sort.DESC))
			} else {
				result.Query += fmt.Sprintf(" %s %s, ", sortName, helpers.GetSortingDirectionString(sort.DESC))
			}
		}
	} else {
		result.Query += " ORDER BY lastUpdateTime DESC "
	}

	result.Query += `
		SKIP toInteger($skip)
		LIMIT toInteger($limit)
	}

	// Fan out details only for paged items
	OPTIONAL MATCH (itm)-[propVal:HAS_CATALOGUE_PROPERTY]->(prop)
	OPTIONAL MATCH (prop)-[:HAS_UNIT]->(unit)
	OPTIONAL MATCH (prop)-[:IS_PROPERTY_TYPE]->(propType)
	OPTIONAL MATCH (group)-[:CONTAINS_PROPERTY]->(prop)
	OPTIONAL MATCH (itm)-[:HAS_SUPPLIER]->(supp)

	WITH itm, cat, propType, supp, prop, group.name AS groupName,
		 toString(propVal.value) AS value, unit, lastUpdateTime, lastUpdateUser

	RETURN {
		uid: itm.uid,
		name: itm.name,
		catalogueNumber: itm.catalogueNumber,
		miniImageUrl: split(itm.miniImageUrl, ';'),
		description: itm.description,
		category: { uid: cat.uid, name: cat.name },
		supplier: CASE WHEN supp IS NOT NULL THEN { uid: supp.uid, name: supp.name } ELSE NULL END,
		manufacturerUrl: itm.manufacturerUrl,
		lastUpdateTime: lastUpdateTime,
		lastUpdateBy: lastUpdateUser,
		details: CASE WHEN COUNT(prop) > 0 THEN collect(DISTINCT {
			property: {
				uid: prop.uid,
				name: prop.name,
				listOfValues: CASE
					WHEN prop.listOfValues IS NOT NULL AND prop.listOfValues <> ''
					THEN apoc.text.split(prop.listOfValues, ';') ELSE NULL END,
				type: CASE WHEN propType IS NOT NULL THEN { uid: propType.uid, name: propType.name, code: propType.code } ELSE NULL END,
				unit: CASE WHEN unit IS NOT NULL THEN { uid: unit.uid, name: unit.name } ELSE NULL END
			},
			propertyGroup: groupName,
			value: value
		}) ELSE NULL END
	} AS items
	`

	result.ReturnAlias = "items"

	return result
}

func CatalogueItemsFiltersTotalCountQuery(search string, categoryUid string, filering *[]helpers.ColumnFilter) (result helpers.DatabaseQuery) {

	result = CatalogueItemsFiltersPrepareQuery(search, categoryUid, filering, 0, 0)

	// Add filters if any
	filterConditions, filterParams := buildCatalogueItemsFilterConditions(filering)
	if len(filterConditions) > 0 {
		// Add filter parameters to the result
		for key, value := range filterParams {
			result.Parameters[key] = value
		}
		result.Query += " AND (" + strings.Join(filterConditions, " AND ") + ")"
	}

	// Close the CALL subquery and count distinct items
	result.Query += `
		RETURN DISTINCT itm
	}
	RETURN count(itm) as itemsCount`

	result.ReturnAlias = "itemsCount"

	return result
}

func CatalogueCategoriesByParentPathQuery(parentPath string) (result helpers.DatabaseQuery) {

	result.Query = `match(category:CatalogueCategory)
	with category
	optional match(parent)-[:HAS_SUBCATEGORY*1..50]->(category) 
	with category, apoc.text.join(reverse(collect(parent.code)),"/") as parentPath
	where parentPath = $parentPath
	return {uid:category.uid,code:category.code, name:category.name,parentPath: parentPath, miniImageUrl: split(category.miniImageUrl, ";")} as categories order by id(category)`

	result.ReturnAlias = "categories"

	result.Parameters = make(map[string]interface{})
	result.Parameters["parentPath"] = parentPath

	return result
}

func CatalogueCategoryImageByUidQuery(uid string) (result helpers.DatabaseQuery) {

	result.Query = `match(category:CatalogueCategory{uid: $uid})
	
	return category.image as image`

	result.ReturnAlias = "image"

	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid

	return result
}

func CatalogueItemImageByUidQuery(uid string) (result helpers.DatabaseQuery) {

	result.Query = `match(r:CatalogueItem{uid: $uid})
	
	return r.image as image`

	result.ReturnAlias = "image"

	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid

	return result
}

func DeleteCatalogueCategoryByUidQuery(uid string) (result helpers.DatabaseQuery) {

	result.Query = `MATCH(p:CatalogueCategory) WHERE p.uid = $uid WITH p
	OPTIONAL MATCH (p)-[:HAS_SUBCATEGORY*1..50]->(child)
	WITH p, child, p.uid as uid
	detach delete p, child`

	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid

	return result
}

func CatalogueItemWithDetailsByUidQuery(uid string) (result helpers.DatabaseQuery) {

	result.Query = `match(itm:CatalogueItem{uid: $uid})-[:BELONGS_TO_CATEGORY]->(cat)			
	OPTIONAL MATCH(itm)-[:HAS_SUPPLIER]->(supp)
	OPTIONAL MATCH(itm)-[propVal:HAS_CATALOGUE_PROPERTY]->(prop)
	OPTIONAL MATCH(prop)-[:HAS_UNIT]->(unit)
	OPTIONAL MATCH(prop)-[:IS_PROPERTY_TYPE]->(propType)
	OPTIONAL MATCH(group)-[:CONTAINS_PROPERTY]->(prop)
	WITH itm, cat, prop, propType, supp, unit,  group.name as groupName, toString(propVal.value) as value
	RETURN {
	uid: itm.uid,
	name: itm.name,
	lastUpdateTime: itm.lastUpdateTime,
	catalogueNumber: itm.catalogueNumber,
	description: itm.description,
	category: {uid: cat.uid, name: cat.name},
	supplier: case when supp is not null then {uid: supp.uid, name: supp.name} else null end,
	manufacturerUrl: itm.manufacturerUrl,	
	miniImageUrl: split(itm.miniImageUrl, ";"),
	details: case when count(prop) > 0 then collect(DISTINCT { 
					property:{
						uid: prop.uid,
						name: prop.name, 
						listOfValues: case when prop.listOfValues is not null and prop.listOfValues <> "" then apoc.text.split(prop.listOfValues, ";") else null end,
						type: { uid: propType.uid, name: propType.name, code: propType.code},
						unit: case when unit is not null then { uid: unit.uid, name: unit.name } else null end 
						},
						propertyGroup: groupName, 
						value: value}) else null end
	} as catalogueItem;`

	result.ReturnAlias = "catalogueItem"

	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid

	return result
}

func CatalogueCategoryWithDetailsQuery(uid string) (result helpers.DatabaseQuery) {

	result.Query = `MATCH(category:CatalogueCategory{uid:$uid})
	OPTIONAL MATCH(category)-[:HAS_GROUP]->(group)-[:CONTAINS_PROPERTY]->(property)
	WITH category, group,property
	OPTIONAL MATCH(category)-[:HAS_SYSTEM_TYPE]->(systemType)
	OPTIONAL MATCH(property)-[:IS_PROPERTY_TYPE]->(propertyType)	
	OPTIONAL MATCH(property)-[:HAS_UNIT]->(unit)
	OPTIONAL MATCH(category)-[:CONTAINS_PHYSICAL_ITEM_PROPERTY]->(physicalItemProperty)
	OPTIONAL MATCH(physicalItemProperty)-[:IS_PROPERTY_TYPE]->(propertyTypePhysical)
	OPTIONAL MATCH(physicalItemProperty)-[:HAS_UNIT]->(unitPhysical)
	WITH category,systemType, group, property, propertyType, physicalItemProperty, propertyTypePhysical, unitPhysical, unit order by id(property)
	WITH category, group, systemType, CASE WHEN physicalItemProperty IS NULL THEN NULL ELSE {
		uid: physicalItemProperty.uid,
		name: physicalItemProperty.name,
		defaultValue: physicalItemProperty.defaultValue,
		type: { uid: propertyTypePhysical.uid, name: propertyTypePhysical.name, code: propertyTypePhysical.code},
		unit: case when unitPhysical is not null then { uid: unitPhysical.uid, name: unitPhysical.name } else null end,
		listOfValues: apoc.text.split(case when physicalItemProperty.listOfValues = "" then null else  physicalItemProperty.listOfValues END, ";")
	} END as physicalItemProperties,	
	collect({
		uid: property.uid, 
		name: property.name,
		defaultValue: property.defaultValue, 
		type: {uid: propertyType.uid, name: propertyType.name, code: propertyType.code}, 		
		unit: case when unit is not null then { uid: unit.uid, name: unit.name } else null end, 
		listOfValues: apoc.text.split(case when property.listOfValues = "" then null else  property.listOfValues END, ";")}) as properties order by id(group)
	WITH category,systemType, CASE WHEN group IS NOT NULL THEN { group: group, properties: properties } ELSE NULL END as groups, physicalItemProperties
	WITH category,systemType, CASE WHEN groups IS NOT NULL THEN  collect({uid: groups.group.uid, name: groups.group.name, properties: groups.properties }) ELSE NULL END as groups,  physicalItemProperties
	WITH { uid: category.uid, name: category.name, code: category.code, miniImageUrl: split(category.miniImageUrl, ";"), groups: groups, physicalItemProperties: collect(physicalItemProperties), systemType: case when systemType is not null then { uid: systemType.uid, name: systemType.name, code: systemType.code } else null end } as category
	return category`

	result.ReturnAlias = "category"

	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid

	return result
}

func CatalogueCategoryPropertiesQuery(uid string) (result helpers.DatabaseQuery) {

	result.Query = `MATCH(category)-[:HAS_SUBCATEGORY*1..50]->(childs:CatalogueCategory{uid:$uid})
	OPTIONAL MATCH(category)-[:HAS_GROUP]->(group)-[:CONTAINS_PROPERTY]->(property)
	WITH category, group,property	
	OPTIONAL MATCH(property)-[:IS_PROPERTY_TYPE]->(propertyType)
	OPTIONAL MATCH(property)-[:HAS_UNIT]->(unit)
	WITH group,property, propertyType, unit order by property.name
	WITH case when count(property) > 0 then { property:{
						uid: property.uid,
						name: property.name, 
						defaultValue: property.defaultValue,
						listOfValues: case when property.listOfValues is not null and property.listOfValues <> "" then apoc.text.split(property.listOfValues, ";") else null end,
						type: { uid: propertyType.uid, name: propertyType.name, code: propertyType.code},
						unit: case when unit is not null then { uid: unit.uid, name: unit.name } else null end 
						},
						propertyGroup: group.name, 
						value: null} else null end as properties
						return properties
	UNION	
	MATCH(category:CatalogueCategory{uid:$uid})
	OPTIONAL MATCH(category)-[:HAS_GROUP]->(group)-[:CONTAINS_PROPERTY]->(property)
	WITH category, group,property 
	OPTIONAL MATCH(property)-[:IS_PROPERTY_TYPE]->(propertyType)
	OPTIONAL MATCH(property)-[:HAS_UNIT]->(unit)
	WITH group,property, propertyType, unit order by property.name
	WITH case when count(property) > 0 then { 
					property:{
						uid: property.uid,
						name: property.name, 
						defaultValue: property.defaultValue,
						listOfValues: case when property.listOfValues is not null and property.listOfValues <> "" then apoc.text.split(property.listOfValues, ";") else null end,
						type: { uid: propertyType.uid, name: propertyType.name, code: propertyType.code},
						unit: case when unit is not null then { uid: unit.uid, name: unit.name } else null end 
						},
						propertyGroup: group.name, 
						value: null} else null end as properties
						return properties;`

	result.ReturnAlias = "properties"
	result.Parameters = make(map[string]interface{}, 1)
	result.Parameters["uid"] = uid

	return result
}

func CatalogueCategoryPhysicalItemPropertiesQuery(uid string) (result helpers.DatabaseQuery) {

	result.Query = `
	MATCH(category:CatalogueCategory{uid:$uid})
	OPTIONAL MATCH(category)-[:CONTAINS_PHYSICAL_ITEM_PROPERTY]->(property)
	WITH category, property 
	OPTIONAL MATCH(property)-[:IS_PROPERTY_TYPE]->(propertyType)
	OPTIONAL MATCH(property)-[:HAS_UNIT]->(unit)
	WITH property, propertyType, unit order by property.name
	WITH case when count(property) > 0 then { 
					property:{
						uid: property.uid,
						name: property.name, 
						defaultValue: property.defaultValue,
						listOfValues: case when property.listOfValues is not null and property.listOfValues <> "" then apoc.text.split(property.listOfValues, ";") else null end,
						type: { uid: propertyType.uid, name: propertyType.name, code: propertyType.code},
						unit: case when unit is not null then { uid: unit.uid, name: unit.name } else null end 
						},
						propertyGroup: null, 
						value: null} else null end as properties
						return properties
						UNION	
	MATCH(category)-[:HAS_SUBCATEGORY*1..50]->(childs:CatalogueCategory{uid:$uid})
	OPTIONAL MATCH(category)-[:CONTAINS_PHYSICAL_ITEM_PROPERTY]->(property)
	WITH category, property 
	OPTIONAL MATCH(property)-[:IS_PROPERTY_TYPE]->(propertyType)
	OPTIONAL MATCH(property)-[:HAS_UNIT]->(unit)
	WITH property, propertyType, unit order by property.name
	WITH case when count(property) > 0 then { 
					property:{
						uid: property.uid,
						name: property.name, 
						defaultValue: property.defaultValue,
						listOfValues: case when property.listOfValues is not null and property.listOfValues <> "" then apoc.text.split(property.listOfValues, ";") else null end,
						type: { uid: propertyType.uid, name: propertyType.name, code: propertyType.code},
						unit: case when unit is not null then { uid: unit.uid, name: unit.name } else null end 
						},
						propertyGroup: null, 
						value: null} else null end as properties
						return properties;`

	result.ReturnAlias = "properties"
	result.Parameters = make(map[string]interface{}, 1)
	result.Parameters["uid"] = uid

	return result
}

func CatalogueCategoryWithDetailsForCopyQuery(uid string) (result helpers.DatabaseQuery) {

	result.Query = `
	MATCH(category:CatalogueCategory{uid:$uid})
	OPTIONAL MATCH(category)-[:HAS_GROUP]->(group)-[:CONTAINS_PROPERTY]->(property)
	WITH category, group,property
	OPTIONAL MATCH(category)-[:HAS_SYSTEM_TYPE]->(systemType)
	OPTIONAL MATCH(property)-[:IS_PROPERTY_TYPE]->(propertyType)
	OPTIONAL MATCH(property)-[:HAS_UNIT]->(unit)
	OPTIONAL MATCH(parent)-[:HAS_SUBCATEGORY]->(category)
	WITH category,systemType, group,property, parent, propertyType, unit order by id(property)
	WITH category,parent,systemType, group, collect({uid: "", name: property.name,defaultValue: property.defaultValue, typeUID: propertyType.uid, unitUID: unit.uid, listOfValues: apoc.text.split(case when property.listOfValues = "" then null else  property.listOfValues END, ";")}) as properties order by id(group)
	WITH category,parent,systemType, CASE WHEN group IS NOT NULL THEN { group: group, properties: properties } ELSE NULL END as groups
	WITH category,parent,systemType, CASE WHEN groups IS NOT NULL THEN  collect({uid: "", name: groups.group.name, properties: groups.properties }) ELSE NULL END as groups	
	WITH { uid: "", name: category.name, code: category.code,miniImageUrl: split(category.miniImageUrl, ";"), groups: groups, parentUID: parent.uid, image: category.image, systemType: case when systemType is not null then { uid: systemType.uid, name: systemType.name, code: systemType.code } else null end } as category
	return category`

	result.ReturnAlias = "category"

	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid

	return result
}

func UpdateCatalogueCategoryQuery(category *models.CatalogueCategory, categoryOld *models.CatalogueCategory) (result helpers.DatabaseQuery) {
	result.Parameters = make(map[string]interface{})

	result.Parameters["name"] = category.Name
	result.Parameters["code"] = category.Code

	result.Query = `
	MERGE(category:CatalogueCategory{uid:$uid})`

	if category.SystemType != nil && category.SystemType.UID != "" {
		result.Parameters["systemTypeUID"] = category.SystemType.UID

		if categoryOld != nil && categoryOld.SystemType != nil {
			result.Parameters["oldSystemTypeUID"] = categoryOld.SystemType.UID
			result.Query += ` WITH category
		MATCH(oldSystemType:SystemType{uid: $oldSystemTypeUID})
		MATCH(category)-[hstOld:HAS_SYSTEM_TYPE]->(oldSystemType)
		DELETE hstOld		
		`
		}

		result.Query += ` WITH category
		MATCH(systemType:SystemType{uid: $systemTypeUID})
		MERGE(category)-[:HAS_SYSTEM_TYPE]->(systemType)
		WITH category `
	} else {
		result.Query += `  WITH category
		OPTIONAL MATCH(category)-[hstOld:HAS_SYSTEM_TYPE]->(oldSystemType)
		DELETE hstOld
		WITH category
		`
	}

	result.Query += `SET category.name = $name, category.code = $code
	`
	if category.Image != "" {
		result.Query += ", category.image = $image "
		if category.Image == "deleted" {
			result.Parameters["image"] = nil
		} else {
			result.Parameters["image"] = category.Image
		}

	}

	// its a new item lets generate new uuid and set parent category
	if category.UID == "" {
		category.UID = uuid.NewString()
		result.Parameters["uid"] = category.UID
		// we could get both parent path or parent uid to identify parent category
		if category.ParentPath != "" {

			if strings.Index(category.ParentPath, "/") == 0 {
				category.ParentPath = strings.Replace(category.ParentPath, "/", "", 1)
			}

			result.Parameters["parentPath"] = category.ParentPath
			result.Query += `WITH category match(parentCategory:CatalogueCategory)	
			optional match(parent)-[:HAS_SUBCATEGORY*1..50]->(parentCategory) 
			with category, parentCategory, apoc.text.join(reverse(collect(parent.code)),"/") as parentPath
			with category, parentCategory, case when parentPath = "" then parentCategory.code else parentPath + "/" + parentCategory.code end as path
			where path = $parentPath
			with parentCategory, category
			MERGE(parentCategory)-[:HAS_SUBCATEGORY]->(category)
			WITH category
			`
		} else if category.ParentUID != "" {
			result.Parameters["parentUID"] = category.ParentUID
			result.Query += `WITH category match(parentCategory:CatalogueCategory{uid: $parentUID})				
			MERGE(parentCategory)-[:HAS_SUBCATEGORY]->(category)
			WITH category
			`
		}
	} else {
		result.Parameters["uid"] = category.UID
	}

	// check if some group exits - then try to merge all the groups
	if category.Groups != nil {
		//merge all groups
		for idg, group := range category.Groups {
			idGroupString := strconv.Itoa(idg)
			result.Parameters["group_name_"+idGroupString] = group.Name

			groupUid := group.UID
			if groupUid == "" {
				groupUid = uuid.NewString()
			}

			result.Query += fmt.Sprintf("MERGE(group%d:CatalogueCategoryPropertyGroup{uid: '%s'}) SET group%d.name=$group_name_%d  MERGE(category)-[:HAS_GROUP]->(group%d) ", idg, groupUid, idg, idg, idg)
			// merge all properties
			if group.Properties != nil {
				for idp, property := range group.Properties {
					propertyUID := property.UID
					if propertyUID == "" {
						propertyUID = uuid.NewString()
					}
					idPropString := idGroupString + "_" + strconv.Itoa(idp)
					result.Parameters["prop_name_"+idPropString] = property.Name
					result.Parameters["prop_defaultValue_"+idPropString] = property.DefaultValue
					result.Parameters["prop_listOfValues_"+idPropString] = strings.Join(property.ListOfValues, ";")

					result.Query += fmt.Sprintf("MERGE(prop_%s:CatalogueCategoryProperty{uid: '%s'}) SET prop_%s.name=$prop_name_%s, prop_%s.defaultValue=$prop_defaultValue_%s, prop_%s.listOfValues=$prop_listOfValues_%s MERGE(group%d)-[:CONTAINS_PROPERTY]->(prop_%s) ", idPropString, propertyUID, idPropString, idPropString, idPropString, idPropString, idPropString, idPropString, idg, idPropString)

					result.Query += fmt.Sprintf("WITH group%[2]d, prop_%[1]v, category OPTIONAL MATCH(prop_%[1]v)-[r_prop_%[1]v:IS_PROPERTY_TYPE]->() DELETE r_prop_%[1]v WITH group%[2]d, prop_%[1]v, category ", idPropString, idg)

					result.Query += fmt.Sprintf("MATCH(type%s:CatalogueCategoryPropertyType{uid:'%s'}) MERGE(prop_%s)-[:IS_PROPERTY_TYPE]->(type%s) ", idPropString, property.Type.UID, idPropString, idPropString)

					result.Query += fmt.Sprintf("WITH group%[2]d, prop_%[1]v, category OPTIONAL MATCH(prop_%[1]v)-[r_prop_%[1]v:HAS_UNIT]->() DELETE r_prop_%[1]v WITH group%[2]d, prop_%[1]v, category ", idPropString, idg)
					if property.Unit != nil && property.Unit.UID != "" {
						result.Query += fmt.Sprintf("MERGE(unit%s:Unit{uid:'%s'}) MERGE(prop_%s)-[:HAS_UNIT]->(unit%s) ", idPropString, property.Unit.UID, idPropString, idPropString)
					}
				}

			}
		}
	}

	// merge all physical item properties
	if category.PhysicalItemProperties != nil {
		for idp, property := range category.PhysicalItemProperties {
			propertyUID := property.UID
			if propertyUID == "" {
				propertyUID = uuid.NewString()
			}
			idPropString := "PIP_" + strconv.Itoa(idp)
			result.Parameters["prop_name_"+idPropString] = property.Name
			result.Parameters["prop_defaultValue_"+idPropString] = property.DefaultValue
			result.Parameters["prop_listOfValues_"+idPropString] = strings.Join(property.ListOfValues, ";")

			result.Query += fmt.Sprintf("MERGE(prop_%s:CatalogueCategoryProperty{uid: '%s'}) SET prop_%s.name=$prop_name_%s, prop_%s.defaultValue=$prop_defaultValue_%s, prop_%s.listOfValues=$prop_listOfValues_%s MERGE(category)-[:CONTAINS_PHYSICAL_ITEM_PROPERTY]->(prop_%s) ", idPropString, propertyUID, idPropString, idPropString, idPropString, idPropString, idPropString, idPropString, idPropString)

			result.Query += fmt.Sprintf("WITH prop_%[1]v, category OPTIONAL MATCH(prop_%[1]v)-[r_prop_%[1]v:IS_PROPERTY_TYPE]->() DELETE r_prop_%[1]v WITH prop_%[1]v, category ", idPropString)

			result.Query += fmt.Sprintf("MATCH(type%s:CatalogueCategoryPropertyType{uid:'%s'}) MERGE(prop_%s)-[:IS_PROPERTY_TYPE]->(type%s) ", idPropString, property.Type.UID, idPropString, idPropString)

			result.Query += fmt.Sprintf("WITH prop_%[1]v, category OPTIONAL MATCH(prop_%[1]v)-[r_prop_%[1]v:HAS_UNIT]->() DELETE r_prop_%[1]v WITH prop_%[1]v, category ", idPropString)
			if property.Unit != nil && property.Unit.UID != "" {
				result.Query += fmt.Sprintf("MERGE(unit%s:Unit{uid:'%s'}) MERGE(prop_%s)-[:HAS_UNIT]->(unit%s) ", idPropString, property.Unit.UID, idPropString, idPropString)
			}
		}

	}

	//only if updating existing item
	if categoryOld != nil {
		result.Query += "WITH category "
		//process deleted items
		deletedGroupsUids, deletedPropsUids := getCatalogueCategoryDeletedItems(category, categoryOld)

		if len(deletedGroupsUids) > 0 {
			result.Parameters["groupsToDelte"] = deletedGroupsUids
			result.Query += "MATCH(groupsToDelete:CatalogueCategoryPropertyGroup) WHERE groupsToDelete.uid IN $groupsToDelte "
		}
		if len(deletedPropsUids) > 0 {
			result.Parameters["propsToDelte"] = deletedPropsUids
			result.Query += "MATCH(propsToDelete:CatalogueCategoryProperty) WHERE propsToDelete.uid IN $propsToDelte "
		}
		if len(deletedGroupsUids) > 0 {
			result.Query += "DETACH DELETE groupsToDelete "
		}
		if len(deletedPropsUids) > 0 {
			result.Query += "DETACH DELETE propsToDelete "
		}
	}
	result.Query += "return category.uid as uid limit 1"

	result.ReturnAlias = "uid"

	return result
}

func GetCatalogueCategoryItemsCountRecursiveQuery(categoryUID string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH (n:CatalogueCategory{uid:$categoryUID})
	WiTH n
	OPTIONAL MATCH(n)-[:HAS_SUBCATEGORY*1..50]->(subs)<-[:BELONGS_TO_CATEGORY*1..20]-(items) 
	with count(items) as subItemsCount, n
	OPTIONAL MATCH (n)<-[:BELONGS_TO_CATEGORY]-(myItems) 
	WITH subItemsCount, count(myItems) as myItemsCount
	return subItemsCount + myItemsCount as result`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["categoryUID"] = categoryUID
	return result
}

func getCatalogueCategoryDeletedItems(category *models.CatalogueCategory, categoryOld *models.CatalogueCategory) (deletedGroupsUids []string, deletedPropsUids []string) {

	for _, oldGroup := range categoryOld.Groups {
		existsGroup := existsGroupByUid(category, oldGroup.UID)
		if existsGroup == nil {
			deletedGroupsUids = append(deletedGroupsUids, oldGroup.UID)
			for _, oldProperty := range oldGroup.Properties {
				deletedPropsUids = append(deletedPropsUids, oldProperty.UID)
			}
		} else {
			for _, oldProperty := range oldGroup.Properties {
				existsProperty := existsPropertyByUid(existsGroup, oldProperty.UID)
				if !existsProperty {
					deletedPropsUids = append(deletedPropsUids, oldProperty.UID)
				}
			}
		}
	}

	// check for deleted physical item properties
	for _, oldProperty := range categoryOld.PhysicalItemProperties {
		existsProperty := existsPhysicalItemPropertyByUid(category, oldProperty.UID)
		if !existsProperty {
			deletedPropsUids = append(deletedPropsUids, oldProperty.UID)
		}
	}

	return deletedGroupsUids, deletedPropsUids
}

func existsGroupByUid(category *models.CatalogueCategory, uid string) (result *models.CatalogueCategoryPropertyGroup) {

	for _, group := range category.Groups {
		if group.UID == uid {
			result = &group
			break
		}
	}

	return result
}

func existsPropertyByUid(category *models.CatalogueCategoryPropertyGroup, uid string) (result bool) {

	for _, property := range category.Properties {
		if property.UID == uid {
			result = true
			break
		}
	}

	return result
}

func existsPhysicalItemPropertyByUid(category *models.CatalogueCategory, uid string) (result bool) {

	for _, property := range category.PhysicalItemProperties {
		if property.UID == uid {
			result = true
			break
		}
	}

	return result
}

func GetUnitsCodebookQuery() (result helpers.DatabaseQuery) {
	result.Query = `MATCH(r:Unit) RETURN {uid: r.uid,name:r.name} as units ORDER BY units.name`
	result.ReturnAlias = "units"
	result.Parameters = make(map[string]interface{})
	return result
}

func GetPropertyTypesCodebookQuery() (result helpers.DatabaseQuery) {
	result.Query = `MATCH(r:CatalogueCategoryPropertyType) RETURN {uid: r.uid,name:r.name} as result ORDER BY id(r)`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	return result
}

func CatalogueSubCategoriesByParentQuery(uid string) (result helpers.DatabaseQuery) {
	result.Query = `MATCH(category:CatalogueCategory{uid:$uid})
OPTIONAL MATCH p = (category)-[:HAS_SUBCATEGORY*1..20]->(childs)
with collect(p) as paths
CALL apoc.convert.toTree(paths, true, { nodes: {CatalogueCategory: ['uid']}, rels:{HAS_SUBCATEGORY: [""]}}) yield value as result
return result`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid
	return result
}

func CatalogueCategoriesForAutocompleteQuery(search string, limit int) (result helpers.DatabaseQuery) {

	result.Query = `MATCH (n:CatalogueCategory)
	WHERE NOT (n)-[:HAS_SUBCATEGORY]->()
	WITH n
	WHERE toLower(n.name) CONTAINS toLower($search)
	OPTIONAL MATCH (parent)-[:HAS_SUBCATEGORY*1..50]->(n)
	WITH n, collect(parent.name) AS parentNames
	RETURN {uid: n.uid, name: n.name + " - " + apoc.text.join(reverse(parentNames), " > ")} AS result
	ORDER BY result.name
	LIMIT $limit;`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["search"] = search
	result.Parameters["limit"] = limit

	return result
}

func ManufacturersForAutocompletQuery(search string, limit int) (result helpers.DatabaseQuery) {

	result.Query = `MATCH (n:Manufacturer)
	WHERE toLower(n.name) STARTS WITH toLower($search)
	RETURN {uid: n.uid, name: n.name} AS result
	ORDER BY result.name
	LIMIT $limit;`
	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["search"] = search
	result.Parameters["limit"] = limit

	return result
}

// save new catalogue item query
func NewCatalogueItemQuery(item *models.CatalogueItem, userUID string) (result helpers.DatabaseQuery) {

	result.Parameters = make(map[string]interface{})
	result.Parameters["name"] = strings.TrimSpace(item.Name)
	result.Parameters["description"] = item.Description
	result.Parameters["categoryUid"] = item.Category.UID
	result.Parameters["userUID"] = userUID
	result.Parameters["manufacturerUrl"] = item.ManufacturerUrl
	result.Parameters["catalogueNumber"] = strings.TrimSpace(item.CatalogueNumber)

	result.Query = `
	MATCH(u:User{uid: $userUID})
	WITH u
	MATCH(cat:CatalogueCategory{uid: $categoryUid})
	WITH u, cat
	CREATE(item:CatalogueItem{  uid: apoc.create.uuid(),
								name: $name, 
								catalogueNumber: $catalogueNumber,
								description: $description,	
								lastUpdateTime: datetime(),						
								manufacturerUrl: $manufacturerUrl })
	CREATE(item)-[:BELONGS_TO_CATEGORY]->(cat)
	CREATE(item)-[:WAS_UPDATED_BY{ at: datetime(), action: "INSERT" }]->(u)
	`
	if item.Supplier != nil && item.Supplier.UID != "" {
		result.Query += `
		WITH item
		MATCH(s:Supplier{uid: $supplierUid})
		CREATE(item)-[:HAS_SUPPLIER]->(s)
		`
		result.Parameters["supplierUid"] = item.Supplier.UID
	}

	for idxProp, prop := range item.Details {

		propIdx := fmt.Sprintf("prop%d", idxProp)
		propUidIdx := fmt.Sprintf("propUID%d", idxProp)
		propValueIdx := fmt.Sprintf("propValue%d", idxProp)
		propValueRelIdx := fmt.Sprintf("propValueRel%d", idxProp)
		result.Parameters[propUidIdx] = prop.Property.UID

		result.Query += fmt.Sprintf(`
		WITH item
		MATCH(%[4]s:CatalogueCategoryProperty{uid: $%[2]s})
		MERGE(item)-[%[1]s:HAS_CATALOGUE_PROPERTY]->(%[4]s)
		SET %[1]s.value = $%[3]s
		`, propValueRelIdx, propUidIdx, propValueIdx, propIdx)

		if prop.Value != nil && prop.Value != "" {

			// check prop.Value type, if its map[string]interface{} then marshal it to json string
			if prop.Property.Type.Code == "range" {
				jsonValue, err := json.Marshal(prop.Value)
				if err == nil {
					result.Parameters[propValueIdx] = string(jsonValue)
				} else {
					// this will throw error in the query
					result.Parameters[propValueIdx] = prop.Value
				}
			} else {
				result.Parameters[propValueIdx] = prop.Value
			}
		} else {
			result.Parameters[propValueIdx] = ""
		}
	}

	result.Query += `RETURN DISTINCT item.uid as uid;`

	result.ReturnAlias = "uid"

	return result
}

// save existing catalogue item query
func UpdateCatalogueItemQuery(item *models.CatalogueItem, oldItem *models.CatalogueItem, userUID string) (result helpers.DatabaseQuery) {

	result.Parameters = make(map[string]interface{})
	result.Parameters["userUID"] = userUID
	result.Parameters["uid"] = item.Uid

	result.Query = `
	MATCH(u:User{uid: $userUID})
	WITH u	
	MATCH(item:CatalogueItem{uid: $uid})	
	SET item.lastUpdateTime = datetime()
	CREATE(item)-[:WAS_UPDATED_BY{ at: datetime(), action: "UPDATE" }]->(u)
	WITh item
	`

	helpers.AutoResolveObjectToUpdateQuery(&result, *item, *oldItem, "item")

	for idxProp, prop := range item.Details {

		propIdx := fmt.Sprintf("prop%d", idxProp)
		propUidIdx := fmt.Sprintf("propUID%d", idxProp)
		propValueIdx := fmt.Sprintf("propValue%d", idxProp)
		propValueRelIdx := fmt.Sprintf("propValueRel%d", idxProp)
		result.Parameters[propUidIdx] = prop.Property.UID

		result.Query += fmt.Sprintf(`
		WITH item
		MATCH(%[4]s:CatalogueCategoryProperty{uid: $%[2]s})
		MERGE(item)-[%[1]s:HAS_CATALOGUE_PROPERTY]->(%[4]s)
		SET %[1]s.value = $%[3]s
		`, propValueRelIdx, propUidIdx, propValueIdx, propIdx)

		if prop.Value != nil && prop.Value != "" {

			// check prop.Value type, if its map[string]interface{} then marshal it to json string
			if prop.Property.Type.Code == "range" {
				jsonValue, err := json.Marshal(prop.Value)
				if err == nil {
					result.Parameters[propValueIdx] = string(jsonValue)
				} else {
					// this will throw error in the query
					result.Parameters[propValueIdx] = prop.Value
				}
			} else {
				result.Parameters[propValueIdx] = prop.Value
			}
		} else {
			result.Parameters[propValueIdx] = ""
		}
	}

	//finally delete all properties that are not in the new item
	for _, oldProp := range oldItem.Details {
		delete := true
		for _, newProp := range item.Details {
			if oldProp.Property.UID == newProp.Property.UID {
				delete = false
				break
			}
		}
		if delete {
			result.Query += fmt.Sprintf(`
			WITH item
			MATCH(propToDelete:CatalogueCategoryProperty{uid: "%s"})
			OPTIONAL MATCH(item)-[r_propToDelete:HAS_CATALOGUE_PROPERTY]->(propToDelete)
			DELETE r_propToDelete
			`, oldProp.Property.UID)
		}
	}

	result.Query += `RETURN DISTINCT item.uid as uid;`

	result.ReturnAlias = "uid"

	return result
}

// delete catalogue item query
func DeleteCatalogueItemQuery(itemUid string, userUID string) (result helpers.DatabaseQuery) {

	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = itemUid
	result.Parameters["userUID"] = userUID

	result.Query = `
	MATCH(u:User{uid: $userUID})
	OPTIONAL MATCH(ci:CatalogueItem{uid: $uid})<-[:IS_BASED_ON]-(i) 
	WITH count(i) as itemsCount, ci, u
	CALL apoc.do.when(itemsCount = 0, 'MATCH(ci:CatalogueItem{uid: uid}) DETACH DELETE ci CREATE(i:CatalogueItemDeleted{uid: uid})-[:WAS_UPDATED_BY{ at: datetime(), action: "DELETE" }]->(u) return itemsCount', 'return itemsCount', {itemsCount: itemsCount, u:u, uid: $uid})
	YIELD value	
	RETURN value.itemsCount as itemsCount;`

	result.ReturnAlias = "itemsCount"

	return result
}

func CatalogueCategoriesTreeQuery(search string) (result helpers.DatabaseQuery) {

	//get catalogue categories as tree
	result.Query = `
	MATCH (parentCat:CatalogueCategory)
	WHERE NOT (parentCat)<-[:HAS_SUBCATEGORY]-()
	WITH parentCat
	OPTIONAL MATCH path = (parentCat)-[:HAS_SUBCATEGORY*0..50]->(children)
	WHERE toLower(parentCat.name) contains $search or toLower(children.name) contains $search
	WITH collect(path) AS paths
	CALL apoc.convert.toTree(paths, true, {
	  nodes: {CatalogueCategory: ['uid', 'name']}
	})
	YIELD value
	RETURN value as tree;`

	result.ReturnAlias = "tree"
	result.Parameters = make(map[string]interface{})
	result.Parameters["search"] = strings.ToLower(search)

	return result
}

func CatalogueItemStatisticsQuery(uid string) (result helpers.DatabaseQuery) {

	result.Query = `
	MATCH (ci:CatalogueItem{uid: $uid})<-[:IS_BASED_ON]-(itm)<-[:HAS_ORDER_LINE]-(o{deleted: false})-[:BELONGS_TO_FACILITY]->(f)
	OPTIONAL MATCH (itm)-[:HAS_ITEM_USAGE]->(usg)
	WITH
	f.name AS facility,
	COALESCE(COUNT(itm), 0) AS facilityCount,
	SUM(CASE WHEN usg.code = "spare-part" THEN 1 ELSE 0 END) AS sparePartsCount,
	SUM(CASE WHEN usg.code = "in-system-part" THEN 1 ELSE 0 END) AS inSystemPartsCount,
	SUM(CASE WHEN usg.code = "experimental-loan-pool-part" THEN 1 ELSE 0 END) AS experimentalLoanPoolPartsCount,
	SUM(CASE WHEN usg.code = "test-and-measurement-equipment" THEN 1 ELSE 0 END) AS testAndMeasurementPartsCount,
	SUM(CASE WHEN usg.code = "stock-item" THEN 1 ELSE 0 END) AS stockItemsCount,
	SUM(CASE WHEN usg.code = "other" OR usg IS NULL THEN 1 ELSE 0 END) AS othersCount
	RETURN{
		facilityName: facility, 
		total: facilityCount,
		sparePartsCount: sparePartsCount,
		inSystemPartsCount: inSystemPartsCount,
		experimentalLoanPoolPartsCount: experimentalLoanPoolPartsCount,
		testAndMeasurementPartsCount: testAndMeasurementPartsCount, 
		stockItemsCount: stockItemsCount,
		othersCount: othersCount} as result;`

	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid

	return result
}

func CatalogueItemsOverallStatisticsQuery() (result helpers.DatabaseQuery) {

	result.Query = `
	MATCH (ci)<-[:IS_BASED_ON]-(itm)<-[:HAS_ORDER_LINE]-(o{deleted: false})-[:BELONGS_TO_FACILITY]->(f)
	OPTIONAL MATCH (itm)-[:HAS_ITEM_USAGE]->(usg)
	WITH
	f.name AS facility,
	COALESCE(COUNT(itm), 0) AS facilityCount,
	SUM(CASE WHEN usg.code = "spare-part" THEN 1 ELSE 0 END) AS sparePartsCount,
	SUM(CASE WHEN usg.code = "in-system-part" THEN 1 ELSE 0 END) AS inSystemPartsCount,
	SUM(CASE WHEN usg.code = "experimental-loan-pool-part" THEN 1 ELSE 0 END) AS experimentalLoanPoolPartsCount,
	SUM(CASE WHEN usg.code = "test-and-measurement-equipment" THEN 1 ELSE 0 END) AS testAndMeasurementPartsCount,
	SUM(CASE WHEN usg.code = "stock-item" THEN 1 ELSE 0 END) AS stockItemsCount,
	SUM(CASE WHEN usg.code = "other" OR usg IS NULL THEN 1 ELSE 0 END) AS othersCount
	RETURN{
		facilityName: facility, 
		total: facilityCount,
		sparePartsCount: sparePartsCount,
		inSystemPartsCount: inSystemPartsCount,
		experimentalLoanPoolPartsCount: experimentalLoanPoolPartsCount,
		testAndMeasurementPartsCount: testAndMeasurementPartsCount, 
		stockItemsCount: stockItemsCount,
		othersCount: othersCount} as result;`

	result.ReturnAlias = "result"

	return result
}
