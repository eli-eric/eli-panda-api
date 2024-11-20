package publicationsservice

import (
	"panda/apigateway/helpers"
	"panda/apigateway/services/publications-service/models"
)

func UpdatePublicationQuery(publication *models.Publication) (result helpers.DatabaseQuery) {
	result.Query = `
		MATCH (n:Publication {uid:$uid})
		SET n.articleTitle = $articleTitle,
			n.abstract = $abstract,		
			n.keywords = $keywords,
			n.year = $year,
			n.publicationDOI = $publicationDOI,
			n.longJournalTitle = $longJournalTitle,
			n.pages = $pages,
			n.pdfFile = $pdfFile
		RETURN n.uid as uid`
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = publication.Uid
	result.Parameters["articleTitle"] = publication.ArticleTitle
	result.Parameters["abstract"] = publication.Abstract
	result.Parameters["keywords"] = publication.Keywords
	result.Parameters["year"] = publication.Year
	result.Parameters["publicationDOI"] = publication.PublicationDOI
	result.Parameters["longJournalTitle"] = publication.LongJournalTitle
	result.Parameters["pages"] = publication.Pages
	result.Parameters["pdfFile"] = publication.PdfFile

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
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = publication.Uid
	result.Parameters["articleTitle"] = publication.ArticleTitle
	result.Parameters["abstract"] = publication.Abstract
	result.Parameters["keywords"] = publication.Keywords
	result.Parameters["year"] = publication.Year
	result.Parameters["publicationDOI"] = publication.PublicationDOI
	result.Parameters["longJournalTitle"] = publication.LongJournalTitle
	result.Parameters["pages"] = publication.Pages
	result.Parameters["pdfFile"] = publication.PdfFile

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
