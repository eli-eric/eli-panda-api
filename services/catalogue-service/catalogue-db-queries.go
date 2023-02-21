package catalogueService

import (
	"fmt"
	"panda/apigateway/helpers"
	"panda/apigateway/services/catalogue-service/models"
	"strconv"

	"github.com/google/uuid"
)

// Get catalogue items with pagination and filters
func CatalogueItemsFiltersPaginationQuery(search string, categoryPath string, skip int, limit int) (result helpers.DatabaseQuery) {

	result.Query = `MATCH(category:CatalogueCategory)
	with category
	optional match(parent)-[:HAS_SUBCATEGORY*1..15]->(category) 
	with category, apoc.text.join(reverse(collect(parent.code)),"/") + "/" + category.code as categoryPath
	where $categoryPath = '' or (categoryPath = $categoryPath or categoryPath = '/' + $categoryPath)
	optional match(category)-[:HAS_SUBCATEGORY*1..15]->(subs)		
	with categoryPath, collect(subs.uid) as subCategoryUids, category.uid as itmCategoryUid
	match(itm:CatalogueItem)-[:BELONGS_TO_CATEGORY]->(cat)		
	where cat.uid in subCategoryUids or cat.uid = itmCategoryUid		
	OPTIONAL MATCH(itm)-[:HAS_MANUFACTURER]->(manu)
	OPTIONAL MATCH(itm)-[propVal:HAS_CATALOGUE_PROPERTY]->(prop)
	OPTIONAL MATCH(prop)-[:HAS_UNIT]->(unit)
	OPTIONAL MATCH(prop)-[:IS_PROPERTY_TYPE]->(propType)
	OPTIONAL MATCH(group)-[:CONTAINS_PROPERTY]->(prop)
	WITH itm,cat,categoryPath, propType.code as propTypeCode, manu, prop.name as propName, group.name as groupName, toString(propVal.value) as value, unit.name as unit
	ORDER BY itm.name
	WHERE $searchText = '' or 
	(toLower(itm.name) CONTAINS $searchText OR toLower(itm.description) CONTAINS $searchText or toLower(manu.name) CONTAINS $searchText or toLower(value) CONTAINS $searchText)
	RETURN {
	uid: itm.uid,
	name: itm.name,
	description: itm.description,
	categoryName: cat.name,
	categoryPath: max(categoryPath),
	manufacturer: manu.name,
	manufacturerUrl: itm.manufacturerUrl,
	manufacturerNumber: itm.catalogueNumber,
	details: collect({ propertyName: propName, propertyType: propTypeCode,propertyUnit: unit, propertyGroup: groupName, value: value})
	} as items
	SKIP $skip
	LIMIT $limit`

	result.ReturnAlias = "items"

	result.Parameters = make(map[string]interface{})
	result.Parameters["searchText"] = search
	result.Parameters["skip"] = skip
	result.Parameters["limit"] = limit
	result.Parameters["categoryPath"] = categoryPath

	return result
}

func CatalogueItemsFiltersTotalCountQuery(search string, categoryPath string) (result helpers.DatabaseQuery) {

	result.Query = `MATCH(category:CatalogueCategory)
	with category
	optional match(parent)-[:HAS_SUBCATEGORY*1..15]->(category) 
	with category, apoc.text.join(reverse(collect(parent.code)),"/") + "/" + category.code as categoryPath
	where $categoryPath = '' or (categoryPath = $categoryPath or categoryPath = '/' + $categoryPath)
	optional match(category)-[:HAS_SUBCATEGORY*1..15]->(subs)		
	with categoryPath, collect(subs.uid) as subCategoryUids, category.uid as itmCategoryUid
	match(itm:CatalogueItem)-[:BELONGS_TO_CATEGORY]->(cat)		
	where cat.uid in subCategoryUids or cat.uid = itmCategoryUid		
	OPTIONAL MATCH(itm)-[:HAS_MANUFACTURER]->(manu)
	OPTIONAL MATCH(itm)-[propVal:HAS_CATALOGUE_PROPERTY]->(prop)
	OPTIONAL MATCH(prop)-[:HAS_UNIT]->(unit)
	OPTIONAL MATCH(prop)-[:IS_PROPERTY_TYPE]->(propType)
	OPTIONAL MATCH(group)-[:CONTAINS_PROPERTY]->(prop)
	WITH itm,cat,categoryPath, propType.code as propTypeCode, manu, prop.name as propName, group.name as groupName, toString(propVal.value) as value, unit.name as unit
	ORDER BY itm.name
	WHERE $searchText = '' or 
	(toLower(itm.name) CONTAINS $searchText OR toLower(itm.description) CONTAINS $searchText or toLower(manu.name) CONTAINS $searchText or toLower(value) CONTAINS $searchText)
	RETURN count(distinct itm.uid) as itemsCount`

	result.ReturnAlias = "itemsCount"

	result.Parameters = make(map[string]interface{})
	result.Parameters["searchText"] = search
	result.Parameters["categoryPath"] = categoryPath

	return result
}

