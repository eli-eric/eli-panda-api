package publicationsservice

import (
	"panda/apigateway/helpers"
	"panda/apigateway/services/publications-service/models"
)

func UpdatePublicationQuery(newPublication *models.Publication, oldPublication *models.Publication, userUID string) (result helpers.DatabaseQuery) {
	result.Parameters = make(map[string]interface{})

	result.Query = `
		MATCH (n:Publication {uid:$uid}) `

	helpers.AutoResolveObjectToUpdateQuery(&result, *newPublication, *oldPublication, "n")

	result.Query += `	
		WITH n
		MATCH(u:User{uid: $lastUpdateBy})
		WITH n, u
		SET n.lastUpdateTime = datetime(), n.lastUpdateBy = u.username
		WITH n, u
		CREATE(n)-[:WAS_UPDATED_BY{at: datetime(), action: "UPDATE" }]->(u)
		RETURN n.uid as uid
		`

	result.Parameters["uid"] = oldPublication.Uid
	result.Parameters["lastUpdateBy"] = userUID

	result.ReturnAlias = "uid"

	return result
}

func GetPublicationByUidQuery(uid string) (result helpers.DatabaseQuery) {
	result.Query = `
		MATCH (n:Publication {uid:$uid})
		RETURN {
			uid: n.uid,
			articleTitle: n.articleTitle,
			abstract: n.abstract,
			longJournalTitle: n.longJournalTitle,
			keywords: n.keywords,
			year: n.year,
			publicationDOI: n.publicationDOI,
			pages: n.pages,
			pdfFile: n.pdfFile
			} as publication`
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid
	result.ReturnAlias = "publication"

	return result
}

func CreatePublicationQuery(publication *models.Publication) (result helpers.DatabaseQuery) {
	result.Query = `
		CREATE (n:Publication {
			uid: $uid,
			articleTitle: $articleTitle,
			abstract: $abstract,
			longJournalTitle: $longJournalTitle,
			keywords: $keywords,
			year: $year,
			publicationDOI: $publicationDOI,
			pages: $pages,
			pdfFile: $pdfFile
		})
		RETURN n.uid as uid`

	result.ReturnAlias = "uid"

	return result
}

func DeletePublicationQuery(uid string) (result helpers.DatabaseQuery) {
	result.Query = `
		MATCH (n:Publication {uid:$uid})
		DETACH DELETE n`
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid

	return result
}

func GetAllPublicationsQuery() (result helpers.DatabaseQuery) {
	result.Query = `
		MATCH (n:Publication)
		RETURN {
			uid: n.uid,
			articleTitle: n.articleTitle,
			abstract: n.abstract,
			longJournalTitle: n.longJournalTitle,
			keywords: n.keywords,
			year: n.year,
			publicationDOI: n.publicationDOI,
			pages: n.pages,
			pdfFile: n.pdfFile
		} as publication`
	result.ReturnAlias = "publication"

	return result
}
