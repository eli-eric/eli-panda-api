package models

import (
	codebookModels "panda/apigateway/services/codebook-service/models"
	"time"
)

type Publication struct {
	Uid                 string                   `json:"uid" neo4j:"key,uid"`
	Doi                 string                   `json:"doi" neo4j:"prop,doi"`
	ArticleTitle        string                   `json:"articleTitle" neo4j:"prop,articleTitle"`
	Abstract            *string                  `json:"abstract" neo4j:"prop,abstract"`
	JournalTitle        *string                  `json:"journalTitle" neo4j:"prop,journalTitle"`
	Volume              *int                     `json:"volume" neo4j:"prop,volume"`
	Issue               *int                     `json:"issue" neo4j:"prop,issue"`
	PagesFrom           *int                     `json:"pagesFrom" neo4j:"prop,pagesFrom"`
	PagesTo             *int                     `json:"pagesTo" neo4j:"prop,pagesTo"`
	PagesTotal          *int                     `json:"pagesTotal" neo4j:"prop,pagesTotal"`
	Citations           *[]string                `json:"citations" neo4j:"prop,citations"`
	ImpactFactor        *float64                 `json:"impactFactor" neo4j:"prop,impactFactor"`
	Quartile            *string                  `json:"quartile" neo4j:"prop,quartile"`
	Year                *string                  `json:"year" neo4j:"prop,year"`
	PdfFileName         *string                  `json:"pdfFileName" neo4j:"prop,pdfFileName"`
	PdfFileUrl          *string                  `json:"pdfFileUrl" neo4j:"prop,pdfFileUrl"`
	PublishDate         *time.Time               `json:"publishDate" neo4j:"prop,publishDate"`
	Keywords            *string                  `json:"keywords" neo4j:"prop,keywords"`
	OecdFord            *string                  `json:"oecdFord" neo4j:"prop,oecdFord"`
	WosNumber           *string                  `json:"wosNumber" neo4j:"prop,wosNumber"`
	Issn                *string                  `json:"issn" neo4j:"prop,issn"`
	EIssn               *string                  `json:"eissn" neo4j:"prop,eissn"`
	Url                 *string                  `json:"url" neo4j:"prop,url"`
	EidScopus           *string                  `json:"eidScopus" neo4j:"prop,eidScopus"`
	State               *string                  `json:"state" neo4j:"prop,state"`
	Language            *string                  `json:"language" neo4j:"prop,language"`
	UserCall            *codebookModels.Codebook `json:"userCall" neo4j:"rel,UserCall,HAS_USER_CALL,uid,userCall"`
	UserExperiment      *codebookModels.Codebook `json:"userExperiment" neo4j:"rel,UserExperiment,HAS_USER_EXPERIMENT,uid,userExperiment"`
	PublicationCategory *codebookModels.Codebook `json:"publicationCategory" neo4j:"rel,PublicationCategory,HAS_PUBLICATION_CATEGORY,uid,publicationCategory"`
	OpenAccessType      *codebookModels.Codebook `json:"openAccessType" neo4j:"rel,OpenAccessType,HAS_OPEN_ACCESS_TYPE,uid,openAccessType"`
	PublicationSupport  *codebookModels.Codebook `json:"publicationSupport" neo4j:"rel,PublicationSupport,HAS_PUBLICATION_SUPPORT,uid,publicationSupport"`
	StatisticsTotal     *int                     `json:"statisticsTotal" neo4j:"prop,statisticsTotal"`
	StatisticsEric      *int                     `json:"statisticsEric" neo4j:"prop,statisticsEric"`
	StatisticsBeamlines *int                     `json:"statistics Beamlines" neo4j:"prop,statisticsBeamlines"`
	StatisticsAlps      *int                     `json:"statisticsAlps" neo4j:"prop,statisticsAlps"`
}

type WosAPIResponse struct {
	WosMetadata WosMetadata `json:"metadata"`
	WosHits     []WosHit    `json:"hits"`
}

type WosMetadata struct {
	WosTotal int `json:"total"`
	WosPage  int `json:"page"`
	WosLimit int `json:"limit"`
}

type WosHit struct {
	WosUID         string         `json:"uid"`
	WosTitle       string         `json:"title"`
	WosTypes       []string       `json:"types"`
	WosSourceTypes []string       `json:"sourceTypes"`
	WosSource      WosSource      `json:"source"`
	WosNames       WosNames       `json:"names"`
	WosCitations   []WosCitation  `json:"citations"`
	WosIdentifiers WosIdentifiers `json:"identifiers"`
	WosKeywords    WosKeywords    `json:"keywords"`
}

type WosSource struct {
	WosSourceTitle  string   `json:"sourceTitle"`
	WosPublishYear  int      `json:"publishYear"`
	WosPublishMonth string   `json:"publishMonth"`
	WosVolume       string   `json:"volume,omitempty"`
	WosIssue        string   `json:"issue,omitempty"`
	WosPages        WosPages `json:"pages"`
}

type WosPages struct {
	WosRange string `json:"range"`
	WosBegin string `json:"begin"`
	WosEnd   string `json:"end"`
	WosCount int    `json:"count"`
}

type WosNames struct {
	WosAuthors     []WosAuthor `json:"authors"`
	WosBookEditors []WosEditor `json:"bookEditors,omitempty"`
}

type WosAuthor struct {
	WosDisplayName  string `json:"displayName"`
	WosStandard     string `json:"wosStandard"`
	WosResearcherID string `json:"researcherId"`
}

type WosEditor struct {
	WosDisplayName string `json:"displayName"`
}

type WosCitation struct {
	WosDB    string `json:"db"`
	WosCount int    `json:"count"`
}

type WosIdentifiers struct {
	WosDOI  string `json:"doi,omitempty"`
	WosISSN string `json:"issn"`
}

type WosKeywords struct {
	WosAuthorKeywords []string `json:"authorKeywords"`
}
