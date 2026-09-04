package models

import codebookModels "panda/apigateway/services/codebook-service/models"

// WosPreviewResponse is the PANDA-owned response returned by the Web of Science
// preview endpoint. It never persists a publication.
type WosPreviewResponse struct {
	Status                  string                  `json:"status"`
	Doi                     string                  `json:"doi"`
	ExistingPublication     *WosExistingPublication `json:"existingPublication,omitempty"`
	Values                  *WosImportValues        `json:"values,omitempty"`
	Authors                 []WosImportAuthor       `json:"authors"`
	MissingImportableFields []string                `json:"missingImportableFields"`
	UnavailableFields       []string                `json:"unavailableFields"`
}

// WosExistingPublication identifies the active PANDA record that already owns a DOI.
type WosExistingPublication struct {
	Uid   string `json:"uid"`
	Code  string `json:"code"`
	Title string `json:"title"`
	Doi   string `json:"doi"`
}

// WosImportValues contains only fields that can be safely mapped to the
// publication form. Pointer fields distinguish unavailable values from zeroes.
type WosImportValues struct {
	Doi               *string                  `json:"doi,omitempty"`
	Title             *string                  `json:"title,omitempty"`
	WosNumber         *string                  `json:"wosNumber,omitempty"`
	LongJournalTitle  *string                  `json:"longJournalTitle,omitempty"`
	Volume            *int                     `json:"volume,omitempty"`
	Issue             *int                     `json:"issue,omitempty"`
	Pages             *string                  `json:"pages,omitempty"`
	PagesCount        *int                     `json:"pagesCount,omitempty"`
	YearOfPublication *string                  `json:"yearOfPublication,omitempty"`
	DateOfPublication *string                  `json:"dateOfPublication,omitempty"`
	Issn              *string                  `json:"issn,omitempty"`
	EIssn             *string                  `json:"eissn,omitempty"`
	Isbn              *string                  `json:"isbn,omitempty"`
	WebLink           *string                  `json:"webLink,omitempty"`
	Keywords          *string                  `json:"keywords,omitempty"`
	AllAuthors        *string                  `json:"allAuthors,omitempty"`
	AllAuthorsCount   *int                     `json:"allAuthorsCount,omitempty"`
	MediaTypeCb       *codebookModels.Codebook `json:"mediaTypeCb,omitempty"`
}

type WosImportAuthor struct {
	SourceIndex  int            `json:"sourceIndex"`
	DisplayName  string         `json:"displayName"`
	WosStandard  string         `json:"wosStandard,omitempty"`
	ResearcherID string         `json:"researcherId,omitempty"`
	Match        WosAuthorMatch `json:"match"`
}

type WosAuthorMatch struct {
	Kind       string                   `json:"kind"`
	Candidates []WosResearcherCandidate `json:"candidates"`
}

// WosResearcherCandidate is a PANDA researcher a Web of Science author might
// be. It carries the researcher's current ResearcherID so the review dialog can
// say what promoting the incoming one would replace, rather than asking the
// user to toggle a value they cannot see.
type WosResearcherCandidate struct {
	Uid                 string `json:"uid"`
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	CurrentResearcherID string `json:"currentResearcherId,omitempty"`
}

// RememberResearcherIDRequest explicitly links one Web of Science ResearcherID
// to an existing PANDA researcher. MakePrimary additionally promotes it to the
// researcher's current ID, which is the one RIV export sends to the government.
type RememberResearcherIDRequest struct {
	ResearcherID string `json:"researcherId"`
	MakePrimary  bool   `json:"makePrimary,omitempty"`
}

// ResearcherIDsResponse reports every ID the researcher now holds, which of
// them is current, and — when the caller did not promote — whether a newer one
// is sitting unused. SuggestedPrimary is empty when the current ID is already
// the newest, or when the newest cannot be determined without a human.
type ResearcherIDsResponse struct {
	ResearcherIDs    []string `json:"researcherIds"`
	PrimaryID        string   `json:"primaryResearcherId"`
	SuggestedPrimary string   `json:"suggestedPrimaryResearcherId,omitempty"`
}

// PublicationAPIError is the stable error envelope returned by publication
// integration endpoints.
type PublicationAPIError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Retryable bool   `json:"retryable"`
}
