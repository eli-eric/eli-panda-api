package filesservice

import (
	"panda/apigateway/helpers"
	"panda/apigateway/services/files-service/models"
	"strings"
)

func GetFileLinksByParentUidQuery(parentUid string) (result helpers.DatabaseQuery) {

	result.Query = `MATCH(n{uid: $parentUid})-[:HAS_FILE_LINK]->(fl)
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

	result.Query = `MATCH(n{uid: $parentUid})
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

func SetMiniImageUrlToNodeQuery(uid string, imageUrls *[]string, forceAll bool) (result helpers.DatabaseQuery) {

	result.Query = `MATCH(n{uid: $uid})
	SET n.miniImageUrl = $imageUrl
	WITH n
	OPTIONAL MATCH (n)<-[:BELONGS_TO_CATEGORY]-(ci:CatalogueItem)<-[:IS_BASED_ON]-(physicalItem)<-[:CONTAINS_ITEM]-(sys:System)
	WITH n,ci,sys, case when ci.miniImageUrl is not null then ci.miniImageUrl else $imageUrl end as CIminiImageUrl, case when sys.miniImageUrl is not null then sys.miniImageUrl else $imageUrl end as SYSminiImageUrl 
	SET ci.miniImageUrl = CIminiImageUrl , sys.miniImageUrl = SYSminiImageUrl 
	WITH n
	OPTIONAL MATCH (n)<-[:IS_BASED_ON]-(physicalItem)<-[:CONTAINS_ITEM]-(sys:System)
	WITH n,sys, case when sys.miniImageUrl is not null then sys.miniImageUrl else $imageUrl end as SYSminiImageUrl
	SET sys.miniImageUrl = SYSminiImageUrl
	RETURN n.name, n.miniImageUrl`

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

	result.ReturnAlias = "n"

	return result
}
