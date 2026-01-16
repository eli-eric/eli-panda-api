package models

import (
	codebookModels "panda/apigateway/services/codebook-service/models"
	"time"
)

type Publication struct {
	Uid                     string                   `json:"uid" neo4j:"key,uid"`                                                                // uid is the unique identifier of the publication
	Doi                     string                   `json:"doi" neo4j:"prop,doi"`                                                               // doi is the unique identifier of the publication
	Code                    string                   `json:"code" neo4j:"prop,code"`                                                             // code is the internal code of the publication
	Title                   string                   `json:"title" neo4j:"prop,title"`                                                           // title is the title of the publication
	Abstract                string                   `json:"abstract" neo4j:"prop,abstract"`                                                     // abstract is the abstract of the publication
	MediaType               string                   `json:"mediaType" neo4j:"prop,mediaType"`                                                   // mediaType is the media type of the publication
	LongJournalTitle        string                   `json:"longJournalTitle" neo4j:"prop,longJournalTitle"`                                     // longJournalTitle is the long journal title of the publication
	ShortJournalTitle       *string                  `json:"shortJournalTitle" neo4j:"prop,shortJournalTitle"`                                   // shortJournalTitle is the short journal title of the publication
	Volume                  int                      `json:"volume" neo4j:"prop,volume"`                                                         // volume is the volume of the publication
	Issue                   *int                     `json:"issue" neo4j:"prop,issue"`                                                           // issue is the issue of the publication
	Pages                   string                   `json:"pages" neo4j:"prop,pages"`                                                           // pages is the pages of the publication
	PagesCount              int                      `json:"pagesCount" neo4j:"prop,pagesCount"`                                                 // pagesCount is the pages count of the publication
	CiteAs                  string                   `json:"citeAs" neo4j:"prop,citeAs"`                                                         // citeAs is the citation of the publication
	ImpactFactor            *float64                 `json:"impactFactor" neo4j:"prop,impactFactor"`                                             // impactFactor is the impact factor of the publication
	QuartilBasis            *string                  `json:"quartilBasis" neo4j:"prop,quartilBasis"`                                             // quartilBasis is the quartil basis of the publication
	Quartil                 *string                  `json:"quartil" neo4j:"prop,quartil"`                                                       // quartil is the quartil of the publication
	YearOfPublication       string                   `json:"yearOfPublication" neo4j:"prop,yearOfPublication"`                                   // yearOfPublication is the year of publication of the publication
	PdfFileName             *string                  `json:"pdfFileName" neo4j:"prop,pdfFileName"`                                               // pdfFileName is the name of the pdf file of the publication
	PdfFileUrl              *string                  `json:"pdfFileUrl" neo4j:"prop,pdfFileUrl"`                                                 // pdfFileUrl is the url of the pdf file of the publication
	DateOfPublication       *string                  `json:"dateOfPublication" neo4j:"prop,dateOfPublication"`                                   // dateOfPublication is the date of publication of the publication
	Keywords                string                   `json:"keywords" neo4j:"prop,keywords"`                                                     // keywords is the keywords of the publication
	OecdFord                *string                  `json:"oecdFord" neo4j:"prop,oecdFord"`                                                     // oecdFord is the oecd ford of the publication
	WosNumber               *string                  `json:"wosNumber" neo4j:"prop,wosNumber"`                                                   // wosNumber is the wos number of the publication
	Issn                    *string                  `json:"issn" neo4j:"prop,issn"`                                                             // issn is the issn of the publication
	EIssn                   *string                  `json:"eissn" neo4j:"prop,eissn"`                                                           // eissn is the eissn of the publication
	WebLink                 string                   `json:"webLink" neo4j:"prop,webLink"`                                                       // webLink is the web link of the publication
	EidScopus               *string                  `json:"eidScopus" neo4j:"prop,eidScopus"`                                                   // eidScopus is the eid scopus of the publication
	Language                string                   `json:"language" neo4j:"prop,language"`                                                     // language is the language of the publication
	Grant                   *string                  `json:"grant" neo4j:"prop,grant"`                                                           // grant is the grant of the publication
	Note                    *string                  `json:"note" neo4j:"prop,note"`                                                             // note is the note of the publication
	AllAuthors              string                   `json:"allAuthors" neo4j:"prop,allAuthors"`                                                 // allAuthors is the all authors of the publication
	AllAuthorsCount         int                      `json:"allAuthorsCount" neo4j:"prop,allAuthorsCount"`                                       // allAuthorsCount is the all authors count of the publication
	EliAuthors              string                   `json:"eliAuthors" neo4j:"prop,eliAuthors"`                                                 // eliAuthors is the eli authors of the publication
	EliAuthorsCount         int                      `json:"eliAuthorsCount" neo4j:"prop,eliAuthorsCount"`                                       // eliAuthorsCount is the eli authors count of the publication
	AuthorsDepartments      []AuthorsDepartment      `json:"authorsDepartments"`                                                                 // authorsDepartments is the authors departments of the publication
	AuthorsDepartmentsArray []string                 `json:"authorsDepartmentsArray" neo4j:"prop,authorsDepartmentsArray"`                       // authorsDepartmentsArray is the authors departments array of the publication
	OpenAccessType          *codebookModels.Codebook `json:"openAccessType" neo4j:"rel,OpenAccessType,HAS_OPEN_ACCESS_TYPE,uid,openAccessType"`  // openAccessType is the open access type of the publication
	PublishingCountry       *codebookModels.Codebook `json:"publishingCountry" neo4j:"rel,Country,HAS_PUBLISHING_COUNTRY,uid,publishingCountry"` // publishingCountry is the publishing country of the publication
	UserCall                *codebookModels.Codebook `json:"userCall" neo4j:"rel,UserCall,HAS_USER_CALL,uid,userCall"`                           // userCall is the user call of the publication
	UserExperiment          *string                  `json:"userExperiment" neo4j:"prop,userExperiment"`                                         // userExperiment is the user experiment of the publication
	ExperimentalSystem      *string                  `json:"experimentalSystem" neo4j:"prop,experimentalSystem"`                                 // experimentalSystem is the experimental system of the publication
	UpdatedAt               *time.Time               `json:"updatedAt" neo4j:"prop,updatedAt"`                                                   // updatedAt is the time when the publication was last updated
	EliResearchers          []ResearcherRef          `json:"eliResearchers"`                                                                     // eliResearchers are the connected researchers via HAS_RESEARCHER relationship
}

type AuthorsDepartment struct {
	Department   codebookModels.Codebook `json:"department" neo4j:"rel,Department,HAS_DEPARTMENT,uid,department"` // department is the department of the publication
	AuthorsCount int                     `json:"authorsCount" neo4j:"prop,authorsCount"`                          // authorsCount is the authors count of the publication
}

type ResearcherRef struct {
	Uid       string `json:"uid"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
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
