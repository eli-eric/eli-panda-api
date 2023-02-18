package catalogueService

import "panda/apigateway/helpers"

// Get catalogue items with pagination and filters
func CatalogueItemsPaginationFiltersQuery(search string, categoryPath string, skip int, limit int) (result helpers.DatabaseQuery) {

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

	result.Parameters["searchText"] = search
	result.Parameters["skip"] = skip
	result.Parameters["limit"] = limit
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

	result.Parameters["parentPath"] = parentPath

	return result
}
