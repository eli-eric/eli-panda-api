package catalogueService

import (
	"fmt"
	"panda/apigateway/helpers"
	"panda/apigateway/services/catalogue-service/models"
	"strconv"
	"strings"

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

func CatalogueCategoryImageByUidQuery(uid string) (result helpers.DatabaseQuery) {

	result.Query = `match(category:CatalogueCategory{uid: $uid})
	
	return category.image as image`

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
	WITH category, group, collect({uid: property.uid, name: property.name,defaultValue: property.defaultValue, typeUID: propertyType.uid, unitUID: unit.uid, listOfValues: apoc.text.split(case when property.listOfValues = "" then null else  property.listOfValues END, ";")}) as properties
	WITH category, CASE WHEN group IS NOT NULL THEN { group: group, properties: properties } ELSE NULL END as groups
	WITH category, CASE WHEN groups IS NOT NULL THEN  collect({uid: groups.group.uid, name: groups.group.name, properties: groups.properties }) ELSE NULL END as groups	
	WITH { uid: category.uid, name: category.name, code: category.code, groups: groups } as category
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
	MERGE(category:CatalogueCategory{uid:$uid})
	SET category.name = $name, category.code = $code
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
		result.Parameters["uid"] = uuid.NewString()
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
					if property.TypeUID != "" {
						result.Query += fmt.Sprintf("MERGE(type%s:CatalogueCategoryPropertyType{uid:'%s'}) MERGE(prop_%s)-[:IS_PROPERTY_TYPE]->(type%s) ", idPropString, property.TypeUID, idPropString, idPropString)
					}
					if property.UnitUID != "" {
						result.Query += fmt.Sprintf("MERGE(unit%s:Unit{uid:'%s'}) MERGE(prop_%s)-[:HAS_UNIT]->(unit%s) ", idPropString, property.UnitUID, idPropString, idPropString)
					}
				}

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
