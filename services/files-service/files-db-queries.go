package filesservice

import (
	"panda/apigateway/helpers"
	"panda/apigateway/services/files-service/models"
	"strings"
)

func GetFileLinksByParentUidQuery(parentUid string) (result helpers.DatabaseQuery) {

	result.Query = `MERGE(n{uid: $parentUid})-[:HAS_FILE_LINK]->(fl)
					RETURN {
							uid: fl.uid, 
							name: fl.name, 
							tags: case when fl.tags is not null and fl.tags <> "" then apoc.text.split(fl.tags, ";") else null end, 
							url: fl.url} as result`

	result.Parameters = map[string]interface{}{}
	result.Parameters["parentUid"] = parentUid
	result.ReturnAlias = "result"

	return result
}

func CreateFileLinkQuery(parentUid string, fileLink models.FileLink) (result helpers.DatabaseQuery) {

	result.Query = `MERGE(n{uid: $parentUid})
	CREATE(fl:FileLink{ 
		uid: apoc.create.uuid(), 
		name: $name, 
		url: $url, 
		tags: $tags })
	CREATE(n)-[:HAS_FILE_LINK{createdAt: datetime()}]->(fl)
	RETURN {
		uid: fl.uid, 
		name: fl.name, 
		tags: case when fl.tags is not null and fl.tags <> "" then apoc.text.split(fl.tags, ";") else null end, 
		url: fl.url} as result`

	result.Parameters = map[string]interface{}{}
	result.Parameters["parentUid"] = parentUid
	result.Parameters["name"] = fileLink.Name
	result.Parameters["url"] = fileLink.Url
	result.Parameters["tags"] = strings.Join(fileLink.Tags, ";")

	result.ReturnAlias = "result"

	return result
}

func UpdateFileLinkQuery(fileLink models.FileLink) (result helpers.DatabaseQuery) {

	result.Query = `MATCH(fl:FileLink{uid: $uid})
					SET fl.name = $name, fl.url = $url, fl.tags = $tags
					RETURN {
							uid: fl.uid,
							name: fl.name,
							tags: case when fl.tags is not null and fl.tags <> "" then apoc.text.split(fl.tags, ";") else null end,
							url: fl.url} as result`

	result.Parameters = map[string]interface{}{}
	result.Parameters["uid"] = fileLink.UID
	result.Parameters["name"] = fileLink.Name
	result.Parameters["url"] = fileLink.Url
	result.Parameters["tags"] = strings.Join(fileLink.Tags, ";")

	result.ReturnAlias = "result"

	return result
}

func DeleteFileLinkQuery(uid string) (result helpers.DatabaseQuery) {

	result.Query = `MATCH(fl:FileLink{uid: $uid}) DETACH DELETE fl`

	result.Parameters = map[string]interface{}{}
	result.Parameters["uid"] = uid

	return result
}

func SetMiniImageUrlToNodeQuery(uid string, imageUrls *[]string, nodeLabel string) (result helpers.DatabaseQuery) {

	if nodeLabel == "CatalogueItem" {
		result.Query = `
		MATCH(n{uid: $uid})
		OPTIONAL MATCH (n)<-[:IS_BASED_ON]-(physicalItem)<-[:CONTAINS_ITEM]-(sys:System)	
		SET n.miniImageUrl = $imageUrl, sys.miniImageUrl = coalesce(sys.miniImageUrl, $imageUrl)
		RETURN 1 as result`
	} else if nodeLabel == "System" {
		result.Query = `
		MATCH(n{uid: $uid})
		SET n.miniImageUrl = $imageUrl
		RETURN 1 as result`
	} else if nodeLabel == "CatalogueCategory" {
		result.Query = `
		MATCH(n{uid: $uid})
		SET n.miniImageUrl = $imageUrl
		WITH n
		OPTIONAL MATCH (n)<-[:BELONGS_TO_CATEGORY]-(ci:CatalogueItem)
		OPTIONAL MATCH (n)-[:HAS_SUBCATEGORY*0..20]->(subs)<-[:BELONGS_TO_CATEGORY]-(subCi:CatalogueItem)
		WITH n, collect(ci) + collect(subCi) AS allCatalogueItems
		UNWIND allCatalogueItems AS ci
		SET ci.miniImageUrl = coalesce(ci.miniImageUrl, $imageUrl)
		WITH ci
		OPTIONAL MATCH (ci)<-[:IS_BASED_ON]-(physicalItem)<-[:CONTAINS_ITEM]-(sys:System)
		SET sys.miniImageUrl = coalesce(sys.miniImageUrl, $imageUrl)
		RETURN 1 as result`
	}

	result.Parameters = map[string]interface{}{}
	result.Parameters["uid"] = uid
	// convert slice of strings to single string
	var imageUrl *string
	if imageUrls != nil {
		joinedUrls := strings.Join(*imageUrls, ";")
		imageUrl = &joinedUrls
	} else {
		imageUrl = nil
	}

	result.Parameters["imageUrl"] = imageUrl

	result.ReturnAlias = "result"

	return result
}