func CatalogueCategoriesByParentPathQuery(parentPath string) (result helpers.DatabaseQuery) {

	result.Query = `match(category:CatalogueCategory)
	with category
	optional match(parent)-[:HAS_SUBCATEGORY*1..50]->(category) 
	with category, apoc.text.join(reverse(collect(parent.code)),"/") as parentPath
	where parentPath = $parentPath
	return {uid:category.uid,code:category.code, name:category.name,parentPath: parentPath} as categories`

	result.ReturnAlias = "categories"

	result.Parameters = make(map[string]interface{})
	result.Parameters["parentPath"] = parentPath

	return result
}

func CatalogueItemWithDetailsByUidQuery(uid string) (result helpers.DatabaseQuery) {

	result.Query = `match(itm:CatalogueItem{uid: $uid})-[:BELONGS_TO_CATEGORY]->(cat)			
	OPTIONAL MATCH(itm)-[:HAS_MANUFACTURER]->(manu)
	OPTIONAL MATCH(itm)-[propVal:HAS_CATALOGUE_PROPERTY]->(prop)
	OPTIONAL MATCH(prop)-[:HAS_UNIT]->(unit)
	OPTIONAL MATCH(prop)-[:IS_PROPERTY_TYPE]->(propType)
	OPTIONAL MATCH(group)-[:CONTAINS_PROPERTY]->(prop)
	WITH itm,cat, propType.code as propTypeCode, manu, prop.name as propName, group.name as groupName, toString(propVal.value) as value, unit.name as unit
	RETURN {
	uid: itm.uid,
	name: itm.name,
	description: itm.description,
	categoryName: cat.name,
	manufacturer: manu.name,
	manufacturerUrl: itm.manufacturerUrl,
	manufacturerNumber: itm.catalogueNumber,
	details: collect({ propertyName: propName, propertyType: propTypeCode,propertyUnit: unit, propertyGroup: groupName, value: value})
	} as catalogueItem`

	result.ReturnAlias = "catalogueItem"

	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid

	return result
}

func CatalogueCategoryWithDetailsQuery(uid string) (result helpers.DatabaseQuery) {

	result.Query = `MATCH(category:CatalogueCategory{uid:$uid})
	OPTIONAL MATCH(category)-[:HAS_GROUP]->(group)-[:CONTAINS_PROPERTY]->(property)
	WITH category, group,property
	OPTIONAL MATCH(property)-[:IS_PROPERTY_TYPE]->(propertyType)
	OPTIONAL MATCH(property)-[:HAS_UNIT]->(unit)
	WITH category, group, collect({uid: property.uid, name: property.name,default: property.defaultValue, typeUID: propertyType.uid, unitUID: unit.uid, listOfValues: apoc.text.split(case when property.listOfValues = "" then null else  property.listOfValues END, ";")}) as properties
	WITH category, { group: group, properties: properties } as groups
	WITH category, collect({uid: groups.group.uid, name: groups.group.name, properties: groups.properties }) as groups
	WITH { uid: category.uid, name: category.name, code: category.code, groups: groups } as category
	return category`

	result.ReturnAlias = "category"

	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid

	return result
}

func UpdateCatalogueCategoryQuery(category *models.CatalogueCategory) (result helpers.DatabaseQuery) {

	result.Query = `
	MERGE(category:CatalogueCategory{uid:$uid})
	SET category.name = $name, category.code = $code
	`

	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = category.UID
	result.Parameters["name"] = category.Name
	result.Parameters["code"] = category.Code

	//add groups queries - merge
	if category.Groups != nil {
		for idg, group := range category.Groups {
			idGroupString := strconv.Itoa(idg)
			result.Parameters["group_name_"+idGroupString] = group.Name

			if group.UID == "" {
				//new group
				newGroupUid := uuid.NewString()
				result.Query += fmt.Sprintf("MERGE(group%d:CatalogueCategoryPropertyGroup{uid: '%s',name: $group_name_%d}) MERGE(category)-[:HAS_GROUP]->(group%d) ", idg, newGroupUid, idg, idg)

			} else {
				// existing group
				result.Query += fmt.Sprintf("MERGE(group%d:CatalogueCategoryPropertyGroup{uid: '%s'}) SET group%d.name=$group_name_%d MERGE(category)-[:HAS_GROUP]->(group%d) ", idg, group.UID, idg, idg, idg)
			}
		}
	}

	result.Query += "return category.uid as uid"

	result.ReturnAlias = "uid"

	return result
}
