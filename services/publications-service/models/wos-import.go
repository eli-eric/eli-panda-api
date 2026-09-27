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
	// Warnings describe degraded source values (non-numeric volume, missing
	// day, ...) so the dialog can explain an omitted field instead of
	// silently dropping it (ELIPANDA-501).
	Warnings []WosImportWarning `json:"warnings"`
}

// WosImportWarning points at one form field whose raw upstream value could
// not be mapped losslessly.
type WosImportWarning struct {
	Code    string `json:"code"`
	Field   string `json:"field"`
	Raw     string `json:"raw,omitempty"`
	Message string `json:"message"`
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
	Orcid        string         `json:"orcid,omitempty"`
	Match        WosAuthorMatch `json:"match"`
}

type WosAuthorMatch struct {
	Kind string `json:"kind"`
	// Confidence is the stable dialog contract: EXACT_ID | NAME | AMBIGUOUS |
	// NONE. Kind stays the internal provenance (researcher-id, orcid, name...).
	Confidence string `json:"confidence"`
	// KnownResearcherId is true when the author's WoS ResearcherID is already
	// registered on the matched researcher; false invites the dialog to offer
	// remembering it (ELIPANDA-501).
	KnownResearcherId bool            `json:"knownResearcherId"`
	Candidates        []ResearcherRef `json:"candidates"`
}

// PublicationAPIError is the stable error envelope returned by publication
// integration endpoints.
type PublicationAPIError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Retryable bool   `json:"retryable"`
}

// ResearcherIdsRequest is the body of PATCH /v1/researcher/{uid}/researcher-ids.
type ResearcherIdsRequest struct {
	ResearcherIDs []string `json:"researcherIds"`
}
