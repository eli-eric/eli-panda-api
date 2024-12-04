package publicationsservice

import (
	"panda/apigateway/helpers"
)

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
