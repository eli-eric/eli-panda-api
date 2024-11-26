package models

type Publication struct {
	Uid              string `json:"uid"`
	ArticleTitle     string `json:"articleTitle"`
	LongJournalTitle string `json:"longJournalTitle"`
	PdfFile          string `json:"pdfFile"`
	Abstract         string `json:"abstract"`
	Keywords         string `json:"keywords"`
	PublicationDOI   string `json:"publicationDOI"`
	Year             string `json:"year"`
	Pages            int    `json:"pages"`
}
