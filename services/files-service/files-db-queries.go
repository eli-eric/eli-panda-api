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
	return result
}

func DeleteFileLinkQuery(uid string) (result helpers.DatabaseQuery) {
	return result
}
