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

func SetMiniImageUrlToNodeQuery(uid string, imageUrl string) (result helpers.DatabaseQuery) {

	result.Query = `MATCH(n{uid: $uid})
					SET n.miniImageUrl = $imageUrl
					RETURN n`

	result.Parameters = map[string]interface{}{}
	result.Parameters["uid"] = uid
	result.Parameters["imageUrl"] = imageUrl

	result.ReturnAlias = "n"

	return result
}
